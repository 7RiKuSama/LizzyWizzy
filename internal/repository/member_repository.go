package repository

import (
	"context"

	"github.com/7RikuSama/liz.git/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MembershipRepository struct {
	Collection *mongo.Collection
}

func NewMembershipRepo(client *mongo.Client) *MembershipRepository {
	return &MembershipRepository{
		Collection: client.Database("LizzyWizzy").Collection("members"),
	}
}

func (m *MembershipRepository) Insert(ctx context.Context, member *models.Membership) error {
	_, err := m.Collection.InsertOne(ctx, member)

	if err != nil {
		return err
	}
	return nil
}

func (m *MembershipRepository) FindByID(ctx context.Context, userID string) (*models.Membership, error) {
	var member *models.Membership
	if err := m.Collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&member); err != nil {
		return nil, err
	}
	return member, nil
}

func (m *MembershipRepository) Find(ctx context.Context, filter bson.M, sort bson.D, limit int64) ([]*models.Membership, error) {
	opts := options.Find()
	if sort != nil {
		opts.SetSort(sort)
	}
	if limit > 0 {
		opts.SetLimit(limit)
	}
	cursor, err := m.Collection.Find(ctx, filter, opts)
	
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	var data []*models.Membership
	
	if err := cursor.All(ctx, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (m *MembershipRepository) Delete(ctx context.Context, userID string) error {
	_, err := m.Collection.DeleteOne(ctx, bson.M{"user_id": userID})
	
	if err != nil {
		return err
	}
	return nil
}

func (m *MembershipRepository) Update(ctx context.Context, userID string, new *models.Membership) error {
	_, err := m.Collection.UpdateOne(ctx, bson.M{"user_id": userID}, bson.M{"$set" : new})
	if err != nil {
		return err
	}
	return nil
}

func (m *MembershipRepository) Count(ctx context.Context, userID string) (int, error) {
	result, err := m.Collection.CountDocuments(ctx, bson.M{"user_id": userID})
	if err != nil {
		return 0, err
	}

	return int(result), nil
}
