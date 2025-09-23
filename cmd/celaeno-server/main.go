package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/seandisero/celaeno/internal/server/chat"
	"github.com/seandisero/celaeno/internal/server/database"
	"github.com/seandisero/celaeno/internal/server/srvapi"

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
	}

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("libsql", dbURL)
	if err != nil {
		log.Fatal("could not open db")
		os.Exit(1)
	}
	defer db.Close()

	chatServer := chat.NewChatServer()

	api := srvapi.ApiHandler{}
	api.DB = database.New(db)
	api.ChatService = chatServer
	api.JwtSecret = jwtSecret

	mux := http.NewServeMux()
	mux.Handle("/", api)

	// TODO: chage api/chat to chat/api since I'll probably be routing from my site to multiple programs.
	mux.HandleFunc("GET /startup", api.HandlerStartup)

	mux.HandleFunc("POST /api/users", api.HandlerCreateUser)
	mux.HandleFunc("PUT /api/users/{id}", api.MiddlewareValidateUser(api.HandlerSetDisplayName))
	mux.HandleFunc("DELETE /api/users/{id}", api.MiddlewareValidateUser(api.HandlerDeleteUser))

	mux.HandleFunc("POST /api/login", api.HandlerLogin)
	mux.HandleFunc("GET /api/login", api.MiddlewareValidateUser(api.HandlerLoggedIn))

	mux.HandleFunc("/api/chat/create", api.MiddlewareValidateUser(api.HandlerCreateChat))
	mux.HandleFunc("/api/chat/connect/{name}", api.MiddlewareValidateUser(api.HandlerConnectToChat))

	mux.HandleFunc("POST /api/chat/publish/{name}", api.MiddlewareValidateUser(api.HandlerPostMessage))

	mux.HandleFunc("GET /status", api.HandlerStatus)

	server := http.Server{
		Handler:           mux,
		Addr:              ":" + port,
		ReadHeaderTimeout: 5 * time.Second,
	}

	fmt.Printf("serving app on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
