package holidays

import (
	"fmt"
	"time"

	"github.com/FGasquez/alum-bot/internal/helpers"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

const HolydaysOfMonthName = "holydays-of-month"

var HolydaysOfMonth = discordgo.ApplicationCommand{
	Name:        HolydaysOfMonthName,
	Description: "Get the holydays of the month",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "month",
			Description: "The month number",
			Required:    true,
			Choices: []*discordgo.ApplicationCommandOptionChoice{
				{
					Name:  "Enero",
					Value: helpers.January,
				},
				{
					Name:  "Febrero",
					Value: helpers.February,
				},
				{
					Name:  "Marzo",
					Value: helpers.March,
				},
				{
					Name:  "Abril",
					Value: helpers.April,
				},
				{
					Name:  "Mayo",
					Value: helpers.May,
				},
				{
					Name:  "Junio",
					Value: helpers.June,
				},
				{
					Name:  "Julio",
					Value: helpers.July,
				},
				{
					Name:  "Agosto",
					Value: helpers.August,
				},
				{
					Name:  "Septiembre",
					Value: helpers.September,
				},
				{
					Name:  "Octubre",
					Value: helpers.October,
				},
				{
					Name:  "Noviembre",
					Value: helpers.November,
				},
				{
					Name:  "Diciembre",
					Value: helpers.December,
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "year",
			Description: "The year",
			Required:    false,
		},
	},
}

var HolydaysOfMonthHandlers = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	month := i.ApplicationCommandData().Options[0].IntValue()
	year := time.Now().Year()

	params := helpers.GetParams(i.ApplicationCommandData().Options)
	if _, ok := params["year"]; ok {
		year = params["year"].(int)
	}

	holydaysOfMonth, err := helpers.GetAllHolidaysOfMonth(helpers.Months(month), year)
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

	if len(holydaysOfMonth) == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("No hay feriados en el mes %d del año %d", month, year),
			},
		})
		return
	}

	var holydays string
	for _, holyday := range holydaysOfMonth {
		parsedDate, err := time.Parse("2006-01-02", holyday.Date)
		if err != nil {
			logrus.Errorf("Failed to parse holiday date: %v", err)
			continue
		}
		day, _, _ := helpers.FormatDateToSpanish(parsedDate)
		holydays += fmt.Sprintf("* **%s** - %s\n", holyday.Name, day)
	}

	// TODO: Detect holidays adjacent to weekends and show in the message
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Feriados del mes %d del año %d:\n%s", month, year, holydays),
		},
	})
}
