package dto

type PricingPlanResponse struct {
	ID          string        `json:"id"`
	Slug        string        `json:"slug"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Price       float64       `json:"price"`
	Currency    string        `json:"currency"`
	Interval    string        `json:"interval"`
	Features    []string      `json:"features"`
	Limits      PricingLimits `json:"limits"`
	IsActive    bool          `json:"isActive"`
	SortOrder   int           `json:"sortOrder"`
}

type PricingLimits struct {
	NumberOfTrips int `json:"numberOfTrips"`
	MaxImages     int `json:"maxImages"`
}

type PricingPlanRequest struct {
	Slug string `json:"slug"`
}
