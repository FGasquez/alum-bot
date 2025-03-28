package holidays

import (
	"fmt"
	"time"

	"github.com/FGasquez/alum-bot/internal/helpers"
	"github.com/bwmarrin/discordgo"
)

const HolidaysLargeCommandName = "next-large-holiday"

var HolidaysLargeCommands = discordgo.ApplicationCommand{
	Name:        HolidaysLargeCommandName,
	Description: "Get the next large holiday",
}

func GetNextLargeHoliday(holidays [][]Holiday) []Holiday {
	nextLargeHoliday := holidays[0]
	for _, holiday := range holidays {
		parsedDate, err := time.Parse("2006-01-02", holiday[0].Date)
		if err != nil {
			continue
		}
		if parsedDate.After(time.Now()) {
			nextLargeHoliday = holiday
			break
		}
	}

	return nextLargeHoliday

}

var HolidayLargeCommandHandlers = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	holidays, err := GetHolidays(time.Now().Year())
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Failed to retrieve holidays. Please try again later.",
			},
		})

		return
	}

	largeHolidays := largeHolidays(holidays)

	if len(largeHolidays) == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ No upcoming large holidays found.",
			},
		})
		return
	}

	nextLargeHoliday := GetNextLargeHoliday(largeHolidays)
	date, err := time.Parse("2006-01-02", nextLargeHoliday[0].Date)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Failed to retrieve holidays. Please try again later.",
			},
		})

		return
	}

	_, month, _ := helpers.FormatDateToSpanish(date)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("El siguiente feriado largo arranca el día **%d de %s** y dura un total de **%d días**", date.Day(), month, len(nextLargeHoliday)),
		},
	})
}
