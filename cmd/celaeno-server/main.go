package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/seandisero/celaeno/internal/server/database"
	"github.com/seandisero/celaeno/internal/server/srvapi"

	// "github.com/go-chi/chi/v5"
	// "github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("warning: assuming default configuration. .env unreadable: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("could not get jwt secret")
		os.Exit(1)
	}

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("libsql", dbURL)
	if err != nil {
		log.Fatal("could not open db")
		os.Exit(1)
	}
	defer db.Close()

	api := srvapi.ApiHandler{}
	api.DB = database.New(db)
	api.JwtSecret = jwtSecret

	mux := http.NewServeMux()
	mux.Handle("/", api)

	mux.Handle("POST /app", api.MiddlewareValidateUser(http.HandlerFunc(api.HandlerPostMessage)))

	mux.HandleFunc("POST /api/users", api.HandlerCreateUser)
	mux.HandleFunc("POST /api/login", api.HandlerLogin)
	mux.Handle("GET /api/login", api.MiddlewareValidateUser(http.HandlerFunc(api.HandlerLoggedIn)))

	server := http.Server{
		Handler:           mux,
		Addr:              ":" + port,
		ReadHeaderTimeout: 5 * time.Second,
	}

	fmt.Printf("serving app on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
