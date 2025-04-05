package models

type VM struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Provider string `json:"provider"`
	Region   string `json:"region"`
	Status   string `json:"status"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}
