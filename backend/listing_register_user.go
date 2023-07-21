package backend

import (
	"context"
	"io"
	"path"
	"time"

	"github.com/google/uuid"
	"github.com/twitter-remake/user/repository"
)

// ListingRegisterUserInput represents the parameters to register a new user
type ListingRegisterUserInput struct {
	Name               string
	ScreenName         string
	Email              string
	ProfileImage       io.Reader
	ProfileImageFormat string
	BirthDate          time.Time
}

// ListingRegisterUser saves a new user in the database
func (d *Dependency) ListingRegisterUser(ctx context.Context, input ListingRegisterUserInput) error {
	profileImageURL, err := d.saveProfileImage(ctx, input)
	if err != nil {
		return err
	}

	err = d.repo.ListingRegisterUser(ctx, repository.ListingRegisterUserParams{
		Name:            input.Name,
		ScreenName:      input.ScreenName,
		Email:           input.Email,
		ProfileImageURL: profileImageURL,
		BirthDate:       input.BirthDate,
	})
	if err != nil {
		return err
	}

	return nil
}

func (d *Dependency) saveProfileImage(ctx context.Context, input ListingRegisterUserInput) (string, error) {
	key := path.Join("users", uuid.New().String()+input.ProfileImageFormat)

	uploadOutput, err := d.clients.S3.UploadFile(ctx, key, input.ProfileImage)
	if err != nil {
		return "", err
	}

	return uploadOutput.Location, nil
}
