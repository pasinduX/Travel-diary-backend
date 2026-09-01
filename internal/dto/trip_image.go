package dto

type TripImageResponse struct {
	ID            string `json:"id"`
	TripID        string `json:"tripId"`
	UserID        string `json:"userId"`
	FileName      string `json:"fileName"`
	ContentType   string `json:"contentType"`
	FileSizeBytes int64  `json:"fileSizeBytes"`
	Width         int    `json:"width"`
	Height        int    `json:"height"`
	DimensionName string `json:"dimensionName"`
	S3Key         string `json:"s3Key"`
	S3URL         string `json:"s3Url"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
}

type TripImageUploadResponse struct {
	Uploaded []TripImageResponse `json:"uploaded"`
}
