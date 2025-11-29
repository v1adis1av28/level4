package handlers

import (
	"eventCalendar/internal/models"
	"eventCalendar/internal/storage"
	"eventCalendar/internal/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateEventHandler(storage *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newEvent models.Event
		if err := c.ShouldBindJSON(&newEvent); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Code:  http.StatusBadRequest,
				Error: "Invalid request payload",
			})
			return
		}
		if newEvent.HaveNotification {
			//todo добавляем в фоновый ворекер

		}
		isEventValid, err := utils.IsEventValid(&newEvent)
		if !isEventValid && err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Code:  http.StatusBadRequest,
				Error: err.Error(),
			})
			return
		}
		err = storage.CreateEvent(&newEvent)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Code:  http.StatusInternalServerError,
				Error: "Failed to create event",
			})
			return
		}

		c.JSON(http.StatusOK, models.SuccessResponse{
			Code:    http.StatusOK,
			Message: "Event created successfully",
			Result:  newEvent,
		})
	}
}

func UpdateEventHandler(s *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.EventModificationRequest
		err := c.ShouldBindJSON(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Code: http.StatusBadRequest, Error: "Error on binding json"})
			return
		}

		isValidModRequest, err := utils.IsModRequestValid(&req)
		if !isValidModRequest && err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Code:  http.StatusBadRequest,
				Error: err.Error(),
			})
			return
		}

		err = s.UpdateEvent(&req)
		if err != nil {
			c.JSON(500, models.ErrorResponse{Code: 500, Error: err.Error()})
			return
		}

		c.JSON(200, models.SuccessResponse{Code: 200, Message: "succesfully update event", Result: req})
	}
}
func DeleteEventHandler(s *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			ID int `json:"id" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Code:  http.StatusBadRequest,
				Error: "Invalid request payload or missing 'id'",
			})
			return
		}

		if req.ID <= 0 {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Code:  http.StatusBadRequest,
				Error: "Event ID must be greater than 0",
			})
			return
		}

		err := s.DeleteEvent(req.ID)
		if err != nil {
			if err.Error() == fmt.Sprintf("event with id %d does not exist", req.ID) {
				c.JSON(http.StatusNotFound, models.ErrorResponse{
					Code:  http.StatusNotFound,
					Error: "Event not found",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Code:  http.StatusInternalServerError,
				Error: "Failed to delete event",
			})
			return
		}

		c.JSON(http.StatusOK, models.SuccessResponse{
			Code:    http.StatusOK,
			Message: "Event deleted successfully",
			Result:  fmt.Sprintf("deleted id :%v", req.ID),
		})
	}
}

func GetEventsForDayHandler(s *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		dateStr := c.Query("date")
		if dateStr == "" {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Code:  http.StatusBadRequest,
				Error: "Query parameter 'date' is required (format: YYYY-MM-DD)",
			})
			return
		}

		_, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Code:  http.StatusBadRequest,
				Error: "Invalid date format, expected YYYY-MM-DD",
			})
			return
		}

		events, err := s.GetEventsForDay(dateStr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Code:  http.StatusInternalServerError,
				Error: "Failed to retrieve events for day",
			})
			return
		}

		c.JSON(http.StatusOK, models.SuccessResponse{
			Code:    http.StatusOK,
			Message: "Events retrieved successfully",
			Result:  events,
		})
	}
}

func GetEventsForWeekHandler(s *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		dateStr := c.Query("date")
		if dateStr == "" {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Code:  http.StatusBadRequest,
				Error: "Query parameter 'date' is required (format: YYYY-MM-DD) to determine the week",
			})
			return
		}

		_, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Code:  http.StatusBadRequest,
				Error: "Invalid date format, expected YYYY-MM-DD",
			})
			return
		}

		events, err := s.GetEventsForWeek(dateStr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Code:  http.StatusInternalServerError,
				Error: "Failed to retrieve events for week",
			})
			return
		}

		c.JSON(http.StatusOK, models.SuccessResponse{
			Code:    http.StatusOK,
			Message: "Events retrieved successfully",
			Result:  events,
		})
	}
}

func GetEventsForMonthHandler(s *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		dateStr := c.Query("date")
		if dateStr == "" {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Code:  http.StatusBadRequest,
				Error: "Query parameter 'date' is required (format: YYYY-MM-DD) to determine the month",
			})
			return
		}

		_, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Code:  http.StatusBadRequest,
				Error: "Invalid date format, expected YYYY-MM-DD",
			})
			return
		}

		events, err := s.GetEventsForMonth(dateStr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Code:  http.StatusInternalServerError,
				Error: "Failed to retrieve events for month",
			})
			return
		}

		c.JSON(http.StatusOK, models.SuccessResponse{
			Code:    http.StatusOK,
			Message: "Events retrieved successfully",
			Result:  events,
		})
	}
}
