package dao

import (
	"context"
	"time"

	"travel-diary-backend/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TripImageDAO struct {
	col *mongo.Collection
}

func NewTripImageDAO(db *mongo.Database) TripImageDAO {
	col := db.Collection("trip_images")
	_, _ = col.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "tripId", Value: 1}, {Key: "createdAt", Value: -1}}},
		{Keys: bson.D{{Key: "userId", Value: 1}, {Key: "createdAt", Value: -1}}},
	})
	return TripImageDAO{col: col}
}

func (d TripImageDAO) Create(ctx context.Context, image models.TripImage) (models.TripImage, error) {
	now := time.Now().UTC()
	image.CreatedAt = now
	image.UpdatedAt = now
	_, err := d.col.InsertOne(ctx, image)
	if err != nil {
		return models.TripImage{}, err
	}
	return image, nil
}

func (d TripImageDAO) ListByTripID(ctx context.Context, userID, tripID string) ([]models.TripImage, error) {
	cur, err := d.col.Find(ctx, bson.M{"userId": userID, "tripId": tripID}, options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var images []models.TripImage
	for cur.Next(ctx) {
		var img models.TripImage
		if err := cur.Decode(&img); err != nil {
			return nil, err
		}
		images = append(images, img)
	}
	return images, cur.Err()
}
