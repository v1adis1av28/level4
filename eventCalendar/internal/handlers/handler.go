package handlers

import (
	"eventCalendar/internal/models"
	"eventCalendar/internal/storage"
	"eventCalendar/internal/utils"
	"net/http"

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
