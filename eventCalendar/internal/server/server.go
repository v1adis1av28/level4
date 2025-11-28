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
	//	POST /create_event — создание нового события;
	//
	// POST /update_event — обновление существующего;
	// POST /delete_event — удаление;
	// GET /events_for_day — получить все события на день;
	// GET /events_for_week — события на неделю;
	// GET /events_for_month — события на месяц.
	s.Router.POST("/create_event", handlers.CreateEventHandler(s.Storage))
	s.Router.POST("/update_event", handlers.UpdateEventHandler(s.Storage))
	// s.Router.POST("/delete_event", s.DeleteEventHandler)
	// s.Router.GET("/events_for_day", s.GetEventsForDayHandler)
	// s.Router.GET("/events_for_week", s.GetEventsForWeekHandler)
	// s.Router.GET("/events_for_month", s.GetEventsForMonthHandler)
}
