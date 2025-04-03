package usecase

import (
	"CalorieCompass/internal/entity"
	"context"
)

// FoodRepository defines the interface for food data operations
type FoodRepository interface {
	SearchFoods(ctx context.Context, query string, page, limit int) ([]entity.Food, int, error)
	GetFoodDetails(ctx context.Context, foodID string) (*entity.FoodDetails, error)
}

// FoodUseCase handles business logic for food operations
type FoodUseCase struct {
	repo FoodRepository
}

// NewFoodUseCase creates a new food use case
func NewFoodUseCase(repo FoodRepository) *FoodUseCase {
	return &FoodUseCase{
		repo: repo,
	}
}

// SearchFoods searches for foods by query
func (uc *FoodUseCase) SearchFoods(ctx context.Context, request entity.FoodSearchRequest) (entity.FoodSearchResponse, error) {
	foods, totalResults, err := uc.repo.SearchFoods(ctx, request.Query, request.Page, request.Limit)
	if err != nil {
		return entity.FoodSearchResponse{}, err
	}

	return entity.FoodSearchResponse{
		Foods:        foods,
		TotalResults: totalResults,
		Page:         request.Page,
		Limit:        request.Limit,
	}, nil
}

// GetFoodDetails gets detailed information about a food
func (uc *FoodUseCase) GetFoodDetails(ctx context.Context, foodID string) (*entity.FoodDetails, error) {
	return uc.repo.GetFoodDetails(ctx, foodID)
}
