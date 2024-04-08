package main

import (
	"banner/internal/cache"
	"banner/internal/handler"
	"banner/internal/middleware"
	"banner/internal/repo"
	"banner/internal/service"
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func newDB(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

func getPostgresDB(ctx context.Context) *pgxpool.Pool {
	psgDsn, ok := os.LookupEnv("POSTGRES_DB_DSN")
	if !ok {
		panic("no POSTGRES_DB_DSN in env vars")
	}
	database, err := newDB(ctx, psgDsn)
	if err != nil {
		log.Panic(err)
	}
	return database
}

func register(router *mux.Router, bannerHandler *handler.BannerHandler) http.Handler {
	router.HandleFunc("/user_banner", bannerHandler.GetUserBanner).Methods(http.MethodGet)

	router.HandleFunc("/banner", bannerHandler.BannerList).Methods(http.MethodGet)
	router.HandleFunc("/banner", bannerHandler.CreateBanner).Methods(http.MethodPost)

	router.HandleFunc("/banner/{id:[0-9]+}", bannerHandler.UpdatePatial).Methods(http.MethodPatch)
	router.HandleFunc("/banner/{id:[0-9]+}", bannerHandler.DeleteBanner).Methods(http.MethodDelete)

	return middleware.AuthMiddleware(router)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database := getPostgresDB(ctx)
	defer database.Close()

	bannerRepo := repo.NewBannerRepo(database)
	bannerCache := cache.NewBannerCache()

	bannerService := service.NewBannerService(bannerRepo, bannerCache)
	bannerHandler := handler.NewBannerHandler(bannerService)

	router := mux.NewRouter()
	register(router, &bannerHandler)

	addr, ok := os.LookupEnv("HOST_PORT")
	if !ok {
		panic("no HOST_PORT in env vars")
	}

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Panic(err)
	}
}
