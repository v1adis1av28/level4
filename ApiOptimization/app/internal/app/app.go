package app

import (
	"context"
	"demo/internal/config"
	"demo/internal/database"
	"demo/internal/handlers"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type App struct {
	DB           *database.DB
	Router       *gin.Engine
	Server       *http.Server
	OrderHandler *handlers.OrderHandler
	Config       *config.Config
}

func NewApp(db *database.DB, handler *handlers.OrderHandler, cfg *config.Config) *App {
	router := gin.Default()

	router.Use(func(c *gin.Context) {
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

	router.MaxMultipartMemory = 1 << 20

	server := &http.Server{
		Addr:         cfg.App.Port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	app := &App{
		DB:           db,
		Router:       router,
		Server:       server,
		OrderHandler: handler,
		Config:       cfg,
	}

	app.SetupRoutes()

	return app
}

func (a *App) MustStart() {
	log.Printf("Starting server on %s", a.Config.App.Port)
	if err := a.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Server failed to start: ", err)
	}
}

func (app *App) SetupRoutes() error {
	app.Router.POST("/order", app.OrderHandler.HandleIncomingOrder)
	app.Router.GET("/order/:id", app.OrderHandler.GetOrderById)
	return nil
}

func (app *App) Stop() {
	log.Println("Shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.Server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	} else {
		log.Println("Server stopped gracefully")
	}
}
