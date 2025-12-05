package postgres

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Store struct {
	logger *zerolog.Logger
	db *pgxpool.Pool
}

// PostgresConnection creates a new connection pool to the PostgreSQL database.
func New(ctx context.Context) *Store{
	logger := log.With().Str("store", "postgres").Logger()

	// Build the connection string
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSL_MODE"),
	)

	// Parse the connection string into a config
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		logger.Error().Err(err).Msg("failed to parse postgres config")
		return nil
	}

	// Configure pool settings
	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute
	config.HealthCheckPeriod = time.Minute

	// Create the connection pool with a timeout context
	// defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create postgres pool")
		return nil
	}

	// Verify connection with a ping
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		logger.Error().Err(err).Msg("failed to ping postgres")
		return nil
	}


	return &Store{logger: &logger, db: pool}
}

func (s *Store) Close() {
	s.logger.Info().Msg("closing postgres connection pool")

	if s.db != nil {
		s.db.Close()
		s.db = nil
	}
	
	s.logger.Info().Msg("postgres connection pool closed")
}

