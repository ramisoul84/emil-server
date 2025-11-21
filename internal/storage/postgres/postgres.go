package postgres

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"gitlab.com/ramisoul/emil-server/config"
)

func New(cfg config.DatabaseConfig) (*sqlx.DB, error) {
	driverName := "postgres"
	dataSourceName := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DB, cfg.SSLMode)
	db, err := sqlx.Connect(driverName, dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}

	maxOpenConns := 25
	if cfg.MaxOpenConns > 0 {
		maxOpenConns = cfg.MaxOpenConns
	}

	maxIdleConns := 10
	if cfg.MaxIdleConns > 0 {
		maxIdleConns = cfg.MaxIdleConns
	}

	maxLifetime := 5 * time.Minute
	if cfg.ConnMaxLifetime > 0 {
		maxLifetime = cfg.ConnMaxLifetime
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(maxLifetime)

	return db, nil
}
