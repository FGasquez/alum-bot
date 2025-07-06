package holidays

import (
	"time"

	"github.com/FGasquez/alum-bot/internal/helpers"
	"github.com/FGasquez/alum-bot/internal/messages"
	"github.com/FGasquez/alum-bot/internal/types"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

const HolidaysOfMonthName = "holidays-of-month"

var HolidaysOfMonth = discordgo.ApplicationCommand{
	Name:        HolidaysOfMonthName,
	Description: "Get the holidays of the month",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "month",
			Description: "The month number",
			Required:    true,
			Choices: []*discordgo.ApplicationCommandOptionChoice{
				{
					Name:  "Enero",
					Value: January,
				},
				{
					Name:  "Febrero",
					Value: February,
				},
				{
					Name:  "Marzo",
					Value: March,
				},
				{
					Name:  "Abril",
					Value: April,
				},
				{
					Name:  "Mayo",
					Value: May,
				},
				{
					Name:  "Junio",
					Value: June,
				},
				{
					Name:  "Julio",
					Value: July,
				},
				{
					Name:  "Agosto",
					Value: August,
				},
				{
					Name:  "Septiembre",
					Value: September,
				},
				{
					Name:  "Octubre",
					Value: October,
				},
				{
					Name:  "Noviembre",
					Value: November,
				},
				{
					Name:  "Diciembre",
					Value: December,
				},
			},
		},
	},
}

var HolidaysOfMonthHandlers = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	month := i.ApplicationCommandData().Options[0].IntValue()
	monthName := helpers.MonthsToSpanish(month)
	year := time.Now().Year()

	// TODO: Fix year param
	params := helpers.GetParams(i.ApplicationCommandData().Options)
	if _, ok := params["year"]; ok {
		year = int(params["year"].(int))
	}

	holidaysOfMonth, err := GetAllHolidaysOfMonth(Months(month), year)
	if err != nil {
		logrus.Errorf("Failed to retrieve holidays of the month: %v", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: messages.TemplateMessage(messages.GetMessage(messages.MessageKeys.FailedToParseHolidayDate), nil),
			},
		})
		return
	}

	if len(holidaysOfMonth) == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: messages.TemplateMessage(messages.GetMessage(messages.MessageKeys.NoHolidaysOfMonth), map[string]interface{}{
					"Month": monthName,
				}),
			},
		})
		return
	}

	var adjacentHolidays [][]types.ParsedHolidays
	adjacentMap := make(map[string]bool)

	for _, holiday := range holidaysOfMonth {
		var currentAdjacents []types.ParsedHolidays

		for _, adjacent := range holiday.Adjacent {
			if adjacentMap[adjacent.Date] {
				continue
			}

			currentAdjacents = append(currentAdjacents, adjacent)
			adjacentMap[adjacent.Date] = true
		}

		if len(currentAdjacents) > 1 { // <<-- Important: at least 2 to consider it a long holiday
			adjacentHolidays = append(adjacentHolidays, currentAdjacents)
		}
	}
	holidaysOfMonthFiltered := make([]types.ParsedHolidays, 0, len(holidaysOfMonth))
	for _, holiday := range holidaysOfMonth {
		if holiday.Type != types.Weekend {
			holidaysOfMonthFiltered = append(holidaysOfMonthFiltered, holiday)
		}
	}
	tmpValues := types.MonthTemplateValues{
		Month:        monthName,
		HolidaysList: holidaysOfMonthFiltered,
		Adjacents:    adjacentHolidays,
		Count:        len(holidaysOfMonthFiltered),
	}

	message := messages.TemplateMessage(messages.GetMessage(messages.MessageKeys.HolidaysOfMonth), tmpValues)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})
}
