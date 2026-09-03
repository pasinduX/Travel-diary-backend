package service

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rwcarlsen/goexif/exif"
	"travel-diary-backend/internal/dao"
	"travel-diary-backend/internal/dto"
	"travel-diary-backend/internal/integrations"
	"travel-diary-backend/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type TripImageService struct {
	trips    dao.TripDAO
	images   dao.TripImageDAO
	analyses dao.TripImageAnalysisDAO
	s3       *integrations.S3Client
	analysis *ImageAnalysisQueue
}

func NewTripImageService(trips dao.TripDAO, images dao.TripImageDAO, analyses dao.TripImageAnalysisDAO, s3 *integrations.S3Client, analysis *ImageAnalysisQueue) *TripImageService {
	return &TripImageService{trips: trips, images: images, analyses: analyses, s3: s3, analysis: analysis}
}

func (s *TripImageService) UploadMany(ctx context.Context, userID, tripID string, files []*multipart.FileHeader) ([]dto.TripImageResponse, error) {
	if _, err := s.trips.FindByIDAndUserID(ctx, tripID, userID); err != nil {
		return nil, fiber.ErrNotFound
	}

	out := make([]dto.TripImageResponse, 0, len(files))
	for _, fh := range files {
		resp, err := s.uploadOne(ctx, userID, tripID, fh)
		if err != nil {
			return nil, err
		}
		out = append(out, resp)
	}
	return out, nil
}

func (s *TripImageService) List(ctx context.Context, userID, tripID string) ([]dto.TripImageResponse, error) {
	images, err := s.images.ListByTripID(ctx, userID, tripID)
	if err != nil {
		return nil, err
	}
	out := make([]dto.TripImageResponse, 0, len(images))
	for _, img := range images {
		out = append(out, toTripImageResponse(img))
	}
	return out, nil
}

// RetryFailedAnalysis requeues only images that previously failed analysis.
func (s *TripImageService) RetryFailedAnalysis(ctx context.Context, userID, tripID string) (int, error) {
	if _, err := s.trips.FindByIDAndUserID(ctx, tripID, userID); err != nil {
		return 0, fiber.ErrNotFound
	}
	images, err := s.images.ListByTripID(ctx, userID, tripID)
	if err != nil {
		return 0, err
	}
	if s.analysis == nil {
		return 0, fmt.Errorf("image analysis is unavailable")
	}

	retried := 0
	for _, image := range images {
		if image.AnalysisStatus != "FAILED" {
			continue
		}
		if err := s.analysis.Enqueue(ctx, image); err != nil {
			return retried, err
		}
		retried++
	}
	return retried, nil
}

func (s *TripImageService) DeleteByTripID(ctx context.Context, userID, tripID string) error {
	if _, err := s.trips.FindByIDAndUserID(ctx, tripID, userID); err != nil {
		return fiber.ErrNotFound
	}
	images, err := s.images.ListByTripID(ctx, userID, tripID)
	if err != nil {
		return err
	}
	for _, image := range images {
		if err := s.s3.DeleteObject(ctx, image.S3Key); err != nil {
			return err
		}
	}
	if err := s.analyses.DeleteByTripID(ctx, userID, tripID); err != nil {
		return err
	}
	return s.images.DeleteByTripID(ctx, userID, tripID)
}

func (s *TripImageService) Delete(ctx context.Context, userID, tripID, imageID string) error {
	if _, err := s.trips.FindByIDAndUserID(ctx, tripID, userID); err != nil {
		return fiber.ErrNotFound
	}
	image, err := s.images.FindByID(ctx, userID, imageID)
	if err != nil || image.TripID != tripID {
		return fiber.ErrNotFound
	}
	if err := s.s3.DeleteObject(ctx, image.S3Key); err != nil {
		return err
	}
	if err := s.analyses.DeleteByImageID(ctx, userID, imageID); err != nil {
		return err
	}
	return s.images.DeleteByID(ctx, userID, tripID, imageID)
}

func (s *TripImageService) uploadOne(ctx context.Context, userID, tripID string, fh *multipart.FileHeader) (dto.TripImageResponse, error) {
	file, err := fh.Open()
	if err != nil {
		return dto.TripImageResponse{}, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return dto.TripImageResponse{}, err
	}

	cfg, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return dto.TripImageResponse{}, err
	}
	exifData := extractEXIF(data)

	key := fmt.Sprintf("trips/%s/%s%s", tripID, uuid.NewString(), filepath.Ext(fh.Filename))
	contentType := fh.Header.Get("Content-Type")
	if contentType == "" {
		contentType = httpDetectContentType(data)
	}

	_, err = s.s3.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      &s.s3.Bucket,
		Key:         &key,
		Body:        bytes.NewReader(data),
		ContentType: &contentType,
	})
	if err != nil {
		return dto.TripImageResponse{}, err
	}

	record, err := s.images.Create(ctx, models.TripImage{
		ID:             uuid.NewString(),
		TripID:         tripID,
		UserID:         userID,
		FileName:       fh.Filename,
		ContentType:    contentType,
		FileSizeBytes:  int64(len(data)),
		Width:          cfg.Width,
		Height:         cfg.Height,
		DimensionName:  fmt.Sprintf("%dx%d", cfg.Width, cfg.Height),
		EXIF:           exifData,
		S3Key:          key,
		S3URL:          s.s3.PublicURL(key),
		AnalysisStatus: "UPLOADED",
	})
	if err != nil {
		return dto.TripImageResponse{}, err
	}
	if s.analysis != nil {
		if err := s.analysis.Enqueue(ctx, record); err != nil {
			return dto.TripImageResponse{}, err
		}
		record.AnalysisStatus = "QUEUED"
	}
	return toTripImageResponse(record), nil
}

func httpDetectContentType(data []byte) string {
	return http.DetectContentType(data)
}

func toTripImageResponse(img models.TripImage) dto.TripImageResponse {
	return dto.TripImageResponse{
		ID:              img.ID,
		TripID:          img.TripID,
		UserID:          img.UserID,
		FileName:        img.FileName,
		ContentType:     img.ContentType,
		FileSizeBytes:   img.FileSizeBytes,
		Width:           img.Width,
		Height:          img.Height,
		DimensionName:   img.DimensionName,
		EXIF:            img.EXIF,
		S3Key:           img.S3Key,
		S3URL:           img.S3URL,
		AnalysisStatus:  img.AnalysisStatus,
		AnalysisError:   img.AnalysisError,
		AnalysisVersion: img.AnalysisVersion,
		CreatedAt:       img.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       img.UpdatedAt.Format(time.RFC3339),
	}
}

func extractEXIF(data []byte) *models.ImageEXIF {
	metadata, err := exif.Decode(bytes.NewReader(data))
	if err != nil {
		return nil
	}

	result := &models.ImageEXIF{
		CapturedAt:   normalizeEXIFDate(exifString(metadata, exif.DateTimeOriginal)),
		CameraMake:   exifString(metadata, exif.Make),
		CameraModel:  exifString(metadata, exif.Model),
		LensModel:    exifString(metadata, exif.LensModel),
		ShutterSpeed: exifString(metadata, exif.ExposureTime),
	}
	result.ISO = exifInt(metadata, exif.ISOSpeedRatings)
	result.Aperture = exifFloat(metadata, exif.FNumber)
	result.FocalLength = exifFloat(metadata, exif.FocalLength)

	if *result == (models.ImageEXIF{}) {
		return nil
	}
	return result
}

func normalizeEXIFDate(value string) string {
	parsed, err := time.Parse("2006:01:02 15:04:05", value)
	if err != nil {
		return value
	}
	return parsed.Format("2006-01-02T15:04:05")
}

func exifString(metadata *exif.Exif, tag exif.FieldName) string {
	value, err := metadata.Get(tag)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(value.String())
}

func exifInt(metadata *exif.Exif, tag exif.FieldName) int {
	value := exifString(metadata, tag)
	parsed, _ := strconv.Atoi(value)
	return parsed
}

func exifFloat(metadata *exif.Exif, tag exif.FieldName) float64 {
	value := exifString(metadata, tag)
	if strings.Contains(value, "/") {
		parts := strings.SplitN(value, "/", 2)
		numerator, numErr := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
		denominator, denErr := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
		if numErr == nil && denErr == nil && denominator != 0 {
			return numerator / denominator
		}
	}
	parsed, _ := strconv.ParseFloat(value, 64)
	return parsed
}
