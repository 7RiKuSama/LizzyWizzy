package repository

import (
	"context"

	"github.com/7RikuSama/liz.git/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UsersRepository struct {
	Collection *mongo.Collection
}

func NewUsersRepo(client *mongo.Client) *UsersRepository {
	return &UsersRepository{
		Collection: client.Database("LizzyWizzy").Collection("users"),
	}
}

func (u *UsersRepository) Insert(ctx context.Context, user *models.User) error {
	_, err := u.Collection.InsertOne(ctx, user)

	if err != nil {
		return err
	}
	return nil
}

func (u *UsersRepository) FindByID(ctx context.Context, userID string) (*models.User, error) {
	var user *models.User
	if err := u.Collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&user); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UsersRepository) Find(ctx context.Context, filter bson.M, sort bson.D, limit int64) ([]*models.User, error) {
	opts := options.Find()
	if sort != nil {
		opts.SetSort(sort)
	}
	if limit > 0 {
		opts.SetLimit(limit)
	}

	result, err := u.Collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	var data []*models.User

	if err := result.Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

func (u *UsersRepository) Delete(ctx context.Context, userID string) error {
	_, err := u.Collection.DeleteOne(ctx, bson.M{"user_id": userID})

	if err != nil {
		return err
	}
	return nil
}

func (u *UsersRepository) Upadate(ctx context.Context, userID string, new *models.User) error {
	_, err := u.Collection.UpdateOne(ctx, bson.M{"user_id": userID}, new)
	if err != nil {
		return err
	}
	return nil
}
