package backend

import (
	"context"
	"errors"
	"time"

	"github.com/twitter-remake/user/model"
	"github.com/twitter-remake/user/repository"
)

// RegisterInput represents the parameters to register a new user
type RegisterInput struct {
	UUID            string
	Name            string
	ScreenName      string
	Email           string
	ProfileImageURL string
	BirthDate       time.Time
}

// Register saves a new user in the database
func (d *Dependency) Register(ctx context.Context, input RegisterInput) (model.User, error) {
	_, err := d.repo.GetUser(ctx, input.UUID)
	if err == nil && !errors.Is(err, repository.ErrNoUser) {
		return model.User{}, repository.ErrUserAlreadyExists
	}

	id, err := d.repo.SaveUser(ctx, repository.SaveUserParams{
		UUID:            input.UUID,
		Name:            input.Name,
		ScreenName:      input.ScreenName,
		Email:           input.Email,
		ProfileImageURL: input.ProfileImageURL,
		BirthDate:       input.BirthDate,
	})
	if err != nil {
		return model.User{}, err
	}

	newUser, err := d.repo.GetUser(ctx, id)
	if err != nil {
		return model.User{}, err
	}

	return newUser, nil
}
