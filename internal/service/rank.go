package service

import (
	"context"
	"log"

	"github.com/7RikuSama/liz.git/internal/models"
	"github.com/7RikuSama/liz.git/internal/utils"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/mongo"
)


func levelUp(totalXP, base, growth, level int) (leveledUp bool, progress int, required int) {
    required = base + (level*level)*growth

    prevTotal := 0
    for i := 1; i < level; i++ {
        prevTotal += base + (i*i)*growth
    }

    progress = totalXP - prevTotal
    leveledUp = progress >= required
    return
}


type RankInfo struct {
	Member *models.Membership
	Progress int
	Required int
}

func (s *Services) GetRank(ctx context.Context, userID string) (*RankInfo, error) {
	
	member, err := s.Memberships.FindByID(ctx, userID)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		} else {
			return nil, err
		}
	}

	_, progress, required := levelUp(member.Exp, 100, growth, member.Level)

	return &RankInfo{
		Member: member,
		Progress: progress,
		Required: required,
	}, nil
}


func (s *Services) LoadRank(ctx context.Context, session *discordgo.Session, message *discordgo.MessageCreate) error {

	member, err := s.Memberships.FindByID(ctx, message.Author.ID)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		} else {
			return err
		}
	}

	_, progress, required := levelUp(member.Exp, 100, growth, member.Level)

	buf, err := utils.MemberCard(utils.NewColor(242, 140, 220), member.Level, progress, required, message.Author.ID, message.Author.Username, message.Author.Avatar)
	if err != nil {
		log.Println(err)
	}

	_, err = session.ChannelMessageSendComplex(message.ChannelID, &discordgo.MessageSend{
		Content: "<@" + message.Author.ID + ">",
		Files: []*discordgo.File{
			{
				Name:   "level.png",
				Reader: buf,
			},
		},
	})

	return nil
}




func (s *Services) LevelTracker(ctx context.Context, session *discordgo.Session, message *discordgo.MessageCreate) error {

	member, err := s.Memberships.FindByID(ctx, message.Author.ID)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		} else {
			return err
		}
	}

	levelUp, _, required := levelUp(member.Exp, base, growth, member.Level)

	if levelUp {
		member.Level += 1
		member.Exp += required - int(float64(required)*0.9)
		session.ChannelMessageSendComplex(message.ChannelID, utils.SetSuccessMessage("Level UP!", "You are reaching higher levels. Congrats!"))
	}

	if err := s.Memberships.Update(ctx, message.Author.ID, member); err != nil {
		return err
	}

	return nil
}

func (s *Services) IncrementExp(ctx context.Context, m *discordgo.MessageCreate) error {
	member, err := s.Memberships.FindByID(ctx, m.Author.ID)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		} else {
			return err
		}
	}

	_, _, required := levelUp(member.Exp, base, growth, member.Level)
	factor := 0.05 - (0.002 * float64(member.Level)) // reduces over time
	if factor < 0.01 {
		factor = 0.01
	}
	gain := int(float64(required) * factor)
	member.Exp += gain

	if err := s.Memberships.Update(ctx, m.Author.ID, member); err != nil {
		return err
	}

	return nil
}


