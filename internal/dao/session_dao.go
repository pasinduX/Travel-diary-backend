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

var ErrSessionNotFound = errors.New("session not found")

type SessionDAO struct {
	col *mongo.Collection
}

func NewSessionDAO(db *mongo.Database) SessionDAO {
	col := db.Collection("refresh_sessions")
	_, _ = col.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "tokenId", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "userId", Value: 1}}, Options: options.Index()},
		{Keys: bson.D{{Key: "expiresAt", Value: 1}}, Options: options.Index().SetExpireAfterSeconds(0)},
	})
	return SessionDAO{col: col}
}

func (d SessionDAO) Create(ctx context.Context, s models.RefreshSession) (models.RefreshSession, error) {
	now := time.Now().UTC()
	s.CreatedAt = now
	s.UpdatedAt = now
	_, err := d.col.InsertOne(ctx, s)
	if err != nil {
		return models.RefreshSession{}, err
	}
	return s, nil
}

func (d SessionDAO) FindActiveByTokenID(ctx context.Context, tokenID string) (models.RefreshSession, error) {
	var session models.RefreshSession
	err := d.col.FindOne(ctx, bson.M{"tokenId": tokenID, "revokedAt": bson.M{"$exists": false}}).Decode(&session)
	if err != nil {
		return models.RefreshSession{}, ErrSessionNotFound
	}
	return session, nil
}

func (d SessionDAO) RevokeByTokenID(ctx context.Context, tokenID string) error {
	now := time.Now().UTC()
	_, err := d.col.UpdateOne(ctx, bson.M{"tokenId": tokenID, "revokedAt": bson.M{"$exists": false}}, bson.M{"$set": bson.M{"revokedAt": now, "updatedAt": now}})
	return err
}
