package utils

import (
	"eventCalendar/internal/models"
	"testing"
	"time"
)

func TestIsEventValid(t *testing.T) {
	now := time.Now()
	futureDate := now.Add(24 * time.Hour).Format("2006-01-02")
	pastDate := now.Add(-24 * time.Hour).Format("2006-01-02")

	tests := []struct {
		name        string
		event       *models.Event
		expectValid bool
		expectError bool
	}{
		{
			name: "Valid event",
			event: &models.Event{
				Name:             "Test Event",
				Description:      "A valid event",
				Date:             futureDate,
				Status:           models.StatusScheduled,
				HaveNotification: false,
			},
			expectValid: true,
			expectError: false,
		},
		{
			name: "Invalid status",
			event: &models.Event{
				Name:             "Test Event",
				Description:      "A valid event",
				Date:             futureDate,
				Status:           "wrong status",
				HaveNotification: false,
			},
			expectValid: false,
			expectError: true,
		},
		{
			name: "Empty name",
			event: &models.Event{
				Name:             "",
				Description:      "A valid event",
				Date:             futureDate,
				Status:           models.StatusScheduled,
				HaveNotification: false,
			},
			expectValid: false,
			expectError: true,
		},
		{
			name: "Empty description",
			event: &models.Event{
				Name:             "Test Event",
				Description:      "",
				Date:             futureDate,
				Status:           models.StatusScheduled,
				HaveNotification: false,
			},
			expectValid: false,
			expectError: true,
		},
		{
			name: "Invalid date format",
			event: &models.Event{
				Name:             "Test Event",
				Description:      "A valid event",
				Date:             "invalid-date",
				Status:           models.StatusScheduled,
				HaveNotification: false,
			},
			expectValid: false,
			expectError: true,
		},
		{
			name: "Event in the past",
			event: &models.Event{
				Name:             "Test Event",
				Description:      "A valid event",
				Date:             pastDate,
				Status:           models.StatusScheduled,
				HaveNotification: false,
			},
			expectValid: false,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, err := IsEventValid(tt.event)

			if (err != nil) != tt.expectError {
				t.Errorf("IsEventValid() error = %v, expectError %v", err, tt.expectError)
				return
			}
			if valid != tt.expectValid {
				t.Errorf("IsEventValid() = %v, expectValid %v", valid, tt.expectValid)
			}
		})
	}
}
