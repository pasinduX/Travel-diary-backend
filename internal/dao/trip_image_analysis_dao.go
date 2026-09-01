package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
	"travel-diary-backend/internal/models"
)

type TripImageAnalysisDAO struct{ col *mongo.Collection }

func NewTripImageAnalysisDAO(db *mongo.Database) TripImageAnalysisDAO {
	col := db.Collection("trip_image_analysis")
	_, _ = col.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "imageId", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "userId", Value: 1}, {Key: "tripId", Value: 1}}},
	})
	return TripImageAnalysisDAO{col: col}
}
func (d TripImageAnalysisDAO) Upsert(ctx context.Context, a models.TripImageAnalysis) error {
	now := time.Now().UTC()
	a.UpdatedAt = now
	_, err := d.col.UpdateOne(ctx, bson.M{"imageId": a.ImageID}, bson.M{"$set": a, "$setOnInsert": bson.M{"createdAt": now}}, options.Update().SetUpsert(true))
	return err
}

func (d TripImageAnalysisDAO) DeleteByTripID(ctx context.Context, userID, tripID string) error {
	_, err := d.col.DeleteMany(ctx, bson.M{"userId": userID, "tripId": tripID})
	return err
}
