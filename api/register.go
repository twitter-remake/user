package api

import (
	"context"

	userpb "github.com/twitter-remake/user/proto/gen/go/user"
)

func (h *handler) Register(context.Context, *userpb.RegisterRequest) (*userpb.RegisterResponse, error) {
	return nil, nil
}
