package helpers

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

// ResolveUsernameToMention attempts to convert a plain text username to a Discord Mention ping (<@UserID>).
// If the username is already formatted as a ping (contains <@), it simply returns it.
func ResolveUsernameToMention(s *discordgo.Session, guildID string, username string) string {
	if username == "" {
		return ""
	}

	// If it already looks like a valid discord ping, return it.
	if strings.HasPrefix(strings.TrimSpace(username), "<@") {
		return username
	}

	// Dynamic lookup through Discord servers members
	// Using Query argument which matches both username and nickname prefixes.
	members, err := s.GuildMembersSearch(guildID, username, 100)
	if err != nil {
		logrus.WithError(err).Warnf("Error searching members for %s in guild %s", username, guildID)
		return fmt.Sprintf("**%s**", username) // string fallback
	}

	for _, member := range members {
		// Exact or case-insensitive match on Username or GlobalName
		if strings.EqualFold(member.User.Username, username) || strings.EqualFold(member.User.GlobalName, username) || strings.EqualFold(member.Nick, username) {
			return fmt.Sprintf("<@%s>", member.User.ID)
		}
	}

	logrus.Infof("Could not resolve username %s to a member in guild %s", username, guildID)
	return fmt.Sprintf("**%s**", username) // string fallback
}
