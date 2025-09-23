package commands

import (
	"context"
	"fmt"
	"log"

	"github.com/7RikuSama/liz.git/internal/service"
	"github.com/7RikuSama/liz.git/internal/utils"
	"github.com/bwmarrin/discordgo"
)

func LoadRank(ctx context.Context, session *discordgo.Session, services *service.Services, event any) error {

	var authorID, username, avatar string

	switch v := event.(type) {
	case *discordgo.MessageCreate:
		authorID = v.Author.ID
		username = v.Author.Username
		avatar = v.Author.Avatar
	case *discordgo.InteractionCreate:
		authorID = v.Member.User.ID
		username = v.Member.User.Username
		avatar = v.Member.User.Avatar
	}

	rank, err := services.GetRank(ctx, authorID)

	if err != nil {
		return err
	}

	buf, err := utils.MemberCard(utils.NewColor(242, 140, 220), rank.Member.Level, rank.Progress, rank.Required, authorID, username, avatar)
	if err != nil {
		log.Println("err: ", err)
	}

	switch v := event.(type) {
	case *discordgo.MessageCreate:
		_, err = session.ChannelMessageSendComplex(v.ChannelID, &discordgo.MessageSend{
			Content: "<@" + v.Author.ID + ">",
			Files: []*discordgo.File{
				{
					Name:   "level.png",
					Reader: buf,
				},
			},
		})
	case *discordgo.InteractionCreate:	
		strContent := fmt.Sprintf("<@%s>", v.Member.User.ID)
		session.InteractionResponseEdit(v.Interaction, &discordgo.WebhookEdit{
			Content: &strContent,
			Files: []*discordgo.File{
				{
					Name:   "level.png",
					Reader: buf,
				},
			},
		})
	}

	if err != nil {
		return err
	}

	return nil
}
