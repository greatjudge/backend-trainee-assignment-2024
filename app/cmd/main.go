package main

import (
	"banner/internal/cache"
	"banner/internal/db"
	"banner/internal/handler"
	"banner/internal/middleware"
	"banner/internal/repo"
	"banner/internal/service"

	"context"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func getPostgresDB(ctx context.Context) *db.Database {
	psgDsn, ok := os.LookupEnv("POSTGRES_DB_DSN")
	if !ok {
		panic("no POSTGRES_DB_DSN in env vars")
	}
	database, err := db.NewDB(ctx, psgDsn)
	if err != nil {
		log.Panic(err)
	}
	return database
}
func register(router *mux.Router, bannerHandler *handler.BannerHandler) {
	router.HandleFunc("/user_banner", bannerHandler.GetUserBanner).Methods(http.MethodGet)

	router.Handle(
		"/banner",
		middleware.OnlyAdmin((http.HandlerFunc(bannerHandler.BannerList))),
	).Methods(http.MethodGet)

	router.Handle(
		"/banner",
		middleware.OnlyAdmin((http.HandlerFunc(bannerHandler.CreateBanner))),
	).Methods(http.MethodPost)

	router.Handle(
		"/banner/{id:[0-9]+}",
		middleware.OnlyAdmin((http.HandlerFunc(bannerHandler.UpdatePatial))),
	).Methods(http.MethodPatch)

	router.Handle(
		"/banner/{id:[0-9]+}",
		middleware.OnlyAdmin((http.HandlerFunc(bannerHandler.DeleteBanner))),
	).Methods(http.MethodDelete)
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

	userToken, ok := os.LookupEnv("USER_TOKEN")
	if !ok {
		panic("no USER_TOKEN in env vars")
	}

	adminToken, ok := os.LookupEnv("ADMIN_TOKEN")
	if !ok {
		panic("no ADMIN_TOKEN in env vars")
	}

	appHandler := middleware.AuthMiddleware(userToken, adminToken, router)

	addr, ok := os.LookupEnv("HOST_PORT")
	if !ok {
		panic("no HOST_PORT in env vars")
	}

	if err := http.ListenAndServe(addr, appHandler); err != nil {
		log.Panic(err)
	}
}
