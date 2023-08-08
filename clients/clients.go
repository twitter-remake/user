package clients

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"github.com/twitter-remake/user/config"
	"golang.org/x/sync/errgroup"
)

type Clients struct {
	ctx context.Context

	PostgreSQL      *pgxpool.Pool
	S3              *S3
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
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String(config.AWSRegion()),
			Credentials: credentials.NewStaticCredentials(
				config.AWSAccessKeyID(),
				config.AWSSecretAccessKey(),
				config.AWSSessionToken(),
			),
		})
		if err != nil {
			return err
		}

		c.S3 = NewS3(sess, config.S3Bucket())
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
