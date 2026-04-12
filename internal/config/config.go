package config

import (
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("token", os.Getenv("DISCORD_TOKEN"))
	viper.SetDefault("test-guilds", strings.Split(os.Getenv("TEST_GUILD_ID"), ","))
	viper.SetDefault("messages-file", os.Getenv("MESSAGES_FILE"))
	viper.SetDefault("birthdays-file", os.Getenv("BIRTHDAYS_FILE"))
	viper.SetDefault("advertise-today-birthday", false)
	viper.SetDefault("advertise-channel", os.Getenv("ADVERTISE_CHANNEL"))
	viper.SetDefault("timezone", os.Getenv("TZ"))
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

func GetBirthdaysPath() string {
	return viper.GetString("birthdays-file")
}

func GetAdvertiseTodayBirthday() bool {
	return viper.GetBool("advertise-today-birthday")
}

func GetAdvertiseChannel() string {
	return viper.GetString("advertise-channel")
}

func GetTimezone() *time.Location {
	tz := viper.GetString("timezone")
	if tz == "" {
		tz = "America/Argentina/Buenos_Aires"
	}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		logrus.WithError(err).Warnf("Invalid timezone %s, defaulting to UTC", tz)
		return time.UTC
	}
	return loc
}
