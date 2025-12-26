package web_playlists

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/mightynerd/hit/web"
	"github.com/mightynerd/hit/web/web_utils"
)

type DeletePlaylistReq struct {
	PlaylistID string `path:"playlistId" format:"uuid"`
}

type DeletePlaylistResp struct{}

func deletePlaylist(web *web.Web) func(ctx context.Context, data *DeletePlaylistReq) (*DeletePlaylistResp, error) {
	return func(ctx context.Context, data *DeletePlaylistReq) (*DeletePlaylistResp, error) {
		user, err := web_utils.GetUserFromContext(ctx)
		if err != nil {
			return nil, err
		}

		playlist, err := web.DB.GetPlaylistById(data.PlaylistID)
		if err != nil {
			return nil, huma.Error404NotFound("playlist not found")
		}

		if user.ID != playlist.UserID {
			return nil, huma.Error404NotFound("playlist not found")
		}

		err = web.DB.DeletePlaylist(data.PlaylistID)
		if err != nil {
			return nil, huma.Error500InternalServerError("could not delete playlist")
		}

		return nil, nil
	}
}
