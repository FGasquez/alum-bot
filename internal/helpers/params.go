package helpers

import "github.com/bwmarrin/discordgo"

func GetParams(options []*discordgo.ApplicationCommandInteractionDataOption) map[string]interface{} {
	params := make(map[string]interface{})
	for _, option := range options {
		switch option.Type {
		case discordgo.ApplicationCommandOptionSubCommand, discordgo.ApplicationCommandOptionSubCommandGroup:
			for k, v := range GetParams(option.Options) {
				params[k] = v
			}
		default:
			params[option.Name] = option.Value
		}
	}
	return params
}
