package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	logrus "github.com/sirupsen/logrus"

	"github.com/ahmadjafari86/go-todo-list/internal/config"
	"github.com/ahmadjafari86/go-todo-list/internal/db"
	"github.com/ahmadjafari86/go-todo-list/internal/handlers"
	"github.com/ahmadjafari86/go-todo-list/internal/middleware"
	"github.com/ahmadjafari86/go-todo-list/internal/models"
	"github.com/ahmadjafari86/go-todo-list/internal/repository"
	"github.com/ahmadjafari86/go-todo-list/internal/service"
)

func main() {
	// load env
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found, relying on environment variables")
	}

	cfg := config.New()

	// configure structured logger
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)

	if cfg.DatabaseURL == "" {
		logrus.Fatal("DATABASE_URL is required")
	}

	// init DB with pool settings and retry
	dbConn, err := db.New(cfg.DatabaseURL, cfg.DBMaxOpenConns, cfg.DBMaxIdleConns, cfg.DBConnMaxLifetime)
	if err != nil {
		logrus.Fatalf("failed to connect to db: %v", err)
	}

	// Auto-migrate models (careful in prod)
	dbConn.AutoMigrate(&models.Todo{}, &models.User{})

	// Wire dependencies
	userRepo := repository.NewGormUserRepository(dbConn)
	todoRepo := repository.NewGormTodoRepository(dbConn)

	authSvc := service.NewAuthService(userRepo)
	todoSvc := service.NewTodoService(todoRepo)

	authH := handlers.NewAuthHandler(authSvc)
	todoH := handlers.NewTodoHandler(todoSvc)

	// gin setup
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestLogger())

	// routes
	r.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.POST("/auth/register", authH.Register)
	r.POST("/auth/login", authH.Login)

	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(authSvc))
	{
		api.GET("/todos", todoH.ListTodos)
		api.GET("/todos/:id", todoH.GetTodo)
		api.POST("/todos", todoH.CreateTodo)
		api.PUT("/todos/:id", todoH.UpdateTodo)
		api.PATCH("/todos/:id/complete", todoH.ToggleComplete)
		api.DELETE("/todos/:id", todoH.DeleteTodo)
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.PortOrDefault()),
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.IdleTimeout) * time.Second,
	}

	// start server
	go func() {
		logrus.Infof("server starting on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("listen: %s\n", err)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logrus.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatalf("server forced to shutdown: %v", err)
	}
	logrus.Info("server exiting")
}
