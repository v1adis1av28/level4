package utils

import (
	"eventCalendar/internal/models"
	"fmt"
	"time"
)

func IsEventValid(event *models.Event) (bool, error) {
	if event.Status != models.StatusScheduled {
		return false, fmt.Errorf("wrong type of status, allowed only: scheduled, completed, canceld")
	}
	if len(event.Name) < 1 || len(event.Description) < 1 {
		return false, fmt.Errorf("name or description of event can`t be empty")
	}
	date, err := time.Parse("2006-01-02", event.Date)
	if err != nil {
		return false, fmt.Errorf("wrong format of date string")
	}

	if date.After(time.Now()) {
		return false, fmt.Errorf("you can`t create events in past")
	}

	return true, nil
}

func IsModRequestValid(req *models.EventModificationRequest) (bool, error) {
	if req.ID < 0 {
		return false, fmt.Errorf("id of event can`t be less than 1")
	}
	date, err := time.Parse("2006-01-02", req.NEW_DATE)
	if err != nil {
		return false, fmt.Errorf("wrong format of date string")
	}

	if date.After(time.Now()) {
		return false, fmt.Errorf("you can`t create events in past")
	}

	return true, nil
}
