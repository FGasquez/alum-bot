package main

import (
	"flag"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/FGasquez/alum-bot/internal/cmds"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

// Variables used for command line parameters
var (
	Token          string
	TestGuildID    string
	RemoveCommands bool
)

func getToken() string {
	if Token != "" {
		return Token
	}
	return os.Getenv("DISCORD_TOKEN")
}

func splitIDs(ids string) []string {
	if ids == "" {
		return []string{}
	}
	return strings.Split(ids, ",")
}

func getTestGuildIDs() []string {
	if TestGuildID != "" {
		return splitIDs(TestGuildID)
	}
	return splitIDs(os.Getenv("TEST_GUILD_ID"))
}

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.StringVar(&TestGuildID, "test-guild-id", "", "Test guild ID")
	flag.BoolVar(&RemoveCommands, "remove-commands", false, "Remove all commands")
	flag.Parse()
}

var commands = []*discordgo.ApplicationCommand{
	&cmds.HolydaysCommands,
	&cmds.HowManyDaysToHolyday,
	{
		Name:        "ping",
		Description: "Ping!",
	},
}

var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	cmds.HolydaysCommandName:   cmds.HolydaysCommandHandlers,
	cmds.DaysLeftToHolydayName: cmds.HowManyDaysToHolydayHandlers,
	"ping": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		// Respond to the ping slash command
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Pong!",
			},
		})
	},
}

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
}

func main() {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + getToken())
	if err != nil {
		logrus.WithError(err).Error("Error creating Discord session")
		return
	}

	testGuids := getTestGuildIDs()
	if len(testGuids) == 0 {
		logrus.Error("No test guilds provided")
		testGuids = append(testGuids, "")
	}

	// Handle interactions (slash commands)
	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == discordgo.InteractionApplicationCommand {
			if handler, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				handler(s, i)
			} else {
				logrus.Warnf("No handler for command: %s", i.ApplicationCommandData().Name)
			}
		}
	})

	// Also handle message events (for legacy commands)
	dg.AddHandler(messageCreate)
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord.
	err = dg.Open()
	if err != nil {
		logrus.WithError(err).Error("Error opening connection")
		return
	}

	// Register the slash commands and store the returned command objects.
	var registeredCommands []*discordgo.ApplicationCommand
	for _, guildID := range testGuids {

		for _, cmd := range commands {
			rc, err := dg.ApplicationCommandCreate(dg.State.User.ID, guildID, cmd)
			if err != nil {
				logrus.WithError(err).WithField("command", cmd.Name).Error("Cannot create slash command")
			} else {
				logrus.WithField("command", rc.Name).Info("Registered command")
				registeredCommands = append(registeredCommands, rc)
			}
		}
	}

	logrus.Info("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Remove commands if the flag is provided.
	if RemoveCommands {
		logrus.Info("Removing all commands")
		for _, guildID := range testGuids {

			for _, c := range registeredCommands {
				logrus.WithField("command", c.Name).Info("Removing command")
				if err := dg.ApplicationCommandDelete(dg.State.User.ID, guildID, c.ID); err != nil {
					logrus.WithError(err).Error("Error deleting command")
				}
			}
		}
	}

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called for every new message.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping", reply with "Pong!"
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}
	// If the message is "pong", reply with "Ping!"
	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}
