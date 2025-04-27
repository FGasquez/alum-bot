package config

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("token", os.Getenv("DISCORD_TOKEN"))
	viper.SetDefault("test-guilds", strings.Split(os.Getenv("TEST_GUILD_ID"), ","))
	viper.SetDefault("messages-file", os.Getenv("MESSAGES_FILE"))
}

func GetToken() string {
	return viper.GetString("token")
}

func GetTestGuilds() []string {
	return viper.GetStringSlice("test-guilds")
}

func GetMessagesPath() string {
	return viper.GetString("messages-file")
}
