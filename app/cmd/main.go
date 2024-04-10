package main

import (
	cache "banner/internal/banner_cache"
	"banner/internal/db"
	"banner/internal/handler"
	"banner/internal/middleware"
	"banner/internal/repo"
	"banner/internal/service"
	"time"

	"context"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
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

func getRedisClient(ctx context.Context) *redis.Client {
	redisAddr, ok := os.LookupEnv("REDIS_ADDR")
	if !ok {
		panic("no REDIS_ADDR in env vars")
	}

	redisPassword, ok := os.LookupEnv("REDIS_PASSWORD")
	if !ok {
		panic("no REDIS_PASSWORD in env vars")
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       0,
	})

	_, err := redisClient.Do(ctx, "PING").Result()
	if err != nil {
		panic(err)
	}

	return redisClient
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

	redisClient := getRedisClient(ctx)
	defer redisClient.Close()

	bannerRepo := repo.NewBannerRepo(database)
	bannerCache := cache.NewBannerRedisCahe(redisClient, 5*time.Minute)

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
