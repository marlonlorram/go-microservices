package mongodb

import (
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"

	"github.com/marlonlorram/go-microservices/internal/pkg/config"
)

// DBConfig represents environment variable configuration parameters for database.
type DBConfig struct {
	Host string `envconfig:"DB_HOST" default:"0.0.0.0" required:"true" desc:"Hostname or IP of the database server"`
	Port int    `envconfig:"DB_PORT" default:"27017" required:"true" desc:"Port of the database server"`
	User string `envconfig:"DB_USER" required:"true" desc:"User for authentication in the database"`
	Pass string `envconfig:"DB_PASSWORD" required:"true" desc:"Password for authentication in the database"`
	Name string `envconfig:"DB_NAME" required:"true" desc:"Name of the database"`
	Mech string `envconfig:"DB_MECH" required:"false" desc:"Authentication mechanism of the database (optional)"`
}

// parseConfig parses the environment variables and returns a configured instance of DBConfig.
// Returns an error if the configuration could not be loaded.
func parseConfig() (*DBConfig, error) {
	cfg, err := config.Parse[DBConfig]()
	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("failed to parse database configuration"))
	}

	if cfg.Port < 1 || cfg.Port > 65535 {
		return nil, fault.New("invalid database port: it must be in the range 1-65535")
	}

	return cfg, nil
}
