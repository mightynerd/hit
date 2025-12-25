package web_playlists

import (
	"context"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
	"github.com/mightynerd/hit/db"
	"github.com/mightynerd/hit/library"
	"github.com/mightynerd/hit/spotify"
	"github.com/mightynerd/hit/web"
	"github.com/mightynerd/hit/web/web_utils"
)

type CreatePlaylistReq struct {
	Body CreatePlaylistData
}

type CreatePlaylistData struct {
	Name string `json:"name"`
	From struct {
		Source string `json:"source" enum:"spotify_playlist"`
		ID     string `json:"id"`
	} `json:"from"`
}

type CreatePlaylistResp struct{}

func createPlaylist(web *web.Web) func(ctx context.Context, data *CreatePlaylistReq) (*CreatePlaylistResp, error) {
	return func(ctx context.Context, data *CreatePlaylistReq) (*CreatePlaylistResp, error) {
		user, err := web_utils.GetUserFromContext(ctx)
		if err != nil {
			return nil, err
		}

		playlist := &db.Playlist{
			UserID: user.ID,
			Name:   data.Body.Name,
			Status: db.PlaylistStatusImporting,
		}

		playlistId, err := web.DB.CreatePlaylist(playlist)
		if err != nil {
			return nil, huma.Error500InternalServerError("could not create playlist")
		}

		go handleImport(web, user, data.Body, playlistId)

		return nil, nil
	}
}

func handleImport(web *web.Web, user *db.User, data CreatePlaylistData, playlistId string) error {
	if data.From.Source == "spotify_playlist" {
		spotify := spotify.FromUser(user)
		lib := library.NewLibrary(web.DB, spotify, web.Discogs)
		err := lib.ImportSpotifyPlaylist(playlistId, data.From.ID)
		if err != nil {
			fmt.Println(err)
			web.DB.UpdatePlaylistStatus(playlistId, db.PlaylistStatusFailed)
			return err
		}
	}

	web.DB.UpdatePlaylistStatus(playlistId, db.PlaylistStatusActive)
	return nil
}
