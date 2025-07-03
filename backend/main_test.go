package main

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "os"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func setupTestDB() {
    var err error
    db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        panic("Failed to connect to test database")
    }
    
    db.AutoMigrate(&User{}, &TravelRequest{})
}

func setupTestRouter() *gin.Engine {
    gin.SetMode(gin.TestMode)
    
    r := gin.New()
    
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })
    
    auth := r.Group("/api/auth")
    {
        auth.POST("/register", registerHandler)
        auth.POST("/login", loginHandler)
    }
    
    api := r.Group("/api/travel-requests")
    api.Use(authMiddleware())
    {
        api.POST("", createTravelRequestHandler)
        api.GET("", listTravelRequestsHandler)
        api.GET("/:id", getTravelRequestHandler)
        api.PUT("/:id/status", updateStatusHandler)
        api.DELETE("/:id", cancelTravelRequestHandler)
    }
    
    return r
}

func TestMain(m *testing.M) {
    setupTestDB()
    code := m.Run()
    os.Exit(code)
}

func TestHealthCheck(t *testing.T) {
    router := setupTestRouter()
    
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/health", nil)
    router.ServeHTTP(w, req)
    
    assert.Equal(t, 200, w.Code)
    assert.Contains(t, w.Body.String(), "ok")
}

func TestUserRegistration(t *testing.T) {
    setupTestDB()
    router := setupTestRouter()
    
    user := RegisterRequest{
        Name:     "Test User",
        Email:    "test@example.com",
        Password: "password123",
    }
    
    jsonValue, _ := json.Marshal(user)
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonValue))
    req.Header.Set("Content-Type", "application/json")
    
    router.ServeHTTP(w, req)
    
    assert.Equal(t, 201, w.Code)
    assert.Contains(t, w.Body.String(), "Usuário criado com sucesso")
}

func TestUserLogin(t *testing.T) {
    setupTestDB()
    router := setupTestRouter()
    
    // First register a user
    user := RegisterRequest{
        Name:     "Test User",
        Email:    "test@example.com",
        Password: "password123",
    }
    
    jsonValue, _ := json.Marshal(user)
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonValue))
    req.Header.Set("Content-Type", "application/json")
    router.ServeHTTP(w, req)
    
    // Then login
    loginReq := LoginRequest{
        Email:    "test@example.com",
        Password: "password123",
    }
    
    jsonValue, _ = json.Marshal(loginReq)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonValue))
    req.Header.Set("Content-Type", "application/json")
    
    router.ServeHTTP(w, req)
    
    assert.Equal(t, 200, w.Code)
    assert.Contains(t, w.Body.String(), "token")
}

func TestCreateTravelRequest(t *testing.T) {
    setupTestDB()
    router := setupTestRouter()
    
    // Register and login to get token
    user := RegisterRequest{
        Name:     "Test User",
        Email:    "test@example.com",
        Password: "password123",
    }
    
    jsonValue, _ := json.Marshal(user)
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonValue))
    req.Header.Set("Content-Type", "application/json")
    router.ServeHTTP(w, req)
    
    // Login
    loginReq := LoginRequest{
        Email:    "test@example.com",
        Password: "password123",
    }
    
    jsonValue, _ = json.Marshal(loginReq)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonValue))
    req.Header.Set("Content-Type", "application/json")
    router.ServeHTTP(w, req)
    
    var loginResponse map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &loginResponse)
    token := loginResponse["token"].(string)
    
    // Create travel request
    travelReq := CreateTravelRequest{
        RequesterName: "Test User",
        Destination:   "São Paulo",
        DepartureDate: "2025-08-15",
        ReturnDate:    "2025-08-20",
    }
    
    jsonValue, _ = json.Marshal(travelReq)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("POST", "/api/travel-requests", bytes.NewBuffer(jsonValue))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+token)
    
    router.ServeHTTP(w, req)
    
    assert.Equal(t, 201, w.Code)
    assert.Contains(t, w.Body.String(), "São Paulo")
}

func TestListTravelRequests(t *testing.T) {
    setupTestDB()
    router := setupTestRouter()
    
    // Register and login
    user := RegisterRequest{
        Name:     "Test User",
        Email:    "test@example.com",
        Password: "password123",
    }
    
    jsonValue, _ := json.Marshal(user)
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonValue))
    req.Header.Set("Content-Type", "application/json")
    router.ServeHTTP(w, req)
    
    loginReq := LoginRequest{
        Email:    "test@example.com",
        Password: "password123",
    }
    
    jsonValue, _ = json.Marshal(loginReq)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonValue))
    req.Header.Set("Content-Type", "application/json")
    router.ServeHTTP(w, req)
    
    var loginResponse map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &loginResponse)
    token := loginResponse["token"].(string)
    
    // List travel requests
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("GET", "/api/travel-requests", nil)
    req.Header.Set("Authorization", "Bearer "+token)
    
    router.ServeHTTP(w, req)
    
    assert.Equal(t, 200, w.Code)
    assert.Contains(t, w.Body.String(), "[]") // Empty array initially
}

func TestUnauthorizedAccess(t *testing.T) {
    setupTestDB()
    router := setupTestRouter()
    
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/api/travel-requests", nil)
    
    router.ServeHTTP(w, req)
    
    assert.Equal(t, 401, w.Code)
    assert.Contains(t, w.Body.String(), "Token de autorização necessário")
}

func TestInvalidLogin(t *testing.T) {
    setupTestDB()
    router := setupTestRouter()
    
    loginReq := LoginRequest{
        Email:    "nonexistent@example.com",
        Password: "wrongpassword",
    }
    
    jsonValue, _ := json.Marshal(loginReq)
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonValue))
    req.Header.Set("Content-Type", "application/json")
    
    router.ServeHTTP(w, req)
    
    assert.Equal(t, 401, w.Code)
    assert.Contains(t, w.Body.String(), "Credenciais inválidas")
}

func TestCreatorCannotUpdateStatus(t *testing.T) {
    setupTestDB()
    router := setupTestRouter()
    
    // Register user and create request
    user := RegisterRequest{
        Name:     "Test User",
        Email:    "test@example.com",
        Password: "password123",
    }
    
    jsonValue, _ := json.Marshal(user)
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonValue))
    req.Header.Set("Content-Type", "application/json")
    router.ServeHTTP(w, req)
    
    // Login
    loginReq := LoginRequest{
        Email:    "test@example.com",
        Password: "password123",
    }
    
    jsonValue, _ = json.Marshal(loginReq)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonValue))
    req.Header.Set("Content-Type", "application/json")
    router.ServeHTTP(w, req)
    
    var loginResponse map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &loginResponse)
    token := loginResponse["token"].(string)
    
    // Create travel request
    travelReq := CreateTravelRequest{
        RequesterName: "Test User",
        Destination:   "Test City",
        DepartureDate: "2025-08-15",
        ReturnDate:    "2025-08-20",
    }
    
    jsonValue, _ = json.Marshal(travelReq)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("POST", "/api/travel-requests", bytes.NewBuffer(jsonValue))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+token)
    router.ServeHTTP(w, req)
    
    // Try to update status (should fail - creator cannot update)
    updateReq := UpdateStatusRequest{Status: "aprovado"}
    jsonValue, _ = json.Marshal(updateReq)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("PUT", "/api/travel-requests/1/status", bytes.NewBuffer(jsonValue))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+token)
    router.ServeHTTP(w, req)
    
    assert.Equal(t, 403, w.Code)
    assert.Contains(t, w.Body.String(), "Você não pode alterar o status")
}
