package helpers

import (
	"testing"
)

func TestResolveUsernameToMention(t *testing.T) {
	// Fast path: Empty username
	res := ResolveUsernameToMention(nil, "", "")
	if res != "" {
		t.Errorf("Expected empty string, got %s", res)
	}

	// Fast path: Already ping formatted
	resPing := ResolveUsernameToMention(nil, "", "<@123456>")
	if resPing != "<@123456>" {
		t.Errorf("Expected <@123456>, got %s", resPing)
	}

	// Fast path: with spaces
	resPingSpace := ResolveUsernameToMention(nil, "", " <@123> ")
	if resPingSpace != " <@123> " {
		t.Errorf("Expected ' <@123> ', got %s", resPingSpace)
	}
	
	// Complex mock test skipped because Session requires active connection logic
	// inside discordgo.Session which cannot be easily mocked without interfaces.
}
