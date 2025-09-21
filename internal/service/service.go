package service

import (
	"context"

	"github.com/7RikuSama/liz.git/internal/models"
	"github.com/7RikuSama/liz.git/internal/repository"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	base = 100
	growth = 50
)

type Services struct {
	Users       *repository.UsersRepository
	Memberships *repository.MembershipRepository
	Guilds      *repository.GuildsRepository
}

func NewServices(client *mongo.Client) *Services {
	return &Services{
		Users:       repository.NewUsersRepo(client),
		Memberships: repository.NewMembershipRepo(client),
		Guilds:      repository.NewGuildsRepo(client),
	}
}

func (s *Services) DeleteUserMembership(ctx context.Context, userID string) error {
	result, err := s.Memberships.Count(ctx, userID)
	if err != nil {
		return err
	}

	if err := s.Memberships.Delete(ctx, userID); err != nil {
		return err
	}

	if result <= 1 {
		if err := s.Users.Delete(ctx, userID); err != nil {
			return err
		}
	}

	return nil
}

func (s *Services) DeleteGuild(ctx context.Context, guildID string) error {

	_, err := s.Memberships.Collection.DeleteMany(ctx, bson.M{"guild_id": guildID})

	if err != nil {
		return err
	}

	if err := s.Guilds.Delete(ctx, guildID); err != nil {
		return err
	}

	return nil
}

func (s *Services) Leaderboard(ctx context.Context, guildID string) ([]*models.Membership, error) {
	result, err := s.Memberships.Find(ctx, bson.M{"guild_id": guildID}, bson.D{{Key: "income", Value: -1}}, 10)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Services) CreateMembership(ctx context.Context, m *discordgo.MessageCreate) error {

	user := models.NewUser(m.Message.Author.ID, m.Author.Bot)
	if user.IsBot {
		return nil
	}
	guild := models.NewGuild(m.Message.GuildID)
	member := models.NewMembership(m.Message.Author.ID, m.Message.GuildID, m.Message.Member.Nick)

	_, err := s.Users.FindByID(ctx, user.UserID)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			if err := s.Users.Insert(ctx, user); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	_, err = s.Guilds.FindByID(ctx, member.GuildID)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			if err := s.Guilds.Insert(ctx, guild); err != nil {
				return err
			}
		} else {
			return err
		}
	}


	result, err := s.Memberships.Find(ctx, bson.M{"guild_id": member.GuildID, "user_id": user.UserID}, nil, 0)

	if err != nil {
		return err
	}
	if len(result) == 0 {
		if err := s.Memberships.Insert(ctx, member); err != nil {
			return err
		}
	}
	return nil
}

