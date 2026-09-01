package dto

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type HealthResponse struct {
	Status  string `json:"status"`
	Env     string `json:"env"`
	AppName string `json:"appName"`
}
