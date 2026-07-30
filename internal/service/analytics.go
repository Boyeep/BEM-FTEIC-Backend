package service

import (
	"context"

	"repo-backend/internal/dto"
	"repo-backend/internal/repository"
	"repo-backend/pkg/apperr"
)

type AnalyticsService struct {
	repo repository.AnalyticsRepository
}

func NewAnalytics(repo repository.AnalyticsRepository) *AnalyticsService {
	return &AnalyticsService{repo: repo}
}

func (s *AnalyticsService) Track(ctx context.Context, input dto.TrackVisitor) error {
	if err := s.repo.Track(ctx, input.ID, input.Path, input.UserAgent); err != nil {
		return apperr.Internal(err)
	}
	return nil
}

func (s *AnalyticsService) Count(ctx context.Context) (int64, error) {
	count, err := s.repo.Count(ctx)
	if err != nil {
		return 0, apperr.Internal(err)
	}
	return count, nil
}
