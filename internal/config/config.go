package config

import (
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	ServicePort string `yaml:"service_port"`

	PGConn *pgxpool.Config `yaml:"pg_conn"`
}

func InitConfig() (*Config, error) {
	vp := viper.New()

	vp.AddConfigPath("config")
	vp.SetConfigName("config")

	if err := vp.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "failed to read in config")
	}

	var config Config

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?pool_max_conns=%d",
		vp.Get("config.db_conn_config.user"),
		vp.Get("config.db_conn_config.password"),
		vp.Get("config.db_conn_config.host"),
		vp.Get("config.db_conn_config.port"),
		vp.Get("config.db_conn_config.db_name"),
		vp.Get("config.db_conn_config.pool_max_conns"))

	pgxConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, errors.Wrap(err, "`Init config` failed to parse config")
	}

	config.ServicePort = vp.GetString("config.service_port")
	config.PGConn = pgxConfig

	return &config, nil
}
