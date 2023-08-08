package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/twitter-remake/user/model"
)

func (d *Dependency) GetUser(ctx context.Context, id string) (model.User, error) {
	query, args, err := buildUserReadPredicateQuery(sq.Eq{"id": id})
	if err != nil {
		return model.User{}, err
	}
	row := d.db.QueryRow(ctx, query, args...)
	return getUserFromRow(row)
}
