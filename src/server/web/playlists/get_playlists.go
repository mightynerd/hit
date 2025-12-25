package web_playlists

import (
	"context"
	"time"

	"github.com/mightynerd/hit/db"
	"github.com/mightynerd/hit/web"
	"github.com/mightynerd/hit/web/web_utils"
)

type GetPlaylistsReq struct {
	Page int `query:"page" default:"0"`
	Size int `query:"size" default:"20"`
}

type GetPlaylistsResp struct {
	Body getPlaylistsRespData
}

type getPlaylistsRespData struct {
	Data  []Playlist `json:"data"`
	Total int        `json:"total"`
}

type Playlist struct {
	ID        string    `json:"id" format:"uuid"`
	CreatedAt time.Time `json:"createdAt"`
	Name      string    `json:"name"`
	UserID    string    `json:"userId"`
	Status    string    `json:"status" enum:"active,importing,failed"`
}

func PlaylistFromDB(dbPlaylist db.Playlist) Playlist {
	return Playlist{
		ID:        dbPlaylist.ID,
		CreatedAt: dbPlaylist.CreatedAt,
		Name:      dbPlaylist.Name,
		UserID:    dbPlaylist.UserID,
		Status:    string(dbPlaylist.Status),
	}
}

func getPlaylists(web *web.Web) func(ctx context.Context, data *GetPlaylistsReq) (*GetPlaylistsResp, error) {
	return func(ctx context.Context, data *GetPlaylistsReq) (*GetPlaylistsResp, error) {
		user, err := web_utils.GetUserFromContext(ctx)
		if err != nil {
			return nil, err
		}

		dbPlaylists, err := web.DB.GetPlaylists(user.ID, data.Page, data.Size)
		if err != nil {
			return nil, err
		}

		playlists := []Playlist{}
		for _, dbPlaylist := range *dbPlaylists {
			playlists = append(playlists, PlaylistFromDB(dbPlaylist))
		}

		return &GetPlaylistsResp{
			Body: getPlaylistsRespData{
				Data:  playlists,
				Total: 0,
			},
		}, nil
	}
}
