package service

import (
	"context"
	"errors"
	"strings"

	"travel-diary-backend/internal/dao"
	"travel-diary-backend/internal/dto"
	"travel-diary-backend/internal/models"
)

type PricingService struct {
	pricing dao.PricingDAO
	users   dao.UserDAO
}

func NewPricingService(pricing dao.PricingDAO, users ...dao.UserDAO) *PricingService {
	service := &PricingService{pricing: pricing}
	if len(users) > 0 {
		service.users = users[0]
	}
	return service
}

func (s *PricingService) List(ctx context.Context) ([]dto.PricingPlanResponse, error) {
	plans, err := s.pricing.ListActive(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]dto.PricingPlanResponse, 0, len(plans))
	for _, plan := range plans {
		result = append(result, toPricingResponse(plan))
	}
	return result, nil
}

func (s *PricingService) Get(ctx context.Context, slug string) (dto.PricingPlanResponse, error) {
	plan, err := s.pricing.FindActiveBySlug(ctx, slug)
	if err != nil {
		return dto.PricingPlanResponse{}, err
	}
	return toPricingResponse(plan), nil
}

func (s *PricingService) GetForUser(ctx context.Context, userID string) (dto.PricingPlanResponse, error) {
	if s.users == (dao.UserDAO{}) {
		return dto.PricingPlanResponse{}, errors.New("user pricing is not configured")
	}
	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return dto.PricingPlanResponse{}, err
	}
	if user.PricingPlan == "" {
		user.PricingPlan = "free"
	}
	return s.Get(ctx, user.PricingPlan)
}

func (s *PricingService) ChangeForUser(ctx context.Context, userID, slug string) (dto.PricingPlanResponse, error) {
	if s.users == (dao.UserDAO{}) {
		return dto.PricingPlanResponse{}, errors.New("user pricing is not configured")
	}
	slug = strings.ToLower(strings.TrimSpace(slug))
	if slug == "" {
		return dto.PricingPlanResponse{}, errors.New("pricing plan is required")
	}
	plan, err := s.pricing.FindActiveBySlug(ctx, slug)
	if err != nil {
		return dto.PricingPlanResponse{}, err
	}
	if err := s.users.AssignPricingPlan(ctx, userID, plan.ID, plan.Slug); err != nil {
		return dto.PricingPlanResponse{}, err
	}
	return toPricingResponse(plan), nil
}

func DefaultPricingPlans() []models.PricingPlan {
	return []models.PricingPlan{
		{
			ID: "free", Slug: "free", Name: "Free", Description: "Start turning your travel memories into cinematic albums.",
			Price: 0, Currency: "USD", Interval: "month", Features: []string{"AI-curated travel albums", "Original photo archive", "Shareable album links"},
			Limits: models.PricingLimits{NumberOfTrips: 2, MaxImages: 200}, IsActive: true, SortOrder: 1,
		},
		{
			ID: "starter", Slug: "starter", Name: "Starter", Description: "More room for bigger journeys and growing archives.",
			Price: 9.99, Currency: "USD", Interval: "month", Features: []string{"Everything in Free", "More travel albums", "Expanded image archive"},
			Limits: models.PricingLimits{NumberOfTrips: 10, MaxImages: 1000}, IsActive: true, SortOrder: 2,
		},
	}
}

func toPricingResponse(plan models.PricingPlan) dto.PricingPlanResponse {
	return dto.PricingPlanResponse{
		ID: plan.ID, Slug: plan.Slug, Name: plan.Name, Description: plan.Description,
		Price: plan.Price, Currency: plan.Currency, Interval: plan.Interval, Features: plan.Features,
		Limits:   dto.PricingLimits{NumberOfTrips: plan.Limits.NumberOfTrips, MaxImages: plan.Limits.MaxImages},
		IsActive: plan.IsActive, SortOrder: plan.SortOrder,
	}
}
