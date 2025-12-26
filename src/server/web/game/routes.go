package web_games

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/mightynerd/hit/web"
)

func RegisterRoutes(web *web.Web) {
	group := huma.NewGroup(web.API, "/games")
	group.UseMiddleware(web.AuthMiddleware)
	tags := []string{"Games"}

	huma.Register(group,
		huma.Operation{
			Method: http.MethodPost,
			Path:   "",
			Tags:   tags,
			Security: []map[string][]string{
				{"bearerAuth": {}},
			},
		}, createGame(web))

	huma.Register(group,
		huma.Operation{
			Method: http.MethodPost,
			Path:   "/:gameId/advance",
			Tags:   tags,
			Security: []map[string][]string{
				{"bearerAuth": {}},
			},
		}, advanceGame(web))
}
