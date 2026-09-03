package models

import "time"

type AlbumPlanRecord struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	TripID    string    `bson:"tripId" json:"tripId"`
	UserID    string    `bson:"userId" json:"userId"`
	Plan      AlbumPlan `bson:"plan" json:"plan"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

type AlbumPlan struct {
	Title    string         `bson:"title" json:"title"`
	Subtitle string         `bson:"subtitle" json:"subtitle"`
	Tone     string         `bson:"tone" json:"tone"`
	Chapters []AlbumChapter `bson:"chapters" json:"chapters"`
	Quotes   []AlbumQuote   `bson:"quotes" json:"quotes"`
}

type AlbumQuote struct {
	From  string `bson:"from" json:"from"`
	To    string `bson:"to" json:"to"`
	Text  string `bson:"text" json:"text"`
	Order int    `bson:"order" json:"order"`
}

type AlbumChapter struct {
	ID          string       `bson:"id" json:"id"`
	Order       int          `bson:"order" json:"order"`
	Eyebrow     string       `bson:"eyebrow" json:"eyebrow"`
	Title       string       `bson:"title" json:"title"`
	Quote       string       `bson:"quote,omitempty" json:"quote,omitempty"`
	Description string       `bson:"description,omitempty" json:"description,omitempty"`
	Blocks      []AlbumBlock `bson:"blocks" json:"blocks"`
}

type AlbumBlock struct {
	Type         string   `bson:"type" json:"type"`
	ImageIDs     []string `bson:"imageIds,omitempty" json:"imageIds,omitempty"`
	TextPosition string   `bson:"textPosition,omitempty" json:"textPosition,omitempty"`
	Eyebrow      string   `bson:"eyebrow,omitempty" json:"eyebrow,omitempty"`
	Title        string   `bson:"title,omitempty" json:"title,omitempty"`
	Text         string   `bson:"text,omitempty" json:"text,omitempty"`
	Quote        string   `bson:"quote,omitempty" json:"quote,omitempty"`
	Description  string   `bson:"description,omitempty" json:"description,omitempty"`
	Caption      string   `bson:"caption,omitempty" json:"caption,omitempty"`
}
