package web_tracks

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/mightynerd/hit/web"
	"github.com/mightynerd/hit/web/web_utils"
)

type DeleteTrackReq struct {
	PlaylistID string `path:"playlistId"`
	TrackID    string `path:"trackId"`
}

type DeleteTrackResp struct{}

func deleteTrack(web *web.Web) func(ctx context.Context, data *DeleteTrackReq) (*DeleteTrackResp, error) {
	return func(ctx context.Context, data *DeleteTrackReq) (*DeleteTrackResp, error) {
		user, err := web_utils.GetUserFromContext(ctx)
		if err != nil {
			return nil, huma.Error500InternalServerError("missing user")
		}

		track, err := web.DB.GetTrackById(data.TrackID)
		if err != nil {
			return nil, huma.Error404NotFound("track not found")
		}

		playlist, err := web.DB.GetPlaylistById(track.PlaylistID)
		if err != nil {
			return nil, huma.Error500InternalServerError("could not get playlist")
		}

		if user.ID != playlist.UserID {
			return nil, huma.Error404NotFound("track not found")
		}

		err = web.DB.DeleteTrack(data.TrackID)
		if err != nil {
			return nil, huma.Error500InternalServerError("could not delete track")
		}

		return nil, nil
	}
}
