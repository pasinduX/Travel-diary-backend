package models

import "time"

type PricingPlan struct {
	ID          string        `bson:"_id,omitempty" json:"id"`
	Slug        string        `bson:"slug" json:"slug"`
	Name        string        `bson:"name" json:"name"`
	Description string        `bson:"description" json:"description"`
	Price       float64       `bson:"price" json:"price"`
	Currency    string        `bson:"currency" json:"currency"`
	Interval    string        `bson:"interval" json:"interval"`
	Features    []string      `bson:"features" json:"features"`
	Limits      PricingLimits `bson:"limits" json:"limits"`
	IsActive    bool          `bson:"isActive" json:"isActive"`
	SortOrder   int           `bson:"sortOrder" json:"sortOrder"`
	CreatedAt   time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time     `bson:"updatedAt" json:"updatedAt"`
}

type PricingLimits struct {
	NumberOfTrips int `bson:"numberOfTrips" json:"numberOfTrips"`
	MaxImages     int `bson:"maxImages" json:"maxImages"`
}
