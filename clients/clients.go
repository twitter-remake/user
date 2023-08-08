package clients

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"github.com/twitter-remake/user/config"
	"golang.org/x/sync/errgroup"
)

type Clients struct {
	ctx context.Context

	PostgreSQL      *pgxpool.Pool
	ServiceRegistry *Consul
}

func New(ctx context.Context) (*Clients, error) {
	c := new(Clients)
	c.ctx = ctx

	var group errgroup.Group

	group.Go(func() error {
		var err error
		c.PostgreSQL, err = NewPostgreSQLClient(ctx, config.DatabaseURL())
		if err != nil {
			return errors.Wrap(err, "initializing postgresql")
		}
		return nil
	})

	group.Go(func() error {
		var err error
		c.ServiceRegistry, err = NewConsulAPI()
		if err != nil {
			return errors.Wrap(err, "initializing consul")
		}

		if err := c.ServiceRegistry.Register(&RegisterCfg{
			ID:          config.AppName(),
			Host:        config.Host(),
			Port:        config.Port(),
			Environment: config.Environment(),
		}); err != nil {
			return errors.Wrap(err, "initializing consul")
		}
		return nil
	})

	if err := group.Wait(); err != nil {
		return nil, err
	}

	return c, nil
}
