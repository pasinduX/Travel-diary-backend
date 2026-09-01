package models

import "time"

type TripImage struct {
	ID            string    `bson:"_id,omitempty" json:"id"`
	TripID        string    `bson:"tripId" json:"tripId"`
	UserID        string    `bson:"userId" json:"userId"`
	FileName      string    `bson:"fileName" json:"fileName"`
	ContentType   string    `bson:"contentType" json:"contentType"`
	FileSizeBytes int64     `bson:"fileSizeBytes" json:"fileSizeBytes"`
	Width         int       `bson:"width" json:"width"`
	Height        int       `bson:"height" json:"height"`
	DimensionName string    `bson:"dimensionName" json:"dimensionName"`
	S3Key         string    `bson:"s3Key" json:"s3Key"`
	S3URL         string    `bson:"s3Url" json:"s3Url"`
	CreatedAt     time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time `bson:"updatedAt" json:"updatedAt"`
}
