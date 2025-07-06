package main

import (
	"log"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Variables used for command line parameters
var (
	Token          string
	TestGuildID    string
	RemoveCommands bool
	PruneCommands  bool
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "alum-bot",
		Short: "Discord bot for keeping track of holidays in Argentina",
		Run: func(cmd *cobra.Command, args []string) {
			runBot()
		},
	}

	rootCmd.PersistentFlags().StringP("token", "t", "", "Bot token (default: DISCORD_TOKEN)")
	rootCmd.PersistentFlags().StringSliceP("test-guilds", "g", []string{}, "List of test guild IDs (default: TEST_GUILD_ID)")
	rootCmd.PersistentFlags().String("messages-file", "", "Path to messages file (default: '')")

	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		log.Fatal(err)
	}

	rootCmd.AddCommand(&cobra.Command{
		Use:   "prune-commands",
		Short: "Prune all commands and exit",
		Run: func(cmd *cobra.Command, args []string) {
			pruneCommands()
		},
	})

	viper.BindPFlags(rootCmd.PersistentFlags())
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
