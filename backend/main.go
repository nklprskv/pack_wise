package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"pack_wise/internal/app"
	"pack_wise/internal/store"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	ctx := context.Background()
	if err := runMigrations(databaseURL); err != nil {
		log.Fatal(err)
	}

	packSizesStore, err := store.NewPostgresPackSizesStore(ctx, databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer packSizesStore.Close()

	addr := ":" + port
	server := app.NewServer(addr, packSizesStore)
	server.Handler = withCORS(server.Handler, loadAllowedOrigins())

	log.Printf("listening on %s", addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func runMigrations(databaseURL string) error {
	migrationRunner, err := migrate.New("file://internal/db/migrations", databaseURL)
	if err != nil {
		return fmt.Errorf("create migration runner: %w", err)
	}

	defer func() {
		sourceErr, databaseErr := migrationRunner.Close()
		if sourceErr != nil {
			log.Printf("migration source close error: %v", sourceErr)
		}
		if databaseErr != nil {
			log.Printf("migration database close error: %v", databaseErr)
		}
	}()

	if err := migrationRunner.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("run migrations: %w", err)
	}

	return nil
}

func withCORS(next http.Handler, allowOrigins []string) http.Handler {
	allowedOrigins := make(map[string]struct{}, len(allowOrigins))
	for _, origin := range allowOrigins {
		allowedOrigins[origin] = struct{}{}
	}

	allowAllOrigins := len(allowOrigins) == 1 && allowOrigins[0] == "*"

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Origin")
		w.Header().Add("Vary", "Access-Control-Request-Method")
		w.Header().Add("Vary", "Access-Control-Request-Headers")

		origin := r.Header.Get("Origin")
		switch {
		case allowAllOrigins:
			w.Header().Set("Access-Control-Allow-Origin", "*")
		case origin != "":
			if _, ok := allowedOrigins[origin]; ok {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func loadAllowedOrigins() []string {
	rawOrigins := os.Getenv("CORS_ALLOW_ORIGINS")
	if rawOrigins == "" {
		rawOrigins = os.Getenv("CORS_ALLOW_ORIGIN")
	}
	if rawOrigins == "" {
		rawOrigins = "http://localhost:5173,http://localhost:80,http://localhost"
	}

	parts := strings.Split(rawOrigins, ",")
	origins := make([]string, 0, len(parts))
	seen := make(map[string]struct{}, len(parts))

	for _, part := range parts {
		origin := strings.TrimSpace(part)
		if origin == "" {
			continue
		}
		if _, ok := seen[origin]; ok {
			continue
		}

		seen[origin] = struct{}{}
		origins = append(origins, origin)
	}

	if len(origins) == 0 {
		return []string{"http://localhost:5173"}
	}

	return origins
}
