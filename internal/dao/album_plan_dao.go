package dao

import (
	"context"
	"time"
	"travel-diary-backend/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AlbumPlanDAO struct{ col *mongo.Collection }

func NewAlbumPlanDAO(db *mongo.Database) AlbumPlanDAO {
	col := db.Collection("album_plans")
	_, _ = col.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "userId", Value: 1}, {Key: "tripId", Value: 1}}, Options: options.Index().SetUnique(true)},
	})
	return AlbumPlanDAO{col: col}
}

func (d AlbumPlanDAO) Upsert(ctx context.Context, record models.AlbumPlanRecord) (models.AlbumPlanRecord, error) {
	now := time.Now().UTC()
	update := bson.M{
		"$set": bson.M{
			"plan":      record.Plan,
			"updatedAt": now,
		},
		"$setOnInsert": bson.M{
			"_id":       record.ID,
			"tripId":    record.TripID,
			"userId":    record.UserID,
			"createdAt": now,
		},
	}
	var saved models.AlbumPlanRecord
	err := d.col.FindOneAndUpdate(
		ctx,
		bson.M{"userId": record.UserID, "tripId": record.TripID},
		update,
		options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After),
	).Decode(&saved)
	return saved, err
}

func (d AlbumPlanDAO) FindByTripID(ctx context.Context, userID, tripID string) (models.AlbumPlanRecord, error) {
	var record models.AlbumPlanRecord
	err := d.col.FindOne(ctx, bson.M{"userId": userID, "tripId": tripID}).Decode(&record)
	return record, err
}

func (d AlbumPlanDAO) DeleteByTripID(ctx context.Context, userID, tripID string) error {
	_, err := d.col.DeleteOne(ctx, bson.M{"userId": userID, "tripId": tripID})
	return err
}
