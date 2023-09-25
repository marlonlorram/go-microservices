package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/fx"
)

const connectionTimeout = 10 * time.Second

// newClient creates and returns a new MongoDB client using the provided configuration.
// It utilizes the provided DBConfig to establish a connection to the MongoDB server
// and returns an error if the connection attempt fails.
func newClient(cfg *DBConfig) (*mongo.Client, error) {
	credential := options.Credential{
		Username:      cfg.User,
		Password:      cfg.Pass,
		AuthSource:    cfg.Name,
		AuthMechanism: cfg.Mech,
	}

	uri := fmt.Sprintf("mongodb://%s:%d", cfg.Host, cfg.Port)
	clientOpts := options.Client().ApplyURI(uri).SetAuth(credential)

	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to connect to the database"))
	}

	return client, nil
}

// newDatabase returns a MongoDB database instance from the provided client and configuration.
func newDatabase(client *mongo.Client, cfg *DBConfig) (*mongo.Database, error) {
	return client.Database(cfg.Name), nil
}

// mount appends lifecycle hooks to the fx Lifecycle for starting and stopping
// the MongoDB client. It ensures that the client establishes a connection
// on start and properly disconnects on stop, returning any errors encountered.
func mount(lc fx.Lifecycle, client *mongo.Client) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := client.Ping(ctx, readpref.Primary()); err != nil {
				return fault.Wrap(err, fmsg.With("failed to verify connection to the database"))
			}

			return nil
		},

		OnStop: func(ctx context.Context) error {
			if err := client.Disconnect(ctx); err != nil {
				return fault.Wrap(err, fmsg.With("failed to disconnect from the database"))
			}

			return nil
		},
	})
}

func Build() fx.Option {
	return fx.Options(
		fx.Provide(
			parseConfig,
			newClient,
			newDatabase,
		),
		fx.Invoke(mount),
	)
}
