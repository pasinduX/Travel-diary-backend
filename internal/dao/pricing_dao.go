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

var ErrPricingPlanNotFound = errors.New("pricing plan not found")

type PricingDAO struct {
	col *mongo.Collection
}

func NewPricingDAO(db *mongo.Database) PricingDAO {
	col := db.Collection("pricing_plans")
	_, _ = col.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "slug", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	return PricingDAO{col: col}
}

func (d PricingDAO) ListActive(ctx context.Context) ([]models.PricingPlan, error) {
	cur, err := d.col.Find(ctx, bson.M{"isActive": true}, options.Find().SetSort(bson.D{{Key: "sortOrder", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	plans := make([]models.PricingPlan, 0)
	for cur.Next(ctx) {
		var plan models.PricingPlan
		if err := cur.Decode(&plan); err != nil {
			return nil, err
		}
		plans = append(plans, plan)
	}
	return plans, cur.Err()
}

func (d PricingDAO) FindActiveBySlug(ctx context.Context, slug string) (models.PricingPlan, error) {
	var plan models.PricingPlan
	if err := d.col.FindOne(ctx, bson.M{"slug": slug, "isActive": true}).Decode(&plan); err != nil {
		return models.PricingPlan{}, ErrPricingPlanNotFound
	}
	return plan, nil
}

func (d PricingDAO) SeedDefaults(ctx context.Context, plans []models.PricingPlan) error {
	for _, plan := range plans {
		now := time.Now().UTC()
		_, err := d.col.UpdateOne(ctx, bson.M{"slug": plan.Slug}, bson.M{
			"$set": bson.M{
				"name": plan.Name, "description": plan.Description, "price": plan.Price,
				"currency": plan.Currency, "interval": plan.Interval, "features": plan.Features,
				"limits": plan.Limits, "isActive": plan.IsActive, "sortOrder": plan.SortOrder,
				"updatedAt": now,
			},
			"$setOnInsert": bson.M{"_id": plan.ID, "createdAt": now},
		}, options.Update().SetUpsert(true))
		if err != nil {
			return err
		}
	}
	return nil
}
