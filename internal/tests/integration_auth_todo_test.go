package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/ahmadjafari86/go-todo-list/internal/handlers"
	"github.com/ahmadjafari86/go-todo-list/internal/middleware"
	"github.com/ahmadjafari86/go-todo-list/internal/models"
	"github.com/ahmadjafari86/go-todo-list/internal/repository"
	"github.com/ahmadjafari86/go-todo-list/internal/service"

	"github.com/tidwall/gjson"
)

var dbAuth *gorm.DB
var routerAuth *gin.Engine

func setupAuthDB(t *testing.T) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:15",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "testpass",
			"POSTGRES_USER":     "testuser",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").WithStartupTimeout(60 * 1),
	}

	pgC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("failed to start container: %v", err)
	}

	t.Cleanup(func() {
		pgC.Terminate(ctx)
	})

	host, _ := pgC.Host(ctx)
	port, _ := pgC.MappedPort(ctx, "5432")

	dsn := "postgres://testuser:testpass@" + host + ":" + port.Port() + "/testdb?sslmode=disable"
	os.Setenv("DATABASE_URL", dsn)
	os.Setenv("JWT_SECRET", "testsecret")
	os.Setenv("JWT_EXP_MINUTES", "15")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect db: %v", err)
	}

	db.AutoMigrate(&models.User{}, &models.Todo{})

	dbAuth = db
}

func setupAuthRouter() {
	userRepo := repository.NewGormUserRepository(dbAuth)
	todoRepo := repository.NewGormTodoRepository(dbAuth)

	authSvc := service.NewAuthService(userRepo)
	todoSvc := service.NewTodoService(todoRepo)

	authHandler := handlers.NewAuthHandler(authSvc)
	todoHandler := handlers.NewTodoHandler(todoSvc)

	r := gin.Default()
	r.POST("/auth/register", authHandler.Register)
	r.POST("/auth/login", authHandler.Login)

	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(authSvc))
	{
		api.POST("/todos", todoHandler.CreateTodo)
		api.GET("/todos", todoHandler.ListTodos)
	}
	routerAuth = r
}

func registerUser(t *testing.T, username, password string) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/register", strings.NewReader(`{"email":"`+username+`","password":"`+password+`"}`))
	req.Header.Set("Content-Type", "application/json")
	routerAuth.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func loginUserAndGetToken(t *testing.T, username, password string) string {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/login", strings.NewReader(`{"email":"`+username+`","password":"`+password+`"}`))
	req.Header.Set("Content-Type", "application/json")
	routerAuth.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	token := gjson.Get(w.Body.String(), "token").String()
	return token
}

func TestRegisterLoginAndCreateTodoAndIsolation(t *testing.T) {
	setupAuthDB(t)
	setupAuthRouter()

	// Register user1
	registerUser(t, "u1@example.com", "pass1234")
	token1 := loginUserAndGetToken(t, "u1@example.com", "pass1234")
	assert.NotEmpty(t, token1)

	// Create Todo with user1
	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("POST", "/api/todos", strings.NewReader(`{"title":"Owned Todo"}`))
	req3.Header.Set("Content-Type", "application/json")
	req3.Header.Set("Authorization", "Bearer "+token1)
	routerAuth.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusCreated, w3.Code)

	// List Todos as user1
	w4 := httptest.NewRecorder()
	req4, _ := http.NewRequest("GET", "/api/todos", nil)
	req4.Header.Set("Authorization", "Bearer "+token1)
	routerAuth.ServeHTTP(w4, req4)
	assert.Equal(t, http.StatusOK, w4.Code)
	assert.Contains(t, w4.Body.String(), "Owned Todo")

	// Register user2
	registerUser(t, "u2@example.com", "pass1234")
	token2 := loginUserAndGetToken(t, "u2@example.com", "pass1234")
	assert.NotEmpty(t, token2)

	// List Todos as user2 (should NOT see user1's todo)
	w7 := httptest.NewRecorder()
	req7, _ := http.NewRequest("GET", "/api/todos", nil)
	req7.Header.Set("Authorization", "Bearer "+token2)
	routerAuth.ServeHTTP(w7, req7)
	assert.Equal(t, http.StatusOK, w7.Code)
	assert.NotContains(t, w7.Body.String(), "Owned Todo")
}
