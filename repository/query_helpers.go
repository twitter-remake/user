package repository

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/twitter-remake/user/model"
)

var (
	sqlUserPropertiesList = []string{"id",
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
		"updated_at",
	}
)

func buildUserReadPredicateQuery(predicate any) (string, []any, error) {
	sb := sq.
		Select(sqlUserPropertiesList...).
		From("users").
		Where(predicate).
		PlaceholderFormat(sq.Dollar)
	return sb.ToSql()
}

func getUserFromRow(row pgx.Row) (model.User, error) {
	var user model.User
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.ScreenName,
		&user.Email,
		&user.Bio,
		&user.Location,
		&user.Website,
		&user.BirthDate,
		&user.ProfileImageURL,
		&user.ProfileBannerURL,
		&user.FollowersCount,
		&user.FollowingsCount,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	switch {
	case err == pgx.ErrNoRows:
		return model.User{}, ErrNoUser
	case err != nil:
		return model.User{}, err
	}

	return user, nil
}
