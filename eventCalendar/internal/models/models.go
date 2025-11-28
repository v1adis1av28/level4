package models

const (
	StatusScheduled = "scheduled"
	StatusCompleted = "completed"
	StatusCancelled = "cancelled"
)

type Event struct {
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Date        string `json:"date"`
	Status      string `json:"status"`
}

type ErrorResponse struct {
	Code  int    `json:"code"`
	Error string `json:"message"`
}

type SuccessResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Result  any    `json:"result,omitempty"`
}

type EventModificationRequest struct {
	ID       int    `json:"id"`
	NEW_DATE string `json:"new_date,omitempty"`
	NEW_NAME string `json:"new_name,omitempty"`
}
