package birthday

import (
	"github.com/bwmarrin/discordgo"
)

const NextBdayCommandName = "next-bday"

var NextBdayCommand = discordgo.ApplicationCommand{
	Name:        NextBdayCommandName,
	Description: "Get the next upcoming birthday in the group",
}

var NextBdayCommandHandler = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	birthdays, err := GetBirthdays()
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Failed to retrieve the birthdays list. Check if the configuration file is provided.",
			},
		})
		return
	}

	nextBday, err := GetNextBirthday(birthdays)
	if err != nil || nextBday == nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ No upcoming birthdays found or invalid file formatting.",
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: FormatNextBirthdayMessage(nextBday, s, i.GuildID),
		},
	})
}
