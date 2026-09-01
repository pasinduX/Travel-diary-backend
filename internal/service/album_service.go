package service

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"travel-diary-backend/internal/dao"
	"travel-diary-backend/internal/models"

	"github.com/google/uuid"
)

type AlbumService struct {
	trips    dao.TripDAO
	images   dao.TripImageDAO
	analyses dao.TripImageAnalysisDAO
	plans    dao.AlbumPlanDAO
}

func NewAlbumService(trips dao.TripDAO, images dao.TripImageDAO, analyses dao.TripImageAnalysisDAO, plans dao.AlbumPlanDAO) *AlbumService {
	return &AlbumService{trips: trips, images: images, analyses: analyses, plans: plans}
}

func (s *AlbumService) Generate(ctx context.Context, userID, tripID string) (models.AlbumPlanRecord, error) {
	trip, err := s.trips.FindByIDAndUserID(ctx, tripID, userID)
	if err != nil {
		return models.AlbumPlanRecord{}, err
	}
	images, err := s.images.ListByTripID(ctx, userID, tripID)
	if err != nil {
		return models.AlbumPlanRecord{}, err
	}
	analyses, err := s.analyses.ListByTripID(ctx, userID, tripID)
	if err != nil {
		return models.AlbumPlanRecord{}, err
	}

	analysisByImage := make(map[string]models.TripImageAnalysis, len(analyses))
	for _, analysis := range analyses {
		analysisByImage[analysis.ImageID] = analysis
	}
	sort.SliceStable(images, func(i, j int) bool {
		left := analysisByImage[images[i].ID]
		right := analysisByImage[images[j].ID]
		return left.Story.Importance > right.Story.Importance
	})

	plan := buildAlbumPlan(trip.Title, trip.Destination, images, analysisByImage)
	return s.plans.Upsert(ctx, models.AlbumPlanRecord{
		ID:     uuid.NewString(),
		TripID: tripID,
		UserID: userID,
		Plan:   plan,
	})
}

func (s *AlbumService) Get(ctx context.Context, userID, tripID string) (models.AlbumPlanRecord, error) {
	return s.plans.FindByTripID(ctx, userID, tripID)
}

func buildAlbumPlan(title, destination string, images []models.TripImage, analyses map[string]models.TripImageAnalysis) models.AlbumPlan {
	ids := make([]string, 0, len(images))
	for _, image := range images {
		if image.S3URL != "" {
			ids = append(ids, image.ID)
		}
	}
	if title == "" {
		title = destination
	}
	if title == "" {
		title = "Your journey"
	}

	blocks := make([]models.AlbumBlock, 0, 3)
	if len(ids) > 0 {
		blocks = append(blocks, models.AlbumBlock{Type: "album_cover", ImageIDs: ids[:1], Title: title})
	}
	if destination != "" {
		blocks = append(blocks, models.AlbumBlock{Type: "story_text", Text: fmt.Sprintf("A collection of moments from %s.", destination)})
	}
	for _, image := range images {
		if analysis, ok := analyses[image.ID]; ok && strings.TrimSpace(analysis.Caption) != "" {
			blocks = append(blocks, models.AlbumBlock{Type: "image_caption", ImageIDs: []string{image.ID}, Caption: analysis.Caption})
		}
	}
	if len(ids) > 1 {
		blocks = append(blocks, models.AlbumBlock{Type: "editorial_grid", ImageIDs: ids[1:]})
	}

	return models.AlbumPlan{
		Title: title, Subtitle: "A collection of moments, kept together.", Tone: "cinematic",
		Chapters: []models.AlbumChapter{{ID: "journey", Order: 1, Eyebrow: "Your photographs", Title: "The journey so far", Blocks: blocks}},
	}
}
