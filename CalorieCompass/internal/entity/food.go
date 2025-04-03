package entity

// Food represents a food item
type Food struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	BrandName string  `json:"brand_name,omitempty"`
	Type      string  `json:"type"`
	URL       string  `json:"url"`
	Calories  float64 `json:"calories"`
	Carbs     float64 `json:"carbs"`
	Protein   float64 `json:"protein"`
	Fat       float64 `json:"fat"`
}

// Serving represents a serving size for a food
type Serving struct {
	ID                     string  `json:"id"`
	Description            string  `json:"description"`
	URL                    string  `json:"url"`
	MetricServingAmount    float64 `json:"metric_serving_amount"`
	MetricServingUnit      string  `json:"metric_serving_unit"`
	NumberOfUnits          float64 `json:"number_of_units"`
	MeasurementDescription string  `json:"measurement_description"`
	Calories               float64 `json:"calories"`
	Carbs                  float64 `json:"carbs"`
	Protein                float64 `json:"protein"`
	Fat                    float64 `json:"fat"`
	SaturatedFat           float64 `json:"saturated_fat"`
	Fiber                  float64 `json:"fiber"`
	Cholesterol            float64 `json:"cholesterol"`
	Sodium                 float64 `json:"sodium"`
	Sugar                  float64 `json:"sugar"`
}

// FoodDetails contains detailed information about a food
type FoodDetails struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	BrandName string    `json:"brand_name,omitempty"`
	Servings  []Serving `json:"servings"`
}

// FoodSearchRequest represents a request to search for foods
type FoodSearchRequest struct {
	Query string `json:"query" binding:"required"`
	Page  int    `json:"page" default:"0"`
	Limit int    `json:"limit" default:"50"`
}

// FoodSearchResponse represents the response from a food search
type FoodSearchResponse struct {
	Foods        []Food `json:"foods"`
	TotalResults int    `json:"total_results"`
	Page         int    `json:"page"`
	Limit        int    `json:"limit"`
}

// FoodDetailsRequest represents a request to get food details
type FoodDetailsRequest struct {
	FoodID string `json:"food_id" binding:"required"`
}
