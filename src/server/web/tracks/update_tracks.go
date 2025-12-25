package web_tracks

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/mightynerd/hit/web"
	"github.com/mightynerd/hit/web/web_utils"
)

type UpdateTrackReq struct {
	PlaylistID string `path:"playlistId"`
	TrackID    string `path:"trackId"`
	Body       struct {
		Title  *string `json:"title" required:"false"`
		Artist *string `json:"artist" required:"false"`
		Year   *int    `json:"year" required:"false"`
	}
}

type UpdateTrackResp struct {
	Body Track
}

func updateTrack(web *web.Web) func(ctx context.Context, data *UpdateTrackReq) (*UpdateTrackResp, error) {
	return func(ctx context.Context, data *UpdateTrackReq) (*UpdateTrackResp, error) {
		user, err := web_utils.GetUserFromContext(ctx)
		if err != nil {
			return nil, err
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

		if data.Body.Title != nil {
			track.Title = *data.Body.Title
		}

		if data.Body.Artist != nil {
			track.Artist = *data.Body.Artist
		}

		if data.Body.Year != nil {
			track.Year = *data.Body.Year
		}

		updated, err := web.DB.UpdateTrack(track)
		if err != nil {
			return nil, huma.Error500InternalServerError("could not update track")
		}

		return &UpdateTrackResp{
			Body: TrackFromDBTrack(*updated),
		}, nil
	}
}
