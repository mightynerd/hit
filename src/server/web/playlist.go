package web

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mightynerd/hit/db"
	"github.com/mightynerd/hit/library"
	"github.com/mightynerd/hit/spotify"
)

type CreatePlaylistBody struct {
	Name string `json:"name"`
	From struct {
		Source string `json:"source"`
		Id     string `json:"id"`
	} `json:"from"`
}

func (web *Web) handleImport(user *db.User, body *CreatePlaylistBody, playlistId string) error {
	if body.From.Source == "spotify_playlist" {
		spotify := spotify.FromUser(user)
		lib := library.NewLibrary(web.db, spotify, web.discogs)
		err := lib.ImportSpotifyPlaylist(playlistId, body.From.Id)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}

func (web *Web) CreatePlaylist(w http.ResponseWriter, r *http.Request) {
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

	var body CreatePlaylistBody
	err = json.Unmarshal(rawBody, &body)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "failed to parse body", http.StatusBadRequest)
		return
	}

	user := r.Context().Value(userContextKey).(*db.User)

	playlist := &db.Playlist{
		UserID: user.ID,
		Name:   body.Name,
	}

	playlistId, err := web.db.CreatePlaylist(playlist)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "could not create paylist", http.StatusInternalServerError)
		return
	}

	err = web.handleImport(user, &body, playlistId)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "could not import playlist", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"playlist_id": playlistId,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
