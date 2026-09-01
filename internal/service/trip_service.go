package service

import (
	"context"
	"time"

	"travel-diary-backend/internal/dao"
	"travel-diary-backend/internal/dto"
	"travel-diary-backend/internal/models"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type TripService struct {
	trips dao.TripDAO
}

func NewTripService(trips dao.TripDAO) *TripService {
	return &TripService{trips: trips}
}

func (s *TripService) Create(ctx context.Context, userID string, req dto.TripCreateRequest) (dto.TripResponse, error) {
	trip, err := s.trips.Create(ctx, models.Trip{
		ID:            uuid.NewString(),
		UserID:        userID,
		Title:         req.Title,
		Destination:   req.Destination,
		Departure:     mustParseTime(req.Departure),
		Return:        mustParseTime(req.Return),
		CinematicMood: req.CinematicMood,
		Intention:     req.Intention,
	})
	if err != nil {
		return dto.TripResponse{}, err
	}
	return toTripResponse(trip), nil
}

func (s *TripService) List(ctx context.Context, userID string) ([]dto.TripResponse, error) {
	trips, err := s.trips.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]dto.TripResponse, 0, len(trips))
	for _, trip := range trips {
		out = append(out, toTripResponse(trip))
	}
	return out, nil
}

func (s *TripService) Get(ctx context.Context, userID, tripID string) (dto.TripResponse, error) {
	trip, err := s.trips.FindByIDAndUserID(ctx, tripID, userID)
	if err != nil {
		return dto.TripResponse{}, err
	}
	return toTripResponse(trip), nil
}

func (s *TripService) Update(ctx context.Context, userID, tripID string, req dto.TripUpdateRequest) (dto.TripResponse, error) {
	trip, err := s.trips.Update(ctx, tripID, userID, bson.M{
		"title":         req.Title,
		"destination":   req.Destination,
		"departure":     mustParseTime(req.Departure),
		"return":        mustParseTime(req.Return),
		"cinematicMood": req.CinematicMood,
		"intention":     req.Intention,
	})
	if err != nil {
		return dto.TripResponse{}, err
	}
	return toTripResponse(trip), nil
}

func (s *TripService) Delete(ctx context.Context, userID, tripID string) error {
	return s.trips.Delete(ctx, tripID, userID)
}

func mustParseTime(value string) time.Time {
	t, _ := time.Parse(time.RFC3339, value)
	return t
}

func toTripResponse(trip models.Trip) dto.TripResponse {
	return dto.TripResponse{
		UserID:        trip.UserID,
		ID:            trip.ID,
		Title:         trip.Title,
		Destination:   trip.Destination,
		Departure:     trip.Departure.Format(time.RFC3339),
		Return:        trip.Return.Format(time.RFC3339),
		CinematicMood: trip.CinematicMood,
		Intention:     trip.Intention,
		CreatedAt:     trip.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     trip.UpdatedAt.Format(time.RFC3339),
	}
}
