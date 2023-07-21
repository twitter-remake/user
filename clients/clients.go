package clients

import (
	"context"
	"fmt"
	"net"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"github.com/twitter-remake/user/config"
	"golang.org/x/sync/errgroup"
)

type Clients struct {
	ctx context.Context

	PostgreSQL      *pgxpool.Pool
	S3              *S3
	ServiceRegistry *consulapi.Client
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

		if err := c.registerServiceToConsul(); err != nil {
			return errors.Wrap(err, "initializing consul")
		}
		return nil
	})

	if err := group.Wait(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Clients) registerServiceToConsul() error {
	var address string
	if config.Environment() == "dev" {
		// assuming consul is running in docker
		address = "host.docker.internal"
	} else {
		address = config.Host()
	}

	check := &consulapi.AgentServiceCheck{
		HTTP:                           fmt.Sprintf("http://%s/", net.JoinHostPort(address, config.Port())),
		Interval:                       "10s",
		Timeout:                        "30s",
		CheckID:                        fmt.Sprintf("service:%s:http", config.AppName()),
		DeregisterCriticalServiceAfter: "1m",
		TLSSkipVerify:                  func() bool { return config.Environment() == "dev" }(),
	}

	port, _ := strconv.Atoi(config.Port())

	serviceDefinition := &consulapi.AgentServiceRegistration{
		ID:      config.AppName(),
		Name:    config.AppName() + "_master",
		Port:    port,
		Address: address,
		Tags:    []string{config.Environment(), config.AppName()},
		Check:   check,
	}

	if err := c.ServiceRegistry.Agent().ServiceRegister(serviceDefinition); err != nil {
		return err
	}

	return nil
}
