package web_games

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/mightynerd/hit/game"
	"github.com/mightynerd/hit/spotify"
	"github.com/mightynerd/hit/web"
	web_tracks "github.com/mightynerd/hit/web/tracks"
	"github.com/mightynerd/hit/web/web_utils"
)

type AdvanceGameReq struct {
	GameID string `path:"gameId" format:"uuid"`
}

type AdvanceGameResp struct {
	Body web_tracks.Track
}

func advanceGame(web *web.Web) func(ctx context.Context, data *AdvanceGameReq) (*AdvanceGameResp, error) {
	return func(ctx context.Context, data *AdvanceGameReq) (*AdvanceGameResp, error) {
		user, err := web_utils.GetUserFromContext(ctx)
		if err != nil {
			return nil, huma.Error500InternalServerError("missing user")
		}

		dbGame, err := web.DB.GetGameById(data.GameID)
		if err != nil {
			return nil, huma.Error404NotFound("game not found")
		}

		if dbGame.UserID != user.ID {
			return nil, huma.Error404NotFound("game not found")
		}

		spotify := spotify.FromUser(user)
		game := game.NewGame(spotify, dbGame, web.DB)

		track, err := game.Advance()

		if err != nil {
			return nil, huma.Error500InternalServerError("could not advance game")
		}

		return &AdvanceGameResp{
			Body: web_tracks.TrackFromDBTrack(*track),
		}, nil
	}
}
