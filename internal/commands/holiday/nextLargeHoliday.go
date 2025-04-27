package holidays

import (
	"time"

	"github.com/FGasquez/alum-bot/internal/helpers"
	"github.com/FGasquez/alum-bot/internal/messages"
	"github.com/FGasquez/alum-bot/internal/types"
	"github.com/bwmarrin/discordgo"
)

const HolidaysLargeCommandName = "next-large-holiday"

var HolidaysLargeCommands = discordgo.ApplicationCommand{
	Name:        HolidaysLargeCommandName,
	Description: "Get the next large holiday",
}

func GetNextLargeHoliday(holiday types.ProcessedHolidays) *types.ParsedHolidays {
	for _, holiday := range holiday.All {
		if len(holiday.Adjacent) > 0 {
			return &holiday
		}
	}

	return nil
}

var HolidayLargeCommandHandlers = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	holidays, err := GetHolidays(time.Now().Year(), true)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Failed to retrieve holidays. Please try again later.",
			},
		})

		return
	}

	largeHolidays := GetNextLargeHoliday(holidays)

	if largeHolidays == nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ No upcoming large holidays found.",
			},
		})
		return
	}

	dateFormatted, day, month, _ := helpers.FormatDateToSpanishUnparsed(largeHolidays.Date)

	parsedDate, err := time.Parse("2006-01-02", largeHolidays.Date)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Failed to retrieve holidays. Please try again later.",
			},
		})

		return
	}

	tmpValues := types.TemplateValues{
		HolidayName:   largeHolidays.Name,
		DaysLeft:      largeHolidays.DaysLeftToHoliday,
		FormattedDate: dateFormatted,
		NamedDate: types.NamedDate{
			Day:   day,
			Month: month,
		},
		RawDate: types.RawDate{
			Day:   parsedDate.Day(),
			Month: int(parsedDate.Month()),
			Year:  parsedDate.Year(),
		},
		FullDate:  largeHolidays.Date,
		Adjacents: largeHolidays.Adjacent,
		IsToday:   largeHolidays.IsToday,
	}

	message := messages.TemplateMessage(messages.GetMessage(messages.MessageKeys.NextLargeHoliday), tmpValues)
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})
}
