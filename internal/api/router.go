package api

import (
	"io/fs"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/aevynt/goPanel/internal/apps"
	"github.com/aevynt/goPanel/internal/auth"
	"github.com/aevynt/goPanel/internal/caddy"
	"github.com/aevynt/goPanel/internal/config"
	"github.com/aevynt/goPanel/internal/database"
	"github.com/aevynt/goPanel/internal/docker"
	"github.com/aevynt/goPanel/internal/filemanager"
	"github.com/aevynt/goPanel/internal/middleware"
	"github.com/aevynt/goPanel/internal/ports"
	"github.com/aevynt/goPanel/internal/servicemanager"
	"github.com/aevynt/goPanel/internal/update"
	"github.com/rs/cors"
)

type Server struct {
	cfg           *config.Config
	db            *database.DB
	auth          *auth.JWTManager
	sm            servicemanager.Manager
	fm            *filemanager.Manager
	caddy         *caddy.Client
	pm            *ports.Manager
	docker        *docker.Service
	apps          *apps.Service
	wsHub         *WebSocketHub
	webFS         fs.FS
	updateChecker *update.Checker
}

func NewServer(
	cfg *config.Config,
	db *database.DB,
	sm servicemanager.Manager,
	fm *filemanager.Manager,
	cc *caddy.Client,
	pm *ports.Manager,
	dockerSvc *docker.Service,
	appsSvc *apps.Service,
	webFS fs.FS,
) *Server {
	jm := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiry)
	return &Server{
		cfg:           cfg,
		db:            db,
		auth:          jm,
		sm:            sm,
		fm:            fm,
		caddy:         cc,
		pm:            pm,
		docker:        dockerSvc,
		apps:          appsSvc,
		wsHub:         NewWebSocketHub(),
		webFS:         webFS,
		updateChecker: update.NewChecker(),
	}
}

func (s *Server) Handler() http.Handler {
	r := chi.NewRouter()

	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.RealIP)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})
	r.Use(corsHandler.Handler)

	// Root-level WebSocket endpoint matching the frontend Dashboard.vue connection URL
	r.Get("/ws", s.DashboardWS)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/login", s.Login)
		r.Post("/auth/refresh", s.RefreshToken)

		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthRequired(s.auth))

			r.Get("/auth/me", s.Me)
			r.Post("/auth/change-password", s.ChangePassword)
			r.Post("/auth/2fa/setup", s.Setup2FA)
			r.Post("/auth/2fa/enable", s.Enable2FA)
			r.Post("/auth/2fa/disable", s.Disable2FA)

			r.Route("/users", func(r chi.Router) {
				r.Use(middleware.AdminOnly)
				r.Get("/", s.ListUsers)
				r.Post("/", s.CreateUser)
				r.Put("/{id}", s.UpdateUser)
				r.Delete("/{id}", s.DeleteUser)
			})

			r.Route("/services", func(r chi.Router) {
				r.Get("/", s.ListServices)
				r.Get("/{name}", s.GetService)
				r.Get("/{name}/logs", s.GetServiceLogs)

				r.Group(func(r chi.Router) {
					r.Use(middleware.AdminOnly)
					r.Post("/", s.CreateService)
					r.Post("/{name}/start", s.StartService)
					r.Post("/{name}/stop", s.StopService)
					r.Post("/{name}/restart", s.RestartService)
					r.Post("/{name}/enable", s.EnableService)
					r.Post("/{name}/disable", s.DisableService)
					r.Delete("/{name}", s.RemoveService)
				})
			})

			r.Route("/binaries", func(r chi.Router) {
				r.Use(middleware.AdminOnly)
				r.Get("/", s.ListBinaries)
				r.Post("/", s.UploadBinary)
				r.Delete("/{id}", s.DeleteBinary)
			})

			r.Route("/files", func(r chi.Router) {
				r.Get("/", s.ListFiles)
				r.Get("/read", s.ReadFile)

				r.Group(func(r chi.Router) {
					r.Use(middleware.AdminOnly)
					r.Post("/write", s.WriteFile)
					r.Post("/mkdir", s.MkdirDir)
					r.Post("/upload", s.UploadFile)
					r.Delete("/", s.RemoveFile)
					r.Post("/rename", s.RenameFile)
					r.Post("/zip", s.ZipFile)
					r.Post("/unzip", s.UnzipFile)
				})
			})

			r.Route("/docker", func(r chi.Router) {
				r.Use(middleware.AdminOnly)
				r.Get("/containers", s.ListContainers)
				r.Post("/containers/{id}/start", s.StartContainer)
				r.Post("/containers/{id}/stop", s.StopContainer)
				r.Post("/containers/{id}/restart", s.RestartContainer)
				r.Get("/containers/{id}/logs", s.GetContainerLogs)
				r.Post("/compose", s.DeployCompose)
			})

			r.Route("/apps", func(r chi.Router) {
				r.Use(middleware.AdminOnly)
				r.Get("/", s.ListApps)
				r.Post("/deploy", s.DeployApp)
			})

			r.Route("/ports", func(r chi.Router) {
				r.Get("/", s.ListPorts)
				r.Get("/check/{port}", s.CheckPort)
				r.Post("/find", s.FindPort)
			})

			r.Route("/sites", func(r chi.Router) {
				r.Use(middleware.AdminOnly)
				r.Get("/", s.ListSites)
				r.Post("/", s.AddSite)
				r.Delete("/{domain}", s.RemoveSite)
				r.Get("/health", s.CaddyHealth)
			})

			r.Route("/dashboard", func(r chi.Router) {
				r.Get("/stats", s.DashboardStats)
				r.Get("/ws", s.DashboardWS)
			})

			r.Route("/updates", func(r chi.Router) {
				r.Use(middleware.AdminOnly)
				r.Get("/check", s.CheckUpdate)
			})

			r.Route("/settings", func(r chi.Router) {
				r.Use(middleware.AdminOnly)
				r.Get("/", s.GetSettings)
				r.Put("/", s.UpdateSettings)
			})

			r.Route("/public", func(r chi.Router) {
				r.Get("/shares", s.ListShares)
				r.Post("/shares", s.CreateShare)
				r.Delete("/shares/{id}", s.DeleteShare)
				r.Get("/shares/{id}/files", s.ListShareFiles)
				r.Post("/shares/{id}/upload", s.UploadShareFile)
				r.Delete("/shares/{id}/file", s.DeleteShareFile)
				r.Get("/domain", s.PublicDomainHandler)
				r.Put("/domain", s.PublicDomainHandler)
			})
		})
	})

	r.Handle("/public/*", http.StripPrefix("/public/", http.FileServer(http.Dir(s.cfg.PublicDir))))

	var mimeTypes = map[string]string{
		".js":   "application/javascript",
		".css":  "text/css",
		".svg":  "image/svg+xml",
		".png":  "image/png",
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".gif":  "image/gif",
		".ico":  "image/x-icon",
		".woff": "font/woff",
		".woff2": "font/woff2",
		".ttf":  "font/ttf",
		".json": "application/json",
		".html": "text/html",
		".txt":  "text/plain",
		".wasm": "application/wasm",
		".map":  "application/json",
	}

	if s.webFS != nil {
		r.Handle("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := strings.TrimPrefix(r.URL.Path, "/")
			if path == "" {
				path = "index.html"
			}

			data, err := fs.ReadFile(s.webFS, path)
			if err != nil {
				data, err = fs.ReadFile(s.webFS, "index.html")
				if err != nil {
					http.NotFound(w, r)
					return
				}
			}

			ext := strings.ToLower(filepath.Ext(path))
			if mimeType, ok := mimeTypes[ext]; ok {
				w.Header().Set("Content-Type", mimeType)
			} else if mimeType := mime.TypeByExtension(ext); mimeType != "" {
				w.Header().Set("Content-Type", mimeType)
			}

			w.Write(data)
		}))
	}

	return r
}

func (s *Server) StartWShub() {
	go s.wsHub.Run()
}
