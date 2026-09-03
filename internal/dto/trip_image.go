package dto

import "travel-diary-backend/internal/models"

type TripImageResponse struct {
	ID              string            `json:"id"`
	TripID          string            `json:"tripId"`
	UserID          string            `json:"userId"`
	FileName        string            `json:"fileName"`
	ContentType     string            `json:"contentType"`
	FileSizeBytes   int64             `json:"fileSizeBytes"`
	Width           int               `json:"width"`
	Height          int               `json:"height"`
	DimensionName   string            `json:"dimensionName"`
	EXIF            *models.ImageEXIF `json:"exif,omitempty"`
	S3Key           string            `json:"s3Key"`
	S3URL           string            `json:"s3Url"`
	AnalysisStatus  string            `json:"analysisStatus"`
	AnalysisError   string            `json:"analysisError,omitempty"`
	AnalyzedAt      string            `json:"analyzedAt,omitempty"`
	AnalysisVersion int               `json:"analysisVersion"`
	CreatedAt       string            `json:"createdAt"`
	UpdatedAt       string            `json:"updatedAt"`
}

type TripImageUploadResponse struct {
	Uploaded []TripImageResponse `json:"uploaded"`
}
