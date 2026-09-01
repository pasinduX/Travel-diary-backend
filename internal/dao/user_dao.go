package dao

import (
	"context"
	"errors"
	"time"

	"travel-diary-backend/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrUserNotFound = errors.New("user not found")

type UserDAO struct {
	col *mongo.Collection
}

func NewUserDAO(db *mongo.Database) UserDAO {
	col := db.Collection("users")
	_, _ = col.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "username", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "email", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "googleId", Value: 1}}, Options: options.Index().SetUnique(true).SetSparse(true)},
	})
	return UserDAO{col: col}
}

func (d UserDAO) Create(ctx context.Context, u models.User) (models.User, error) {
	now := time.Now().UTC()
	u.CreatedAt = now
	u.UpdatedAt = now

	_, err := d.col.InsertOne(ctx, u)
	if err != nil {
		return models.User{}, err
	}

	return u, nil
}

func (d UserDAO) FindByUsername(ctx context.Context, username string) (models.User, error) {
	return d.findOne(ctx, "username", username)
}

func (d UserDAO) FindByEmail(ctx context.Context, email string) (models.User, error) {
	return d.findOne(ctx, "email", email)
}

func (d UserDAO) FindByID(ctx context.Context, id string) (models.User, error) {
	return d.findOne(ctx, "_id", id)
}

func (d UserDAO) UpsertGoogleUser(ctx context.Context, u models.User) (models.User, error) {
	now := time.Now().UTC()
	u.UpdatedAt = now

	filter := bson.M{"googleId": u.GoogleID}
	update := bson.M{
		"$set": bson.M{
			"username":     u.Username,
			"email":        u.Email,
			"googleId":     u.GoogleID,
			"authProvider": u.AuthProvider,
			"name":         u.Name,
			"pictureUrl":   u.PictureURL,
			"updatedAt":    now,
		},
		"$setOnInsert": bson.M{
			"_id":       u.ID,
			"createdAt": now,
		},
	}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	var updated models.User
	if err := d.col.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updated); err != nil {
		return models.User{}, err
	}
	return updated, nil
}

func (d UserDAO) findOne(ctx context.Context, column, value string) (models.User, error) {
	var u models.User
	err := d.col.FindOne(ctx, bson.M{column: value}).Decode(&u)
	if err != nil {
		return models.User{}, ErrUserNotFound
	}
	return u, nil
}
