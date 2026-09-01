package service

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
	"travel-diary-backend/internal/ai"
	"travel-diary-backend/internal/config"
	"travel-diary-backend/internal/dao"
	"travel-diary-backend/internal/models"
)

const CurrentAnalysisVersion = 1

type analysisJob struct{ userID, imageID, tripID, imageURL string }
type ImageAnalysisQueue struct {
	jobs     chan analysisJob
	images   dao.TripImageDAO
	analyses dao.TripImageAnalysisDAO
	analyzer ai.ImageAnalyzer
	model    string
}

func NewImageAnalysisQueue(cfg config.Config, db *mongo.Database) *ImageAnalysisQueue {
	workers := cfg.ImageAnalysisWorkers
	if workers < 1 {
		workers = 4
	}
	q := &ImageAnalysisQueue{jobs: make(chan analysisJob, workers*4), images: dao.NewTripImageDAO(db), analyses: dao.NewTripImageAnalysisDAO(db), analyzer: ai.NewOpenAIImageAnalyzer(cfg.OpenAIKey), model: cfg.OpenAIModelImage}
	for i := 0; i < workers; i++ {
		go q.worker()
	}
	return q
}
func (q *ImageAnalysisQueue) Enqueue(ctx context.Context, image models.TripImage) error {
	if err := q.images.SetAnalysisState(ctx, image.ID, "QUEUED", "", nil); err != nil {
		return err
	}
	select {
	case q.jobs <- analysisJob{image.UserID, image.ID, image.TripID, image.S3URL}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
func (q *ImageAnalysisQueue) worker() {
	for job := range q.jobs {
		ctx := context.Background()
		_ = q.images.SetAnalysisState(ctx, job.imageID, "PROCESSING", "", nil)
		started := time.Now()
		result, err := q.analyzer.Analyze(ctx, job.imageURL, q.model)
		if err != nil {
			_ = q.images.SetAnalysisState(ctx, job.imageID, "FAILED", err.Error(), nil)
			log.Printf("image_analysis_failed image_id=%s duration_ms=%d error=%v", job.imageID, time.Since(started).Milliseconds(), err)
			continue
		}
		result.ID = job.imageID
		result.ImageID = job.imageID
		result.TripID = job.tripID
		result.UserID = job.userID
		result.Model = q.model
		result.AnalysisVersion = CurrentAnalysisVersion
		result.PromptVersion = "v1"
		if err := q.analyses.Upsert(ctx, result); err != nil {
			_ = q.images.SetAnalysisState(ctx, job.imageID, "FAILED", "could not save image analysis: "+err.Error(), nil)
			log.Printf("image_analysis_save_failed image_id=%s duration_ms=%d error=%v", job.imageID, time.Since(started).Milliseconds(), err)
			continue
		}
		now := time.Now().UTC()
		_ = q.images.SetAnalysisState(ctx, job.imageID, "ANALYZED", "", &now)
		log.Printf("image_analysis_completed image_id=%s duration_ms=%d", job.imageID, time.Since(started).Milliseconds())
	}
}
