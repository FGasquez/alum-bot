package birthday

import (
	"github.com/FGasquez/alum-bot/internal/helpers"
	"github.com/FGasquez/alum-bot/internal/messages"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

// BirthdayTemplateData holds the data for birthday message templates.
type BirthdayTemplateData struct {
	Name     string
	Username string
	Date     string
	DaysLeft int
}

// NewBirthdayTemplateData builds template data from a ParsedBirthday,
// substituting a resolved mention if provided, or bolding the name.
func NewBirthdayTemplateData(b *ParsedBirthday, s *discordgo.Session, guildID string) BirthdayTemplateData {
	username := b.Username
	if username == "" {
		username = "**" + b.Name + "**"
	} else if s != nil && guildID != "" {
		username = helpers.ResolveUsernameToMention(s, guildID, username)
	}

	return BirthdayTemplateData{
		Name:     b.Name,
		Username: username,
		Date:     b.Date,
		DaysLeft: b.DaysLeft,
	}
}

// FormatNextBirthdayMessage formats the "next birthday" message using the
// messages package, appending a suffix based on how many days are left.
func FormatNextBirthdayMessage(b *ParsedBirthday, s *discordgo.Session, guildID string) string {
	data := NewBirthdayTemplateData(b, s, guildID)
	msg := messages.TemplateMessage(messages.GetMessage(messages.MessageKeys.NextBday), data)

	switch {
	case b.DaysLeft == 0:
		msg = messages.TemplateMessage(messages.GetMessage(messages.MessageKeys.TodayBday), data)
	case b.DaysLeft == 1:
		msg += "\n(It's tomorrow! 🎂)"
	default:
		msg += "\n(" + itoa(b.DaysLeft) + " days left)"
	}
	return msg
}

// FormatTodayBirthdayMessage formats a standalone "today is their birthday" message.
func FormatTodayBirthdayMessage(b *ParsedBirthday, s *discordgo.Session, guildID string) string {
	data := NewBirthdayTemplateData(b, s, guildID)
	return messages.TemplateMessage(messages.GetMessage(messages.MessageKeys.TodayBday), data)
}

// AdvertiseTodayBirthday checks if today is someone's birthday and, if so,
// sends a message to the given Discord channel. Returns true if a message was sent.
func AdvertiseTodayBirthday(dg *discordgo.Session, channelID string) (bool, error) {
	channel, err := dg.State.Channel(channelID)
	if err != nil {
		channel, err = dg.Channel(channelID)
		if err != nil {
			logrus.WithError(err).Errorf("Failed to fetch channel ID %s for advertisement", channelID)
			return false, err
		}
	}

	birthdays, err := GetBirthdays()
	if err != nil {
		return false, err
	}

	nextBday, err := GetNextBirthday(birthdays)
	if err != nil || nextBday == nil {
		logrus.Info("No upcoming birthdays found or invalid file formatting.")
		return false, nil
	}

	if nextBday.DaysLeft != 0 {
		logrus.Infof("Next birthday is in %d days, not advertising today.", nextBday.DaysLeft)
		return false, nil
	}

	msg := FormatTodayBirthdayMessage(nextBday, dg, channel.GuildID)
	if _, err := dg.ChannelMessageSend(channelID, msg); err != nil {
		return false, err
	}
	return true, nil
}

// itoa converts an int to its decimal string representation without importing strconv.
func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	neg := i < 0
	if neg {
		i = -i
	}
	buf := [20]byte{}
	pos := len(buf)
	for i > 0 {
		pos--
		buf[pos] = byte('0' + i%10)
		i /= 10
	}
	if neg {
		pos--
		buf[pos] = '-'
	}
	return string(buf[pos:])
}
