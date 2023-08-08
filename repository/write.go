package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/rs/zerolog/log"
)

// SaveUserParams represents the parameters to register a new user
type SaveUserParams struct {
	UUID            string
	Name            string
	ScreenName      string
	Email           string
	ProfileImageURL string
	BirthDate       time.Time
}

// SaveUser saves a new user in the database
func (d *Dependency) SaveUser(ctx context.Context, params SaveUserParams) (string, error) {
	log.Debug().Any("params", params).Msg("SaveUser")
	sb := sq.Insert("users").
		Columns(
			"id",
			"name",
			"screen_name",
			"email",
			"bio",
			"location",
			"website",
			"birth_date",
			"profile_image_url",
			"profile_banner_url",
			"followers_count",
			"followings_count",
			"created_at",
			"updated_at").
		Values(
			params.UUID,
			params.Name,
			params.ScreenName,
			params.Email,
			"",
			"",
			"",
			params.BirthDate,
			params.ProfileImageURL,
			"",
			0,
			0,
			time.Now(),
			time.Now()).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING id")

	query, args, err := sb.ToSql()
	if err != nil {
		return "", err
	}

	var id string
	err = d.db.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}
