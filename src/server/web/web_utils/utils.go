package web_utils

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/mightynerd/hit/db"
)

func GetUserFromContext(ctx context.Context) (*db.User, error) {
	user, exists := ctx.Value("user").(*db.User)
	if !exists {
		return nil, huma.Error500InternalServerError("missing user")
	}

	return user, nil
}
