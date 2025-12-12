package main

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"farm-time/internal/auth"
	"farm-time/internal/db"
	"farm-time/internal/handlers"
)

//go:embed web/dist/*
var webFS embed.FS

func main() {
	ctx := context.Background()

	// Connect to database
	database, err := db.New(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Run migrations
	if err := database.Migrate(ctx); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Database migrations completed")

	// Setup handlers
	h := handlers.New(database)
	authHandler := auth.NewHandler(database)

	// Setup router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	// Health check
	r.Get("/health", h.Health)

	// Auth routes (public)
	r.Route("/auth", func(r chi.Router) {
		r.Get("/google/login", authHandler.GoogleLogin)
		r.Get("/google/callback", authHandler.GoogleCallback)
		r.Post("/logout", authHandler.Logout)
		r.Get("/me", authHandler.Me)
	})

	// API routes (protected)
	r.Route("/api", func(r chi.Router) {
		r.Use(authHandler.RequireAuth)

		// Admin routes
		r.Route("/admin", func(r chi.Router) {
			r.Get("/users", h.ListUsers)
			r.Put("/users/{userId}", h.UpdateUserPermissions)
		})

		// Events
		r.Route("/events", func(r chi.Router) {
			r.Get("/", h.ListEvents)
			r.Post("/", h.CreateEvent)
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", h.GetEventWithAll)
				r.Put("/", h.UpdateEvent)
				r.Delete("/", h.DeleteEvent)

				// Attendees for an event
				r.Get("/attendees", h.ListAttendees)
				r.Post("/attendees", h.AddAttendee)
				r.Put("/attendees/{attendeeId}", h.UpdateAttendee)
				r.Delete("/attendees/{attendeeId}", h.RemoveAttendee)

				// Meals for an event
				r.Get("/meals", h.ListMeals)
				r.Post("/meals", h.CreateMeal)
				r.Route("/meals/{mealId}", func(r chi.Router) {
					r.Put("/", h.UpdateMeal)
					r.Delete("/", h.DeleteMeal)

					// Items for a meal
					r.Post("/items", h.AddMealItem)
					r.Put("/items/{itemId}", h.UpdateMealItem)
					r.Delete("/items/{itemId}", h.DeleteMealItem)

					// Signups for an item
					r.Post("/items/{itemId}/signup", h.SignupForItem)
					r.Delete("/items/{itemId}/signup", h.RemoveSignup)
				})

				// Todos for an event
				r.Get("/todos", h.ListTodos)
				r.Post("/todos", h.CreateTodo)
				r.Put("/todos/{todoId}", h.UpdateTodo)
				r.Delete("/todos/{todoId}", h.DeleteTodo)
			})
		})
	})

	// Serve static files from embedded filesystem
	webContent, err := fs.Sub(webFS, "web/dist")
	if err != nil {
		log.Fatalf("Failed to get web content: %v", err)
	}
	fileServer := http.FileServer(http.FS(webContent))

	// Serve static files and SPA fallback
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		// Try to serve the file
		path := r.URL.Path
		if path == "/" {
			path = "/index.html"
		}

		// Check if file exists
		if _, err := fs.Stat(webContent, path[1:]); err == nil {
			fileServer.ServeHTTP(w, r)
			return
		}

		// SPA fallback: serve index.html for client-side routing
		r.URL.Path = "/"
		fileServer.ServeHTTP(w, r)
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}()

	log.Printf("Server starting on port %s", port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}
}
