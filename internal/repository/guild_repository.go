package repository

import (
	"context"

	"github.com/7RikuSama/liz.git/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type GuildsRepository struct {
	collection *mongo.Collection
}

func NewGuildsRepo(client *mongo.Client) *GuildsRepository {
	return &GuildsRepository{
		collection: client.Database("LizzyWizzy").Collection("guilds"),
	}
}

func (g *GuildsRepository) Insert(ctx context.Context, guild *models.Guild) error {
	_, err := g.collection.InsertOne(ctx, guild)

	if err != nil {
		return err
	}
	return nil
}

func (g *GuildsRepository) FindByID(ctx context.Context, guildID string) (*models.Guild, error) {
	var guild *models.Guild
	if err := g.collection.FindOne(ctx, bson.M{"guild_id": guildID}).Decode(&guild); err != nil {
		return nil, err
	}
	return guild, nil
}

func (g *GuildsRepository) Find(ctx context.Context, filter bson.M) ([]*models.Guild, error) {
	result, err := g.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var data []*models.Guild
	
	if err := result.Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

func (g *GuildsRepository) Delete(ctx context.Context, guildID string) error {
	_, err := g.collection.DeleteOne(ctx, bson.M{"guild_id": guildID})
	
	if err != nil {
		return err
	}
	return nil
}

func (g *GuildsRepository) Upadate(ctx context.Context, guildID string, new *models.Guild) error {
	_, err := g.collection.UpdateOne(ctx, bson.M{"guild_id": guildID}, new)
	if err != nil {
		return err
	}
	return nil
}

func (g *GuildsRepository) Count(ctx context.Context, guildID string) (int, error) {
	result, err := g.collection.CountDocuments(ctx, bson.M{"guild_id": guildID})
	if err != nil {
		return 0, err
	}

	return int(result), nil
}
