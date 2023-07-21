package api

import (
	"context"

	userpb "github.com/twitter-remake/user/proto/gen/go/user"
)

func (h *handler) GetPublicProfile(context.Context, *userpb.GetPublicProfileRequest) (*userpb.User, error) {
	return nil, nil
}
