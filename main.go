package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/kireeti-28/chirpy/internal/database"
)

type apiConfig struct {
	fileServerHits int
	DB             *database.DB
	jwtSecret      string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	cfg := apiConfig{
		fileServerHits: 0,
		DB:             db,
		jwtSecret:      os.Getenv("JWT_SECRET"),
	}

	router := chi.NewRouter()

	router.Use(middlewareCors)

	fsHandler := cfg.middlewareMerticInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))) // fileserver handler

	router.Handle("/app", fsHandler)
	router.Handle("/app/*", fsHandler)

	apiRouter := chi.NewRouter()

	apiRouter.Get("/healthz", handlerReadniess)
	apiRouter.Post("/chirps", cfg.postChirps)
	apiRouter.Get("/chirps", cfg.getChirps)
	apiRouter.Get("/chirps/{chirpId}", cfg.getChripById)
	apiRouter.Delete("/chirps/{chripId}", cfg.handlerDeleteChirp)
	apiRouter.Post("/users", cfg.createUser)
	apiRouter.Put("/users", cfg.userUpdate)
	apiRouter.Post("/login", cfg.loginUser)
	apiRouter.Post("/refresh", cfg.handlerRefresh)
	apiRouter.Post("/revoke", cfg.handlerRevoke)

	apiRouter.Post("/polka/webhooks", cfg.handlerPolkaWebhook)
	router.Mount("/api", apiRouter)

	adminRouter := chi.NewRouter()
	adminRouter.Get("/metrics", cfg.handlerMetrics)
	router.Mount("/admin", adminRouter)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	log.Printf("Starting server at port: %v \n", port)
	err = srv.ListenAndServe()
	log.Fatal(err)
}
