package models

import "time"

type User struct {
	ID           string    `bson:"_id,omitempty" json:"id"`
	Username     string    `bson:"username" json:"username"`
	Email        string    `bson:"email" json:"email"`
	PasswordHash string    `bson:"passwordHash" json:"-"`
	GoogleID     string    `bson:"googleId,omitempty" json:"googleId,omitempty"`
	AuthProvider string    `bson:"authProvider" json:"provider"`
	Name         string    `bson:"name" json:"name"`
	PictureURL   string    `bson:"pictureUrl" json:"pictureUrl,omitempty"`
	CreatedAt    time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time `bson:"updatedAt" json:"updatedAt"`
}
