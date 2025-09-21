package utils

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func SendMessage(title, description string) *discordgo.MessageSend {
	return &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       title,
				Description: description,
				Color:       0xff76d2,
			},
		},
	}
}

func SetSuccessMessage(title, description string) *discordgo.MessageSend {
	return &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       fmt.Sprintf("✅ %s", title),
				Description: description,
				Color:       0x32ff7f,
			},
		},
	}
}

func SetErrorMessage(title, description string) *discordgo.MessageSend {
	return &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       fmt.Sprintf("❌ %s", title),
				Description: description,
				Color:       0xEA4108,
			},
		},
	}
}

func SetWarningMessage(title, description string) *discordgo.MessageSend {
	return &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       fmt.Sprintf("⚠️ %s", title),
				Description: description,
				Color:       0xFFD700,
			},
		},
	}
}

func SendEmbed(embed *discordgo.MessageSend, channelID string, session *discordgo.Session) {
	session.ChannelMessageSendComplex(channelID, embed)
}

func RoleSelectorEmbed(roles discordgo.Roles) *discordgo.MessageSend {

	return &discordgo.MessageSend{
		Embeds: func() []*discordgo.MessageEmbed {

			description := []string{"select role"}
			for index, role := range roles {
				description = append(description, fmt.Sprintf("%d. %s", index, role.Name))
			}
			finalDescription := strings.Join(description, "\n")

			var embeds []*discordgo.MessageEmbed
			embeds = append(embeds, &discordgo.MessageEmbed{
				Title:       "Choose Your Role",
				Description: finalDescription,
			})

			return embeds
		}(),
	}
}

