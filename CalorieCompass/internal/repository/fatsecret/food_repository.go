package fatsecret

import (
	"CalorieCompass/internal/entity"
	"CalorieCompass/internal/service"
	"context"
	"strconv"
)

// FoodRepository implements the food data operations using the FatSecret API
type FoodRepository struct {
	service *service.FatSecretService
}

// NewFoodRepository creates a new FatSecret food repository
func NewFoodRepository(service *service.FatSecretService) *FoodRepository {
	return &FoodRepository{
		service: service,
	}
}

// SearchFoods searches for foods in the FatSecret API
func (r *FoodRepository) SearchFoods(ctx context.Context, query string, page, limit int) ([]entity.Food, int, error) {
	response, err := r.service.SearchFoods(query, page, limit)
	if err != nil {
		return nil, 0, err
	}

	foods := make([]entity.Food, 0)

	// Check if there are any foods in the response
	if response.Foods.Food == nil {
		totalResults, _ := strconv.Atoi(response.Foods.TotalResults)
		return foods, totalResults, nil
	}

	for _, food := range response.Foods.Food {
		calories, _ := strconv.ParseFloat(food.FoodDescription.Calories, 64)
		carbs, _ := strconv.ParseFloat(food.FoodDescription.Carbs, 64)
		protein, _ := strconv.ParseFloat(food.FoodDescription.Protein, 64)
		fat, _ := strconv.ParseFloat(food.FoodDescription.Fat, 64)

		foods = append(foods, entity.Food{
			ID:        food.FoodID,
			Name:      food.FoodName,
			BrandName: food.BrandName,
			Type:      food.FoodType,
			URL:       food.FoodURL,
			Calories:  calories,
			Carbs:     carbs,
			Protein:   protein,
			Fat:       fat,
		})
	}

	totalResults, _ := strconv.Atoi(response.Foods.TotalResults)
	return foods, totalResults, nil
}

// GetFoodDetails gets detailed information about a specific food
func (r *FoodRepository) GetFoodDetails(ctx context.Context, foodID string) (*entity.FoodDetails, error) {
	response, err := r.service.GetFoodDetails(foodID)
	if err != nil {
		return nil, err
	}

	food := response.Food

	// Parse the serving sizes
	servings := make([]entity.Serving, 0)
	if response.Food.ServingSizes.ServingSize != nil {
		for _, srv := range response.Food.ServingSizes.ServingSize {
			calories, _ := strconv.ParseFloat(srv.Calories, 64)
			carbs, _ := strconv.ParseFloat(srv.Carbohydrate, 64)
			protein, _ := strconv.ParseFloat(srv.Protein, 64)
			fat, _ := strconv.ParseFloat(srv.Fat, 64)
			saturatedFat, _ := strconv.ParseFloat(srv.SaturatedFat, 64)
			fiber, _ := strconv.ParseFloat(srv.Fiber, 64)
			cholesterol, _ := strconv.ParseFloat(srv.Cholesterol, 64)
			sodium, _ := strconv.ParseFloat(srv.Sodium, 64)
			sugar, _ := strconv.ParseFloat(srv.Sugar, 64)
			metricAmount, _ := strconv.ParseFloat(srv.MetricServingAmount, 64)
			numberOfUnits, _ := strconv.ParseFloat(srv.NumberOfUnits, 64)

			servings = append(servings, entity.Serving{
				ID:                     srv.ServingID,
				Description:            srv.ServingDescription,
				URL:                    srv.ServingURL,
				MetricServingAmount:    metricAmount,
				MetricServingUnit:      srv.MetricServingUnit,
				NumberOfUnits:          numberOfUnits,
				MeasurementDescription: srv.MeasurementDescription,
				Calories:               calories,
				Carbs:                  carbs,
				Protein:                protein,
				Fat:                    fat,
				SaturatedFat:           saturatedFat,
				Fiber:                  fiber,
				Cholesterol:            cholesterol,
				Sodium:                 sodium,
				Sugar:                  sugar,
			})
		}
	}

	return &entity.FoodDetails{
		ID:        food.FoodID,
		Name:      food.FoodName,
		BrandName: food.BrandName,
		Servings:  servings,
	}, nil
}
