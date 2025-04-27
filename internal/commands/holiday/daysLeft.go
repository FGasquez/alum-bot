package holidays

import (
	"time"

	"github.com/FGasquez/alum-bot/internal/helpers"
	"github.com/FGasquez/alum-bot/internal/messages"
	"github.com/FGasquez/alum-bot/internal/types"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

const DaysLeftToHolidayName = "days-left"

var HowManyDaysToHoliday = discordgo.ApplicationCommand{
	Name:        DaysLeftToHolidayName,
	Description: "Get how many days are left for the next holiday",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionBoolean,
			Name:        "skip-today",
			Description: "skip today in the calculation",
			Required:    false,
		},
		{
			Type:        discordgo.ApplicationCommandOptionBoolean,
			Name:        "skip-weekend",
			Description: "skip weekend in the calculation",
			Required:    false,
		},
	},
}

var HowManyDaysToHolidayHandlers = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var skipToday bool = false
	var skipWeekend bool = true

	params := helpers.GetParams(i.ApplicationCommandData().Options)
	if _, ok := params["skip-today"]; ok {
		skipToday = params["skip-today"].(bool)
	}
	if _, ok := params["skip-weekend"]; ok {
		skipWeekend = params["skip-weekend"].(bool)
	}

	daysLeftToHoliday, holiday, isToday := DaysLeft(skipWeekend, skipToday)
	if daysLeftToHoliday == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Es hoy! 🎉",
			},
		})
		return
	}

	if daysLeftToHoliday == -1 {
		logrus.Errorf("Failed to parse holiday date")
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Failed to retrieve the next holiday. Please try again later.",
			},
		})
		return
	}

	dateFormatted, day, month, _ := helpers.FormatDateToSpanishUnparsed(holiday.Date)
	parsedDate, err := time.Parse("2006-01-02", holiday.Date)
	if err != nil {
		logrus.Errorf("Failed to parse holiday date: %v", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Failed to retrieve the next holiday. Please try again later.",
			},
		})
		return
	}

	tmpValues := types.TemplateValues{
		HolidayName:   holiday.Name,
		DaysLeft:      daysLeftToHoliday,
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
		FullDate:  holiday.Date,
		Adjacents: holiday.Adjacent,
		IsToday:   isToday,
	}

	message := messages.TemplateMessage(messages.GetMessage(messages.MessageKeys.DaysLeft), tmpValues)
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})
}
