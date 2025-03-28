package holidays

import (
	"fmt"
	"time"

	"github.com/FGasquez/alum-bot/internal/helpers"
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

	holidaysOfMonth, adjacent, err := GetAllHolidaysOfMonth(Months(month), year)
	if err != nil {
		logrus.Errorf("Failed to retrieve holidays of the month: %v", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Failed to retrieve the holidays of the month. Please try again later.",
			},
		})
		return
	}

	if len(holidaysOfMonth) == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("No hay feriados en el mes **%s** del año %d", monthName, year),
			},
		})
		return
	}

	var holidays string
	for _, holiday := range holidaysOfMonth {
		parsedDate, err := time.Parse("2006-01-02", holiday.Date)
		if err != nil {
			logrus.Errorf("Failed to parse holiday date: %v", err)
			continue
		}
		day, _, _ := helpers.FormatDateToSpanish(parsedDate)
		holidays += fmt.Sprintf("* **%s** - %s\n", holiday.Name, day)
	}

	var adjacentHolidays []string
	for _, adj := range adjacent {
		firstDay, err := time.Parse("2006-01-02", adj[0].Date)
		if err != nil {
			logrus.Errorf("Failed to parse holiday date: %v", err)
			continue
		}

		lastDay, err := time.Parse("2006-01-02", adj[len(adj)-1].Date)
		if err != nil {
			logrus.Errorf("Failed to parse holiday date: %v", err)
			continue
		}
		adjacentHolidays = append(adjacentHolidays, fmt.Sprintf("* Del %d al %d\n", firstDay.Day(), lastDay.Day()))

	}

	var message string
	message = fmt.Sprintf("Feriados del mes **%s** del año %d:\n%s", monthName, year, holidays)
	if len(adjacentHolidays) > 0 {
		message += "\n\nFeriados largos:\n"
		for _, adj := range adjacentHolidays {
			message += adj
		}
	}

	// TODO: Detect holidays adjacent to weekends and show in the message
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})
}
