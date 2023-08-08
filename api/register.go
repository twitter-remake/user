package api

import (
	"context"
	"errors"
	"net/mail"
	"net/url"

	"github.com/twitchtv/twirp"
	"github.com/twitter-remake/user/backend"
	userpb "github.com/twitter-remake/user/proto/gen/go/user"
	"github.com/twitter-remake/user/repository"
)

func (h *handler) Register(ctx context.Context, req *userpb.RegisterRequest) (*userpb.RegisterResponse, error) {
	if err := validateRegisterRequest(req); err != nil {
		return nil, err
	}

	newUser, err := h.backend.Register(ctx, backend.RegisterInput{
		UUID:            req.GetUuid(),
		Name:            req.GetName(),
		ScreenName:      req.GetScreenName(),
		Email:           req.GetEmail(),
		ProfileImageURL: req.GetProfileImageUrl(),
		BirthDate:       req.GetBirthDate().AsTime(),
	})
	switch {
	case err != nil && errors.Is(err, repository.ErrUserAlreadyExists):
		return nil, twirp.AlreadyExists.Error(err.Error())
	case err != nil:
		return nil, err
	}

	return &userpb.RegisterResponse{
		User: newUser.PB(),
	}, nil
}

func validateRegisterRequest(req *userpb.RegisterRequest) error {
	if req.GetUuid() == "" {
		return twirp.RequiredArgumentError("uuid")
	}

	if req.GetName() == "" {
		return twirp.RequiredArgumentError("name")
	}

	if req.GetScreenName() == "" {
		return twirp.RequiredArgumentError("screen_name")
	}

	if req.GetEmail() == "" {
		return twirp.RequiredArgumentError("email")
	}

	if _, err := mail.ParseAddress(req.GetEmail()); err != nil {
		return twirp.InvalidArgumentError("email", "invalid email address")
	}

	if !req.GetBirthDate().IsValid() {
		return twirp.RequiredArgumentError("birth_date")
	}

	if req.GetProfileImageUrl() == "" {
		return twirp.RequiredArgumentError("email")
	}

	profileImageURL, err := url.Parse(req.GetProfileImageUrl())
	if err != nil && profileImageURL.Scheme == "" && profileImageURL.Host == "" {
		return twirp.InvalidArgumentError("profile_image_url", "invalid profile image url")
	}

	return nil
}
