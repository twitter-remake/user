package api

import (
	"context"

	userpb "github.com/twitter-remake/user/proto/gen/go/user"
)

func (h *handler) Get(context.Context, *userpb.GetProfileRequest) (*userpb.User, error) {
	return nil, nil
}
