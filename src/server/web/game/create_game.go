package web_games

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/mightynerd/hit/db"
	"github.com/mightynerd/hit/web"
	"github.com/mightynerd/hit/web/web_utils"
)

type CreateGameReq struct {
	Body struct {
		PlaylistID string `json:"playlistId" format:"uuid"`
	}
}

type CreateGameResp struct {
	Body CreateGameRespData
}

type CreateGameRespData struct {
	GameID string `json:"gameId" format:"uuid"`
}

func createGame(web *web.Web) func(ctx context.Context, data *CreateGameReq) (*CreateGameResp, error) {
	return func(ctx context.Context, data *CreateGameReq) (*CreateGameResp, error) {
		user, err := web_utils.GetUserFromContext(ctx)
		if err != nil {
			return nil, huma.Error500InternalServerError("missing user")
		}

		playlist, err := web.DB.GetPlaylistById(data.Body.PlaylistID)
		if err != nil {
			return nil, huma.Error404NotFound("playlist not found")
		}

		if playlist.UserID != user.ID {
			return nil, huma.Error404NotFound("playlist not found")
		}

		game := &db.Game{
			UserID:     user.ID,
			PlaylistID: playlist.ID,
		}

		gameId, err := web.DB.CreateGame(game)
		if err != nil {
			return nil, huma.Error500InternalServerError("could not create game")
		}

		return &CreateGameResp{
			Body: CreateGameRespData{
				GameID: gameId,
			},
		}, nil
	}
}
