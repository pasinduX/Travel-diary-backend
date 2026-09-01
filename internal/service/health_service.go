package service

import (
	"context"

	"travel-diary-backend/internal/config"
	"travel-diary-backend/internal/dao"
	"travel-diary-backend/internal/dto"
)

type HealthService struct {
	cfg config.Config
	dao dao.HealthDAO
}

func NewHealthService(cfg config.Config) *HealthService {
	return &HealthService{
		cfg: cfg,
		dao: dao.NewHealthDAO(),
	}
}

func (s *HealthService) Check(ctx context.Context) dto.HealthResponse {
	_ = ctx

	return dto.HealthResponse{
		Status:  "healthy",
		Env:     s.cfg.Environment,
		AppName: s.cfg.AppName,
	}
}
