package api

import (
	"context"

	userpb "github.com/twitter-remake/user/proto/gen/go/user"
)

func (h *handler) DeleteUser(context.Context, *userpb.DeleteUserRequest) (*userpb.DeleteUserResponse, error) {
	return nil, nil
}
