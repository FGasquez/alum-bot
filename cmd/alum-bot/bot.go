package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/FGasquez/alum-bot/internal/commands"
	birthdayCmd "github.com/FGasquez/alum-bot/internal/commands/birthday"
	holidaysCmd "github.com/FGasquez/alum-bot/internal/commands/holiday"
	"github.com/FGasquez/alum-bot/internal/config"
	"github.com/FGasquez/alum-bot/internal/messages"
	"github.com/FGasquez/alum-bot/internal/types"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

// Commands and handlers are managed via internal/commands registry.

func removeAllCommands(dg *discordgo.Session, guilds []string) {
	logrus.Info("Pruning all commands and exiting")
	for _, guildID := range guilds {
		existingCommands, err := dg.ApplicationCommands(dg.State.User.ID, guildID)
		if err != nil {
			logrus.WithError(err).Error("Error fetching existing commands for pruning")
			continue
		}

		for _, cmd := range existingCommands {
			logrus.WithField("command", cmd.Name).Info("Pruning command")
			if err := dg.ApplicationCommandDelete(dg.State.User.ID, guildID, cmd.ID); err != nil {
				logrus.WithError(err).Error("Error pruning command")
			}
		}
	}
}

func pruneCommands() {
	token := config.GetToken()
	testGuilds := config.GetTestGuilds()

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		logrus.WithError(err).Error("Error creating Discord session")
		return
	}

	err = dg.Open()
	if err != nil {
		logrus.WithError(err).Error("Error opening Discord session")
		return
	}
	defer dg.Close()

	removeAllCommands(dg, testGuilds)
}

func setActivityStatus(dg *discordgo.Session, message string) {

	logrus.WithField("message", message).Debug("Setting activity status")
	err := dg.UpdateStatusComplex(discordgo.UpdateStatusData{
		Activities: []*discordgo.Activity{
			{
				Name:  message,
				Type:  discordgo.ActivityTypeCustom,
				State: message,
			},
		},
	})
	if err != nil {
		logrus.WithError(err).Error("Error setting activity status")
	}
}

func advertiseBirthday() {
	token := config.GetToken()
	channelID := config.GetAdvertiseChannel()

	if channelID == "" {
		logrus.Error("Advertise channel is not set. Use --advertise-channel or ADVERTISE_CHANNEL env var.")
		os.Exit(1)
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		logrus.WithError(err).Error("Error creating Discord session for advertisement job")
		os.Exit(1)
	}

	sent, err := birthdayCmd.AdvertiseTodayBirthday(dg, channelID)
	if err != nil {
		logrus.WithError(err).Error("Failed to advertise birthday")
		os.Exit(1)
	}
	if sent {
		logrus.Info("Successfully advertised today's birthday!")
	}
	os.Exit(0)
}

func runBot() {
	if config.GetAdvertiseTodayBirthday() {
		if config.GetBirthdaysPath() == "" {
			logrus.Error("Birthdays file is not set. The advertise job requires a birthdays file. Use --birthdays-file or BIRTHDAYS_FILE env var.")
			os.Exit(1)
		}

		logrus.Info("Running advertise birthday job...")
		advertiseBirthday()
		return
	}

	token := config.GetToken()
	testGuilds := config.GetTestGuilds()

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		logrus.WithError(err).Error("Error creating Discord session")
		return
	}

	activeCommands, commandHandlers := commands.GetActiveCommands()

	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == discordgo.InteractionApplicationCommand {
			if handler, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				handler(s, i)
			} else {
				logrus.Warnf("No handler for command: %s", i.ApplicationCommandData().Name)
			}
		}
	})

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		logrus.WithError(err).Error("Error opening connection")
		return
	}

	existingCommands, err := dg.ApplicationCommands(dg.State.User.ID, "")
	if err != nil {
		logrus.WithError(err).Error("Error fetching existing commands")
		return
	}

	existingCommandNames := make(map[string]bool)
	for _, cmd := range existingCommands {
		existingCommandNames[cmd.Name] = true
	}

	var registeredCommands []*discordgo.ApplicationCommand
	for _, guildID := range testGuilds {
		for _, cmd := range activeCommands {
			if existingCommandNames[cmd.Name] {
				logrus.WithField("command", cmd.Name).Info("Command already registered, skipping")
				continue
			}

			rc, err := dg.ApplicationCommandCreate(dg.State.User.ID, guildID, cmd)
			if err != nil {
				logrus.WithError(err).WithField("command", cmd.Name).Error("Cannot create slash command")
			} else {
				logrus.WithField("command", rc.Name).Info("Registered command")
				registeredCommands = append(registeredCommands, rc)
			}
		}
	}

	DaysLeft, _, _ := holidaysCmd.DaysLeft(true, false)
	logrus.Infof("Time to next holiday: %d days", DaysLeft)

	go func() {
		DaysLeft, _, _ := holidaysCmd.DaysLeft(true, false)
		setActivityStatus(dg, messages.TemplateMessage(messages.GetMessage(messages.MessageKeys.ActivityStatus), types.TemplateValues{DaysLeft: DaysLeft}))

		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			DaysLeft, _, _ := holidaysCmd.DaysLeft(true, false)
			setActivityStatus(dg, messages.TemplateMessage(messages.GetMessage(messages.MessageKeys.ActivityStatus), types.TemplateValues{DaysLeft: DaysLeft}))
		}
	}()

	logrus.Info("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	// Cleanly close down the Discord session.
	logrus.Info("Graceful shutdown")
	setActivityStatus(dg, "")
	dg.Close()
}
