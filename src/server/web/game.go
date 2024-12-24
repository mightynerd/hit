package web

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mightynerd/hit/db"
	"github.com/mightynerd/hit/game"
	"github.com/mightynerd/hit/spotify"
)

type CreateGameBody struct {
	PlaylistId string `json:"playlist_id"`
}

func (web *Web) CreateGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}

	defer r.Body.Close()
	rawBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "could not read body", http.StatusInternalServerError)
		return
	}

	var body CreateGameBody
	err = json.Unmarshal(rawBody, &body)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "failed to parse body", http.StatusBadRequest)
		return
	}

	user := r.Context().Value(userContextKey).(*db.User)

	playlist, err := web.db.GetPlaylistById(body.PlaylistId)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "failed to get playlist", http.StatusInternalServerError)
		return
	}

	game := &db.Game{
		UserID:     user.ID,
		PlaylistID: playlist.ID,
	}

	gameId, err := web.db.CreateGame(game)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "could not create game", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"game_id": gameId,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (web *Web) AdvanceGame(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameId := vars["gameId"]
	user := r.Context().Value(userContextKey).(*db.User)

	dbGame, err := web.db.GetGameById(gameId)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "failed to get game", http.StatusInternalServerError)
		return
	}

	if dbGame.UserID != user.ID {
		http.Error(w, "game not found", http.StatusNotFound)
		return
	}

	spotify := spotify.FromUser(user)
	game := game.NewGame(spotify, dbGame, web.db)

	track, err := game.Advance()

	if err != nil {
		fmt.Println(err)
		http.Error(w, "could not advance game", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(track)
}
