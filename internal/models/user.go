package models

import "time"

type User struct {
	ID            string    `bson:"_id,omitempty" json:"id"`
	Username      string    `bson:"username" json:"username"`
	Email         string    `bson:"email" json:"email"`
	PasswordHash  string    `bson:"passwordHash" json:"-"`
	GoogleID      string    `bson:"googleId,omitempty" json:"googleId,omitempty"`
	Auth0ID       string    `bson:"auth0Id,omitempty" json:"auth0Id,omitempty"`
	AuthProvider  string    `bson:"authProvider" json:"provider"`
	Name          string    `bson:"name" json:"name"`
	PictureURL    string    `bson:"pictureUrl" json:"pictureUrl,omitempty"`
	PricingPlanID string    `bson:"pricingPlanId" json:"pricingPlanId"`
	PricingPlan   string    `bson:"pricingPlan" json:"pricingPlan"`
	CreatedAt     time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time `bson:"updatedAt" json:"updatedAt"`
}
