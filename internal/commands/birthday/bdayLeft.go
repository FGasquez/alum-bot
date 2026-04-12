package birthday

import (
	"fmt"
	"strings"

	"github.com/FGasquez/alum-bot/internal/messages"
	"github.com/FGasquez/alum-bot/internal/types"
	"github.com/bwmarrin/discordgo"
)

const BdayLeftCommandName = "bday-left"

var BdayLeftCommand = discordgo.ApplicationCommand{
	Name:        BdayLeftCommandName,
	Description: "Get the days left until a specific user's birthday",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionUser,
			Name:        "user",
			Description: "The user to calculate days left for",
			Required:    true,
		},
	},
}

var BdayLeftCommandHandler = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	var targetUserID string
	var targetUser *discordgo.User

	for _, option := range options {
		if option.Name == "user" {
			targetUser = option.UserValue(s)
			if targetUser != nil {
				targetUserID = targetUser.ID
			}
		}
	}

	pingFormat := fmt.Sprintf("<@%s>", targetUserID)

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

	var targetBirthday *types.Birthday
	if targetUser != nil {
		for _, b := range birthdays {
			bUsername := strings.TrimSpace(b.Username)
			bName := strings.TrimSpace(b.Name)
			
			// Match if the discord ping, discord ID, or plain username match
			matchUsername := bUsername != "" && (bUsername == pingFormat || bUsername == targetUserID || strings.EqualFold(bUsername, targetUser.Username))
			// Match if the Name in YAML matches their Display Name (GlobalName) or plain username
			// Discord global names may be nil sometimes, so defend against it
			globalName := targetUser.GlobalName
			matchName := bName != "" && (strings.EqualFold(bName, targetUser.Username) || (globalName != "" && strings.EqualFold(bName, globalName)))

			if matchUsername || matchName {
				bCopy := b
				targetBirthday = &bCopy
				break
			}
		}
	}

	var parsedBirthday *ParsedBirthday
	if targetBirthday != nil {
		parsedBirthday, err = CalculateDaysLeft(targetBirthday)
	} else {
		err = fmt.Errorf("not found")
	}

	if err != nil || parsedBirthday == nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("❌ Birthday for user <@%s> not found in the registry.", targetUserID),
			},
		})
		return
	}

	messageTemplate := messages.GetMessage(messages.MessageKeys.BdayLeft)

	// In the structure messages expects maps with specific keys
	templateData := map[string]interface{}{
		"Username": parsedBirthday.Username,
		"DaysLeft": parsedBirthday.DaysLeft,
	}

	content := messages.TemplateMessage(messageTemplate, templateData)

	// Support for formatting without pinging if needed, but since we want to ping:
	if !strings.Contains(content, "<@") {
		content = fmt.Sprintf("%s <@%s>", content, targetUserID)
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
}
