package dto

type TripCreateRequest struct {
	Title         string `json:"title"`
	Destination   string `json:"destination"`
	Departure     string `json:"departure"`
	Return        string `json:"return"`
	CinematicMood string `json:"cinematicMood"`
	Intention     string `json:"intention,omitempty"`
}

type TripUpdateRequest struct {
	Title         string `json:"title"`
	Destination   string `json:"destination"`
	Departure     string `json:"departure"`
	Return        string `json:"return"`
	CinematicMood string `json:"cinematicMood"`
	Intention     string `json:"intention,omitempty"`
}

type TripResponse struct {
	UserID        string `json:"userId"`
	ID            string `json:"id"`
	Title         string `json:"title"`
	Destination   string `json:"destination"`
	Departure     string `json:"departure"`
	Return        string `json:"return"`
	CinematicMood string `json:"cinematicMood"`
	Intention     string `json:"intention,omitempty"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
}
