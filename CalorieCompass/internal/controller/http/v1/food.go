package v1

import (
	"context"
	"net/http"
	"strconv"

	"CalorieCompass/internal/entity"
	"github.com/gin-gonic/gin"
)

// FoodUseCase defines the interface for food business logic
type FoodUseCase interface {
	SearchFoods(ctx context.Context, request entity.FoodSearchRequest) (entity.FoodSearchResponse, error)
	GetFoodDetails(ctx context.Context, foodID string) (*entity.FoodDetails, error)
}

// FoodController handles HTTP requests for food operations
type FoodController struct {
	foodUseCase FoodUseCase
}

// NewFoodController creates a new food controller
func NewFoodController(foodUseCase FoodUseCase) *FoodController {
	return &FoodController{
		foodUseCase: foodUseCase,
	}
}

// @Summary Search foods
// @Description Search for foods by name
// @Tags food
// @Accept json
// @Produce json
// @Param query query string true "Search query"
// @Param page query int false "Page number (0-based)" default(0)
// @Param limit query int false "Results per page" default(50)
// @Success 200 {object} entity.FoodSearchResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /food/search [get]
func (c *FoodController) SearchFoods(ctx *gin.Context) {
	query := ctx.Query("query")
	if query == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "query parameter is required"})
		return
	}

	// Get page and limit from query parameters, with defaults
	page := 0
	limit := 50

	if pageStr := ctx.Query("page"); pageStr != "" {
		if pageVal, err := strconv.Atoi(pageStr); err == nil {
			page = pageVal
		}
	}

	if limitStr := ctx.Query("limit"); limitStr != "" {
		if limitVal, err := strconv.Atoi(limitStr); err == nil {
			limit = limitVal
		}
	}

	request := entity.FoodSearchRequest{
		Query: query,
		Page:  page,
		Limit: limit,
	}

	response, err := c.foodUseCase.SearchFoods(ctx.Request.Context(), request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// @Summary Get food details
// @Description Get detailed information about a specific food
// @Tags food
// @Accept json
// @Produce json
// @Param food_id path string true "Food ID"
// @Success 200 {object} entity.FoodDetails
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /food/{food_id} [get]
func (c *FoodController) GetFoodDetails(ctx *gin.Context) {
	foodID := ctx.Param("food_id")
	if foodID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "food_id parameter is required"})
		return
	}

	details, err := c.foodUseCase.GetFoodDetails(ctx.Request.Context(), foodID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if details == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "food not found"})
		return
	}

	ctx.JSON(http.StatusOK, details)
}
