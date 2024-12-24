package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/mightynerd/hit/db"
	"github.com/mightynerd/hit/discogs"
	"github.com/mightynerd/hit/web"
)

type Server struct {
	ctx    *context.Context
	db     *db.DB
	config *Config
	wg     *sync.WaitGroup
	server *http.Server
}

func (s *Server) ConnectToDb() {
	db, err := db.Connect(*s.ctx, s.config.PGConnectionString)
	if err != nil {
		log.Fatal("Failed to connect to db", err)
	}

	s.db = db
}

func main() {
	// Load config
	config := LoadConfig("config.json")

	// Migrate
	migrator := db.NewMigrator(config.PGConnectionString)
	migrator.Migrate()

	ctx := context.Background()

	router := mux.NewRouter()
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	server := &Server{
		ctx:    &ctx,
		config: config,
		server: httpServer,
	}

	// Connect to DB
	server.ConnectToDb()

	discogs := discogs.NewDiscogsConfig(config.DiscogsAPIKey)

	routes := web.NewWeb(server.db, server.config.ServiceUrl, server.config.SpotifyClientId, server.config.SpotifyClientSecret, discogs)

	router.HandleFunc("/login", routes.Login).Methods("GET")
	router.HandleFunc("/callback", routes.Callback).Methods("GET")
	router.Handle("/playlists", routes.AuthMiddleware(http.HandlerFunc(routes.CreatePlaylist))).Methods("POST")
	router.Handle("/games", routes.AuthMiddleware(http.HandlerFunc(routes.CreateGame))).Methods("POST")
	router.Handle("/games/{gameId}/next", routes.AuthMiddleware(http.HandlerFunc(routes.AdvanceGame))).Methods("POST")

	fmt.Println("Listening on :8080")
	err := httpServer.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
