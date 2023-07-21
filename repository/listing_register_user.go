package repository

import (
	"context"
	"time"
)

// ListingRegisterUserParams represents the parameters to register a new user
type ListingRegisterUserParams struct {
	Name            string
	ScreenName      string
	Email           string
	ProfileImageURL string
	BirthDate       time.Time
}

// ListingRegisterUser saves a new user in the database
func (d *Dependency) ListingRegisterUser(ctx context.Context, params ListingRegisterUserParams) error {
	query := `INSERT INTO users (name, screen_name, email, birth_date) VALUES ($1, $2, $3, $4)`
	_, err := d.db.Exec(ctx, query,
		params.Name,
		params.ScreenName,
		params.Email,
		params.ProfileImageURL,
		params.BirthDate)
	if err != nil {
		return err
	}

	return nil
}
