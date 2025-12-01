package internal

import (
	"fmt"
	"log"
	"net/http"
	"task-service/internal/config"
	"task-service/internal/entity"
	"task-service/internal/handler"
	"task-service/internal/helper"
	"task-service/internal/infrastructure"
	"task-service/internal/middleware"
	"task-service/internal/repository"
	"task-service/internal/usecase"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Server struct {
	taskUc usecase.TaskUc
	host   string
	router *http.ServeMux
}

// Middleware untuk CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Izinkan domain frontend kamu di sini (misal: http://localhost:5173 atau domain production)
		w.Header().Set("Access-Control-Allow-Origin", "*") // gunakan "*" hanya untuk dev
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight request (OPTIONS)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) initRoute() {
	taskMux := http.NewServeMux()

	cfg := config.NewConfig()
	jwtService := infrastructure.NewJWTService(*cfg)

	jwtMiddleware := middleware.NewMiddleware(jwtService)

	taskHandler := handler.NewTaskHandler(s.taskUc, taskMux, jwtMiddleware)

	taskHandler.SetupRoutes()

	// Bungkus authMux dengan middleware CORS
	s.router.Handle(helper.ApiGrup+"/", corsMiddleware(http.StripPrefix(helper.ApiGrup, taskMux)))

	log.Println("✅ ROUTES SETUP COMPLETE ON PREFIX :", helper.ApiGrup)
}

func (s *Server) Run() {
	s.initRoute()
	log.Printf("🚀 TASK SERVICE STARTING ON HOST %s", s.host)

	if err := http.ListenAndServe(s.host, s.router); err != nil {
		panic(fmt.Errorf("server not running on host %s, because of error %v", s.host, err))
	}
}

func NewServer() *Server {
	// Load Konfigurasi
	cfg := config.NewConfig()

	// Koneksi Database (GORM + Postgres)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.DBConfig.Host, cfg.DBConfig.User, cfg.DBConfig.Password, cfg.DBConfig.Name, cfg.DBConfig.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("failed to connect database: %v", err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Errorf("failed to get underlying *sql.DB: %v", err))
	}

	// Konfigurasi connection pool
	sqlDB.SetMaxIdleConns(cfg.DBConfig.MaxIdle)
	sqlDB.SetMaxOpenConns(cfg.DBConfig.MaxOpen)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.DBConfig.MaxLife) * time.Minute)
	sqlDB.SetConnMaxIdleTime(time.Duration(cfg.DBConfig.MaxIdleTime) * time.Minute)

	// Cek koneksi
	if err = sqlDB.Ping(); err != nil {
		panic(fmt.Errorf("failed to ping database: %v", err))
	}

	// Migrasi Database
	if err := db.AutoMigrate(&entity.Task{}); err != nil {
		panic(fmt.Errorf("GORM AutoMigrate failed: %v", err))
	}
	log.Println("✅ CONNECTING DATABASE SUCCESS")

	// Inisialisasi Server
	host := fmt.Sprintf(":%s", cfg.ServerPort)
	router := http.NewServeMux()

	// Dependency Injection
	taskRepo := repository.NewTaskRepo(db)
	taskUc := usecase.NewTaskUc(taskRepo)

	return &Server{
		taskUc: taskUc,
		host:   host,
		router: router,
	}
}
