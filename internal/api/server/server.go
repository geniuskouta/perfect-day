package server

import (
	"net/http"
	"perfect-day/internal/api/handlers"
	"perfect-day/internal/api/middleware"
	"perfect-day/internal/api/routes"
	"perfect-day/pkg/auth"
	"perfect-day/pkg/config"
	"perfect-day/pkg/places"
	"perfect-day/pkg/search"
	"perfect-day/pkg/storage"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router        *gin.Engine
	config        *config.Config
	Storage       *storage.Storage
	AuthService   *auth.AuthService
	PlacesService *places.PlacesService
	SearchService *search.SearchService
}

func NewServer(cfg *config.Config) *Server {
	// Initialize storage
	storage := storage.NewStorage(cfg.DataDir)
	if err := storage.Initialize(); err != nil {
		panic("Failed to initialize storage: " + err.Error())
	}

	// Initialize services
	authService := auth.NewAuthService(storage.UserStorage)
	placesService, _ := places.NewPlacesService(cfg.GooglePlacesAPIKey)
	searchService := search.NewSearchService()

	// Create server
	server := &Server{
		config:        cfg,
		Storage:       storage,
		AuthService:   authService,
		PlacesService: placesService,
		SearchService: searchService,
	}

	// Setup router
	server.setupRouter()

	return server
}

func (s *Server) setupRouter() {
	// Set gin mode based on environment
	gin.SetMode(gin.ReleaseMode)

	s.router = gin.New()

	// Add middleware
	s.router.Use(gin.Logger())
	s.router.Use(gin.Recovery())
	s.router.Use(middleware.CORS())

	// Create handlers
	handlers := &handlers.Handlers{
		AuthService:   s.AuthService,
		Storage:       s.Storage,
		PlacesService: s.PlacesService,
		SearchService: s.SearchService,
	}

	// Setup routes
	routes.SetupRoutes(s.router, handlers)
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(w, req)
}