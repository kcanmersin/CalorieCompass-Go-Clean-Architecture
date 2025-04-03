package service

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	baseURL            = "https://platform.fatsecret.com/rest/server.api"
	oauthTokenEndpoint = "https://oauth.fatsecret.com/connect/token"
	defaultTimeout     = 10 * time.Second
	methodFoodSearch   = "foods.search"
	methodFoodGet      = "food.get"
	methodRecipeSearch = "recipes.search"
	methodRecipeGet    = "recipe.get"
	responseFormatJSON = "json"
)

type FatSecretService struct {
	clientID       string
	clientSecret   string
	consumerKey    string
	consumerSecret string
	httpClient     *http.Client
	accessToken    string
	tokenExpiry    time.Time
}

type FatSecretResponse struct {
	Foods struct {
		Food []struct {
			FoodID          string `json:"food_id"`
			FoodName        string `json:"food_name"`
			BrandName       string `json:"brand_name,omitempty"`
			FoodType        string `json:"food_type"`
			FoodURL         string `json:"food_url"`
			FoodDescription struct {
				Calories string `json:"calories"`
				Carbs    string `json:"carbohydrate"`
				Protein  string `json:"protein"`
				Fat      string `json:"fat"`
			} `json:"food_description,omitempty"`
		} `json:"food"`
		MaxResults   string `json:"max_results"`
		PageNumber   string `json:"page_number"`
		TotalResults string `json:"total_results"`
	} `json:"foods"`
}

type FoodDetails struct {
	Food struct {
		FoodID       string `json:"food_id"`
		FoodName     string `json:"food_name"`
		BrandName    string `json:"brand_name,omitempty"`
		ServingSizes struct {
			ServingSize []struct {
				ServingID              string `json:"serving_id"`
				ServingDescription     string `json:"serving_description"`
				ServingURL             string `json:"serving_url"`
				MetricServingAmount    string `json:"metric_serving_amount"`
				MetricServingUnit      string `json:"metric_serving_unit"`
				NumberOfUnits          string `json:"number_of_units"`
				MeasurementDescription string `json:"measurement_description"`
				Calories               string `json:"calories"`
				Carbohydrate           string `json:"carbohydrate"`
				Protein                string `json:"protein"`
				Fat                    string `json:"fat"`
				SaturatedFat           string `json:"saturated_fat"`
				Fiber                  string `json:"fiber"`
				Cholesterol            string `json:"cholesterol"`
				Sodium                 string `json:"sodium"`
				Sugar                  string `json:"sugar"`
			} `json:"serving_size"`
		} `json:"servings"`
	} `json:"food"`
}

// NewFatSecretService creates a new FatSecret API service
func NewFatSecretService(clientID, clientSecret, consumerKey, consumerSecret string) *FatSecretService {
	return &FatSecretService{
		clientID:       clientID,
		clientSecret:   clientSecret,
		consumerKey:    consumerKey,
		consumerSecret: consumerSecret,
		httpClient:     &http.Client{Timeout: defaultTimeout},
	}
}

// SearchFoods searches for foods by query string
func (s *FatSecretService) SearchFoods(query string, pageNumber, maxResults int) (*FatSecretResponse, error) {
	// Try to use OAuth 2.0 first
	err := s.ensureToken()
	if err != nil {
		// If OAuth 2.0 fails, try OAuth 1.0
		fmt.Printf("OAuth 2.0 failed: %s, trying OAuth 1.0\n", err)
		return s.searchFoodsWithOAuth1(query, pageNumber, maxResults)
	}

	// OAuth 2.0 was successful
	params := url.Values{}
	params.Add("method", methodFoodSearch)
	params.Add("search_expression", query)
	params.Add("page_number", fmt.Sprintf("%d", pageNumber))
	params.Add("max_results", fmt.Sprintf("%d", maxResults))
	params.Add("format", responseFormatJSON)

	resp, err := s.makeRequestWithOAuth2(params)
	if err != nil {
		return nil, err
	}

	var result FatSecretResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return &result, nil
}

// searchFoodsWithOAuth1 searches for foods using OAuth 1.0
func (s *FatSecretService) searchFoodsWithOAuth1(query string, pageNumber, maxResults int) (*FatSecretResponse, error) {
	params := map[string]string{
		"method":            methodFoodSearch,
		"search_expression": query,
		"page_number":       fmt.Sprintf("%d", pageNumber),
		"max_results":       fmt.Sprintf("%d", maxResults),
		"format":            responseFormatJSON,
	}

	resp, err := s.makeRequestWithOAuth1(params)
	if err != nil {
		return nil, err
	}

	var result FatSecretResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return &result, nil
}

// GetFoodDetails gets detailed information about a specific food
func (s *FatSecretService) GetFoodDetails(foodID string) (*FoodDetails, error) {
	// Try to use OAuth 2.0 first
	err := s.ensureToken()
	if err != nil {
		// If OAuth 2.0 fails, try OAuth 1.0
		fmt.Printf("OAuth 2.0 failed: %s, trying OAuth 1.0\n", err)
		return s.getFoodDetailsWithOAuth1(foodID)
	}

	// OAuth 2.0 was successful
	params := url.Values{}
	params.Add("method", methodFoodGet)
	params.Add("food_id", foodID)
	params.Add("format", responseFormatJSON)

	resp, err := s.makeRequestWithOAuth2(params)
	if err != nil {
		return nil, err
	}

	var result FoodDetails
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return &result, nil
}

// getFoodDetailsWithOAuth1 gets detailed information about a specific food using OAuth 1.0
func (s *FatSecretService) getFoodDetailsWithOAuth1(foodID string) (*FoodDetails, error) {
	params := map[string]string{
		"method":  methodFoodGet,
		"food_id": foodID,
		"format":  responseFormatJSON,
	}

	resp, err := s.makeRequestWithOAuth1(params)
	if err != nil {
		return nil, err
	}

	var result FoodDetails
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return &result, nil
}

// ensureToken makes sure we have a valid access token
func (s *FatSecretService) ensureToken() error {
	// If token is still valid, return
	if s.accessToken != "" && time.Now().Before(s.tokenExpiry) {
		return nil
	}

	// Get a new token
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("scope", "basic")

	req, err := http.NewRequest("POST", oauthTokenEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("error creating token request: %w", err)
	}

	req.SetBasicAuth(s.clientID, s.clientSecret)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error getting token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error getting token: %d %s - %s", resp.StatusCode, resp.Status, string(body))
	}

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return fmt.Errorf("error decoding token response: %w", err)
	}

	s.accessToken = tokenResponse.AccessToken
	s.tokenExpiry = time.Now().Add(time.Duration(tokenResponse.ExpiresIn-60) * time.Second) // 60 seconds buffer

	return nil
}

// makeRequestWithOAuth2 makes an authenticated request using OAuth 2.0
func (s *FatSecretService) makeRequestWithOAuth2(params url.Values) ([]byte, error) {
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.URL.RawQuery = params.Encode()
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error from API: %s - %s", resp.Status, string(body))
	}

	var buffer bytes.Buffer
	if _, err := io.Copy(&buffer, resp.Body); err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	return buffer.Bytes(), nil
}

// makeRequestWithOAuth1 makes an authenticated request using OAuth 1.0
func (s *FatSecretService) makeRequestWithOAuth1(params map[string]string) ([]byte, error) {
	// Add OAuth parameters
	oauth := map[string]string{
		"oauth_consumer_key":     s.consumerKey,
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_timestamp":        strconv.FormatInt(time.Now().Unix(), 10),
		"oauth_nonce":            strconv.FormatInt(time.Now().UnixNano(), 10),
		"oauth_version":          "1.0",
	}

	// Combine params and oauth
	allParams := make(map[string]string)
	for k, v := range params {
		allParams[k] = v
	}
	for k, v := range oauth {
		allParams[k] = v
	}

	// Create signature base string
	baseString := s.createSignatureBaseString("GET", baseURL, allParams)

	// Calculate signature
	signature := s.calculateSignature(baseString, s.consumerSecret, "")

	// Add signature to OAuth params
	oauth["oauth_signature"] = signature

	// Create authorization header
	authHeader := "OAuth "
	i := 0
	for k, v := range oauth {
		if i > 0 {
			authHeader += ", "
		}
		authHeader += fmt.Sprintf("%s=\"%s\"", k, url.QueryEscape(v))
		i++
	}

	// Create URL with params
	reqURL := baseURL + "?"
	i = 0
	for k, v := range params {
		if i > 0 {
			reqURL += "&"
		}
		reqURL += fmt.Sprintf("%s=%s", k, url.QueryEscape(v))
		i++
	}

	// Make request
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Add("Authorization", authHeader)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error from API: %s - %s", resp.Status, string(body))
	}

	var buffer bytes.Buffer
	if _, err := io.Copy(&buffer, resp.Body); err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	return buffer.Bytes(), nil
}

// createSignatureBaseString creates the signature base string for OAuth 1.0
func (s *FatSecretService) createSignatureBaseString(method, baseURL string, params map[string]string) string {
	// Create a slice of parameter name-value pairs
	pairs := make([]string, 0, len(params))
	for k, v := range params {
		pairs = append(pairs, fmt.Sprintf("%s=%s", url.QueryEscape(k), url.QueryEscape(v)))
	}

	// Sort the pairs
	sort.Strings(pairs)

	// Join the pairs with &
	paramsStr := strings.Join(pairs, "&")

	// Create the signature base string
	signatureBaseString := fmt.Sprintf("%s&%s&%s",
		url.QueryEscape(method),
		url.QueryEscape(baseURL),
		url.QueryEscape(paramsStr))

	return signatureBaseString
}

// calculateSignature calculates the OAuth 1.0 signature
func (s *FatSecretService) calculateSignature(baseString, consumerSecret, tokenSecret string) string {
	key := fmt.Sprintf("%s&%s", url.QueryEscape(consumerSecret), url.QueryEscape(tokenSecret))

	// Create HMAC-SHA1 hash
	h := hmac.New(sha1.New, []byte(key))
	h.Write([]byte(baseString))

	// Base64 encode the hash
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return signature
}
