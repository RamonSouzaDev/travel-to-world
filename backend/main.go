package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "time"
    "strings"

    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var db *gorm.DB
var jwtSecret = []byte("travel-requests-jwt-secret-2024-secure-key")

// Models
type User struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Name      string    `json:"name"`
    Email     string    `json:"email" gorm:"uniqueIndex"`
    Password  string    `json:"-"`
    CreatedAt time.Time `json:"created_at"`
}

type TravelRequest struct {
    ID            uint      `json:"id" gorm:"primaryKey"`
    RequesterName string    `json:"requester_name"`
    Destination   string    `json:"destination"`
    DepartureDate time.Time `json:"departure_date"`
    ReturnDate    time.Time `json:"return_date"`
    Status        string    `json:"status" gorm:"default:'solicitado'"`
    UserID        uint      `json:"user_id"`        // Usuário que pode ver o pedido
    CreatedByID   uint      `json:"created_by_id"`  // Usuário que criou (NÃO pode alterar)
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
}

// Request DTOs
type RegisterRequest struct {
    Name     string `json:"name" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

type CreateTravelRequest struct {
    RequesterName string `json:"requester_name" binding:"required"`
    Destination   string `json:"destination" binding:"required"`
    DepartureDate string `json:"departure_date" binding:"required"`
    ReturnDate    string `json:"return_date" binding:"required"`
}

type UpdateStatusRequest struct {
    Status string `json:"status" binding:"required"`
}

func main() {
    print_status("Iniciando Travel Requests Backend...")
    
    setupDatabase()
    setupRoutes()
}

func print_status(msg string) {
    log.Printf("[INFO] %s", msg)
}

func setupDatabase() {
    print_status("Conectando ao banco de dados...")
    
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
        getEnv("DB_HOST", "postgres"),
        getEnv("DB_USER", "postgres"),
        getEnv("DB_PASSWORD", "postgres"),
        getEnv("DB_NAME", "travel_requests"),
        getEnv("DB_PORT", "5432"))

    var err error
    for i := 0; i < 10; i++ {
        db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
        if err == nil {
            break
        }
        print_status(fmt.Sprintf("Tentativa %d de conectar ao banco... erro: %v", i+1, err))
        time.Sleep(time.Second * 3)
    }
    
    if err != nil {
        log.Fatal("Falha ao conectar com o banco após 10 tentativas:", err)
    }

    // Auto migrate
    if err := db.AutoMigrate(&User{}, &TravelRequest{}); err != nil {
        log.Fatal("Falha nas migrations:", err)
    }
    
    print_status("Banco de dados conectado e migrado com sucesso!")
}

func setupRoutes() {
    if getEnv("ENV", "development") == "production" {
        gin.SetMode(gin.ReleaseMode)
    }

    r := gin.Default()

    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000", "http://frontend", "http://frontend:80", "*"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        AllowCredentials: true,
    }))

    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "status": "ok",
            "time":   time.Now(),
            "service": "travel-requests-backend",
        })
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

    port := getEnv("PORT", "8080")
    print_status(fmt.Sprintf("Servidor iniciando na porta %s", port))
    log.Fatal(r.Run(":" + port))
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func authMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Token de autorização necessário"})
            c.Abort()
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
            }
            return jwtSecret, nil
        })

        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
            c.Abort()
            return
        }

        if claims, ok := token.Claims.(jwt.MapClaims); ok {
            userID := uint(claims["user_id"].(float64))
            c.Set("user_id", userID)
        } else {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Claims inválidas"})
            c.Abort()
            return
        }

        c.Next()
    }
}

func registerHandler(c *gin.Context) {
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var existingUser User
    if err := db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
        c.JSON(http.StatusConflict, gin.H{"error": "Email já está em uso"})
        return
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar senha"})
        return
    }

    user := User{
        Name:     req.Name,
        Email:    req.Email,
        Password: string(hashedPassword),
    }

    if err := db.Create(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar usuário"})
        return
    }

    print_status(fmt.Sprintf("Novo usuário registrado: %s (%s)", user.Name, user.Email))
    c.JSON(http.StatusCreated, gin.H{
        "message": "Usuário criado com sucesso",
        "user": gin.H{
            "id":    user.ID,
            "name":  user.Name,
            "email": user.Email,
        },
    })
}

func loginHandler(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var user User
    if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciais inválidas"})
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciais inválidas"})
        return
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID,
        "email":   user.Email,
        "exp":     time.Now().Add(time.Hour * 24).Unix(),
        "iat":     time.Now().Unix(),
    })

    tokenString, err := token.SignedString(jwtSecret)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar token"})
        return
    }

    print_status(fmt.Sprintf("Login realizado: %s", user.Email))
    c.JSON(http.StatusOK, gin.H{
        "token": tokenString,
        "user": gin.H{
            "id":    user.ID,
            "name":  user.Name,
            "email": user.Email,
        },
    })
}

func createTravelRequestHandler(c *gin.Context) {
    userID, _ := c.Get("user_id")
    
    var req CreateTravelRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    departureDate, err := time.Parse("2006-01-02", req.DepartureDate)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de data de ida inválido (use YYYY-MM-DD)"})
        return
    }

    returnDate, err := time.Parse("2006-01-02", req.ReturnDate)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de data de volta inválido (use YYYY-MM-DD)"})
        return
    }

    if returnDate.Before(departureDate) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Data de volta deve ser posterior à data de ida"})
        return
    }

    userIDValue := userID.(uint)
    travelRequest := TravelRequest{
        RequesterName: req.RequesterName,
        Destination:   req.Destination,
        DepartureDate: departureDate,
        ReturnDate:    returnDate,
        Status:        "solicitado",
        UserID:        userIDValue,     // Usuário que pode ver
        CreatedByID:   userIDValue,     // Usuário que criou (não pode alterar status)
    }

    if err := db.Create(&travelRequest).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar pedido de viagem"})
        return
    }

    print_status(fmt.Sprintf("Novo pedido criado: %s para %s", req.RequesterName, req.Destination))
    log.Printf("📧 [NOTIFICATION] Pedido de viagem criado para %s - Destino: %s", 
        travelRequest.RequesterName, travelRequest.Destination)

    c.JSON(http.StatusCreated, travelRequest)
}

func listTravelRequestsHandler(c *gin.Context) {
    var requests []TravelRequest
    query := db.Model(&TravelRequest{})

    // Filtros implementados
    if status := c.Query("status"); status != "" {
        query = query.Where("status = ?", status)
    }
    if destination := c.Query("destination"); destination != "" {
        query = query.Where("destination ILIKE ?", "%"+destination+"%")
    }
    
    // NOVO: Filtros por período
    if startDate := c.Query("start_date"); startDate != "" {
        if parsed, err := time.Parse("2006-01-02", startDate); err == nil {
            query = query.Where("departure_date >= ?", parsed)
        }
    }
    if endDate := c.Query("end_date"); endDate != "" {
        if parsed, err := time.Parse("2006-01-02", endDate); err == nil {
            query = query.Where("return_date <= ?", parsed)
        }
    }
    
    // Filtro por período de criação
    if createdAfter := c.Query("created_after"); createdAfter != "" {
        if parsed, err := time.Parse("2006-01-02", createdAfter); err == nil {
            query = query.Where("created_at >= ?", parsed)
        }
    }
    if createdBefore := c.Query("created_before"); createdBefore != "" {
        if parsed, err := time.Parse("2006-01-02", createdBefore); err == nil {
            query = query.Where("created_at <= ?", parsed.Add(24*time.Hour))
        }
    }

    query.Order("created_at DESC").Find(&requests)
    
    if requests == nil {
        requests = []TravelRequest{}
    }
    
    c.JSON(http.StatusOK, requests)
}

func getTravelRequestHandler(c *gin.Context) {
    id := c.Param("id")

    var request TravelRequest
    if err := db.Where("id = ?", id).First(&request).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Pedido de viagem não encontrado"})
        return
    }

    c.JSON(http.StatusOK, request)
}

func updateStatusHandler(c *gin.Context) {
    userID, _ := c.Get("user_id")
    id := c.Param("id")

    var request TravelRequest
    if err := db.Where("id = ?", id).First(&request).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Pedido de viagem não encontrado"})
        return
    }

    // REGRA DE NEGÓCIO: Usuário que criou o pedido NÃO pode alterar o status
    if request.CreatedByID == userID.(uint) {
        c.JSON(http.StatusForbidden, gin.H{
            "error": "Você não pode alterar o status de um pedido que você mesmo criou. Outro usuário deve fazer essa alteração.",
        })
        return
    }

    var req UpdateStatusRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    validStatuses := map[string]bool{
        "solicitado": true,
        "aprovado":   true,
        "cancelado":  true,
    }
    
    if !validStatuses[req.Status] {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Status inválido. Use: solicitado, aprovado ou cancelado"})
        return
    }

    oldStatus := request.Status
    request.Status = req.Status
    
    if err := db.Save(&request).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar status"})
        return
    }

    print_status(fmt.Sprintf("Status atualizado: %s -> %s para %s", oldStatus, req.Status, request.RequesterName))
    log.Printf("📧 [NOTIFICATION] Status do pedido de %s alterado de '%s' para '%s'", 
        request.RequesterName, oldStatus, request.Status)

    c.JSON(http.StatusOK, request)
}

func cancelTravelRequestHandler(c *gin.Context) {
    id := c.Param("id")

    var request TravelRequest
    if err := db.Where("id = ?", id).First(&request).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Pedido de viagem não encontrado"})
        return
    }

    if request.Status == "cancelado" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Pedido já está cancelado"})
        return
    }

    // Regra de negócio: Permite cancelar pedidos aprovados
    oldStatus := request.Status
    request.Status = "cancelado"
    
    if err := db.Save(&request).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao cancelar pedido"})
        return
    }

    print_status(fmt.Sprintf("Pedido cancelado: %s (era %s)", request.RequesterName, oldStatus))
    log.Printf("📧 [NOTIFICATION] Pedido de viagem de %s foi cancelado (anterior: %s)", 
        request.RequesterName, oldStatus)

    c.JSON(http.StatusOK, gin.H{
        "message": "Pedido cancelado com sucesso",
        "request": request,
    })
}
