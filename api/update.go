package api

import (
	"context"

	userpb "github.com/twitter-remake/user/proto/gen/go/user"
)

func (h *handler) Update(context.Context, *userpb.UpdateProfileRequest) (*userpb.User, error) {
	return nil, nil
}
