package models

import "time"

type Trip struct {
	ID            string    `bson:"_id,omitempty" json:"id"`
	UserID        string    `bson:"userId" json:"userId"`
	Title         string    `bson:"title" json:"title"`
	Destination   string    `bson:"destination" json:"destination"`
	Departure     time.Time `bson:"departure" json:"departure"`
	Return        time.Time `bson:"return" json:"return"`
	CinematicMood string    `bson:"cinematicMood" json:"cinematicMood"`
	Intention     string    `bson:"intention,omitempty" json:"intention,omitempty"`
	CreatedAt     time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time `bson:"updatedAt" json:"updatedAt"`
}
