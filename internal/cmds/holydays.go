package cmds

import (
	"fmt"
	"strconv"
	"time"

	"github.com/FGasquez/alum-bot/internal/helpers"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

const HolydaysCommandName = "next-holyday"
const DaysLeftToHolydayName = "days-left"

var HolydaysCommands = discordgo.ApplicationCommand{
	Name:        HolydaysCommandName,
	Description: "Get the next holyday",
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

var HolydaysCommandHandlers = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var skipToday bool = false
	var skipWeekend bool = true

	if len(i.ApplicationCommandData().Options) == 0 {
		for _, option := range i.ApplicationCommandData().Options {
			switch option.Name {
			case "skip-today":
				skipToday = option.Options[0].BoolValue()
			case "skip-weekend":
				skipWeekend = option.Options[0].BoolValue()
			}

		}
	}

	// logrus.Infof("Getting next holyday, summoned by %s, at guild %s, in channel %s", i.Member.User.Username, i.GuildID, i.ChannelID)
	nextHoliday := helpers.NextHolyday(time.Now(), skipWeekend, skipToday)
	logrus.Infof("Next holyday: %s, date: %s", nextHoliday.Name, nextHoliday.Date)

	// Parse nextHoliday.Date into a time.Time object
	parsedDate, err := time.Parse("2006-01-02", nextHoliday.Date)
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

	day, month, year := helpers.FormatDateToSpanish(parsedDate)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("🎉 El próximo feriado es **%s** el **%s %s %s**. 🎉", nextHoliday.Name, day, month, year),
		},
	})
}

var HowManyDaysToHolyday = discordgo.ApplicationCommand{
	Name:        DaysLeftToHolydayName,
	Description: "Get how many days are left for the next holyday",
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

var HowManyDaysToHolydayHandlers = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var skipToday bool = false
	var skipWeekend bool = true

	// logrus.Infof("Calculating days left, summoned by %s, at guild %s, in channel %s", i.Member.User.Username, i.GuildID, i.ChannelID)

	if len(i.ApplicationCommandData().Options) == 0 {
		for _, option := range i.ApplicationCommandData().Options {
			switch option.Name {
			case "skip-today":
				skipToday = option.Options[0].BoolValue()
			case "skip-weekend":
				skipWeekend = option.Options[0].BoolValue()
			}

		}
	}

	logrus.Info("Getting how many days to the next holyday")
	logrus.Info("Skip today: ", skipToday)
	logrus.Info("Skip weekend: ", skipWeekend)

	daysLeftToHolyday := helpers.DaysLeft()

	if daysLeftToHolyday == -1 {
		logrus.Errorf("Failed to parse holiday date")
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Failed to retrieve the next holiday. Please try again later.",
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("🎉 Para el próximo feriado faltan %s días! 🎉", strconv.Itoa(daysLeftToHolyday)),
		},
	})
}
