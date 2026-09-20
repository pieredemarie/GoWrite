package main

import (
	"context"
	"gowrite/internal/handler"
	"gowrite/internal/repository"
	"gowrite/internal/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var allowedOrigins = map[string]bool{
	"http://localhost:5173":            true,
	"https://gowrite-frontend.vercel.app": true,
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Edit-Token")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	dsn := os.Getenv("DATABASE_URL")
	port := os.Getenv("PORT")

	repo, err := repository.NewPostgresRepo(dsn)
	if err != nil {
		log.Fatalf("failed to init repo: %v", err)
	}
	defer repo.Close()

	articleService := service.NewArticleService(repo)
	articleHandler := handler.NewArticleHandler(articleService)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/articles", articleHandler.CreateHandler)
	mux.HandleFunc("GET /api/articles/{slug}", articleHandler.GetArticleBySlug)
	mux.HandleFunc("PUT /api/articles/{slug}", articleHandler.UpdateArticle)
	mux.HandleFunc("DELETE /api/articles/{slug}", articleHandler.DeleteArticle)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMiddleware(mux),
	}

	go func() {
		log.Printf("listening on: %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("server error : %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("shutdown error: %v", err)
	}

}
