package server

import (
	"eventCalendar/internal/config"
	"eventCalendar/internal/handlers"
	"eventCalendar/internal/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Router     *gin.Engine
	HttpServer *http.Server
	Storage    *storage.Storage
}

func New(Config *config.Config, storage *storage.Storage) *Server {
	server := &Server{Router: gin.New(), Storage: storage}

	server.Router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	server.HttpServer = &http.Server{
		Addr:    Config.App.Port,
		Handler: server.Router,
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	s.Router.POST("/create_event", handlers.CreateEventHandler(s.Storage))
	s.Router.POST("/update_event", handlers.UpdateEventHandler(s.Storage))
	s.Router.POST("/delete_event", handlers.DeleteEventHandler(s.Storage))
	s.Router.GET("/events_for_day", handlers.GetEventsForDayHandler(s.Storage))
	s.Router.GET("/events_for_week", handlers.GetEventsForWeekHandler(s.Storage))
	s.Router.GET("/events_for_month", handlers.GetEventsForMonthHandler(s.Storage))
}
