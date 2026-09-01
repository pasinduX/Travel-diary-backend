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

var ErrTripNotFound = errors.New("trip not found")

type TripDAO struct {
	col *mongo.Collection
}

func NewTripDAO(db *mongo.Database) TripDAO {
	col := db.Collection("trips")
	_, _ = col.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "userId", Value: 1}, {Key: "createdAt", Value: -1}}},
	})
	return TripDAO{col: col}
}

func (d TripDAO) Create(ctx context.Context, trip models.Trip) (models.Trip, error) {
	now := time.Now().UTC()
	trip.CreatedAt = now
	trip.UpdatedAt = now
	_, err := d.col.InsertOne(ctx, trip)
	if err != nil {
		return models.Trip{}, err
	}
	return trip, nil
}

func (d TripDAO) ListByUserID(ctx context.Context, userID string) ([]models.Trip, error) {
	cur, err := d.col.Find(ctx, bson.M{"userId": userID}, options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var trips []models.Trip
	for cur.Next(ctx) {
		var trip models.Trip
		if err := cur.Decode(&trip); err != nil {
			return nil, err
		}
		trips = append(trips, trip)
	}
	return trips, cur.Err()
}

func (d TripDAO) FindByIDAndUserID(ctx context.Context, id, userID string) (models.Trip, error) {
	var trip models.Trip
	err := d.col.FindOne(ctx, bson.M{"_id": id, "userId": userID}).Decode(&trip)
	if err != nil {
		return models.Trip{}, ErrTripNotFound
	}
	return trip, nil
}

func (d TripDAO) Update(ctx context.Context, id, userID string, update bson.M) (models.Trip, error) {
	now := time.Now().UTC()
	update["updatedAt"] = now
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var trip models.Trip
	err := d.col.FindOneAndUpdate(ctx, bson.M{"_id": id, "userId": userID}, bson.M{"$set": update}, opts).Decode(&trip)
	if err != nil {
		return models.Trip{}, ErrTripNotFound
	}
	return trip, nil
}

func (d TripDAO) Delete(ctx context.Context, id, userID string) error {
	res, err := d.col.DeleteOne(ctx, bson.M{"_id": id, "userId": userID})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return ErrTripNotFound
	}
	return nil
}
