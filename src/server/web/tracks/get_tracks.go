package web_tracks

import (
	"context"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/mightynerd/hit/db"
	"github.com/mightynerd/hit/web"
	"github.com/mightynerd/hit/web/web_utils"
)

type GetTracksReq struct {
	PlaylistID string `path:"playlistId"`
	Page       int    `query:"page" default:"0"`
	Size       int    `query:"size" default:"20"`
}

type GetTracksResp struct {
	Body getTracksRespData
}

type getTracksRespData struct {
	Data  []Track `json:"data"`
	Total int     `json:"total"`
}

type Track struct {
	ID         string    `json:"id" format:"uuid"`
	CreatedAt  time.Time `json:"createdAt"`
	PlaylistID string    `json:"playlistId" format:"uuid"`
	Title      string    `json:"title"`
	Artist     string    `json:"artist"`
	Year       int       `json:"year"`
	SpotifyURI string    `json:"spotifyURI"`
}

func TrackFromDBTrack(dbTrack db.Track) Track {
	return Track{
		ID:         dbTrack.ID,
		CreatedAt:  dbTrack.CreatedAt,
		PlaylistID: dbTrack.PlaylistID,
		Title:      dbTrack.Title,
		Artist:     dbTrack.Artist,
		Year:       dbTrack.Year,
		SpotifyURI: dbTrack.SpotifyURI,
	}
}

func getTracks(web *web.Web) func(ctx context.Context, data *GetTracksReq) (*GetTracksResp, error) {
	return func(ctx context.Context, data *GetTracksReq) (*GetTracksResp, error) {
		user, err := web_utils.GetUserFromContext(ctx)
		if err != nil {
			return nil, err
		}

		playlist, err := web.DB.GetPlaylistById(data.PlaylistID)
		if err != nil {
			return nil, huma.Error404NotFound("playlist not found")
		}

		if playlist.UserID != user.ID {
			return nil, huma.Error404NotFound("playlist not found")
		}

		dbTracks, err := web.DB.GetTracks(playlist.ID, data.Page, data.Size)
		if err != nil {
			return nil, huma.Error500InternalServerError("could not get tracks")
		}

		tracks := []Track{}
		for _, dbTrack := range *dbTracks {
			tracks = append(tracks, TrackFromDBTrack(dbTrack))
		}

		return &GetTracksResp{
			Body: getTracksRespData{
				Data:  tracks,
				Total: 0,
			},
		}, nil
	}
}
