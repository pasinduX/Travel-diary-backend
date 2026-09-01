package models

import "time"

type RefreshSession struct {
	ID        string     `bson:"_id,omitempty" json:"id"`
	UserID    string     `bson:"userId" json:"userId"`
	TokenID   string     `bson:"tokenId" json:"tokenId"`
	TokenHash string     `bson:"tokenHash" json:"-"`
	ExpiresAt time.Time  `bson:"expiresAt" json:"expiresAt"`
	RevokedAt *time.Time `bson:"revokedAt,omitempty" json:"revokedAt,omitempty"`
	CreatedAt time.Time  `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time  `bson:"updatedAt" json:"updatedAt"`
}
