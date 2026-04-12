package helpers

import (
	"reflect"
	"testing"

	"github.com/bwmarrin/discordgo"
)

func TestGetParams(t *testing.T) {
	options := []*discordgo.ApplicationCommandInteractionDataOption{
		{
			Name:  "arg1",
			Type:  discordgo.ApplicationCommandOptionString,
			Value: "value1",
		},
		{
			Name: "subcommand",
			Type: discordgo.ApplicationCommandOptionSubCommand,
			Options: []*discordgo.ApplicationCommandInteractionDataOption{
				{
					Name:  "arg2",
					Type:  discordgo.ApplicationCommandOptionInteger,
					Value: 42,
				},
			},
		},
	}

	result := GetParams(options)
	expected := map[string]interface{}{
		"arg1": "value1",
		"arg2": 42,
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %#v, got %#v", expected, result)
	}
}
