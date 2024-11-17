package postgresdb

import (
	"context"
	"embed"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nhuongmh/cvfs/timemn/internal/logger"
	"github.com/pkg/errors"

	_ "github.com/jackc/pgx/v5/stdlib"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
)

//go:embed migrations/*.sql
var migrationFs embed.FS

type DB struct {
	*pgxpool.Pool
	QueryBuilder *squirrel.StatementBuilderType
	url          string
}

func ConnectDB(ctx context.Context, user, pass, host, port, dbname, schema string) (*DB, error) {
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", user, pass, host, port, dbname, schema)

	logger.Log.Info().Msgf("Connecting to database %v", url)
	db, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, errors.Wrap(err, "Failed connect to Database")
	}
	logger.Log.Info().Msgf("Connected to database")

	var greeting string
	err = db.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		errors.Wrap(err, "Failed QueryRow in Database")
	}
	err = db.Ping(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Failed ping Database")
	}

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	return &DB{
		db,
		&psql,
		url,
	}, nil
}

func (db *DB) Migrate() error {
	driver, err := iofs.New(migrationFs, "migrations")
	if err != nil {
		return errors.Wrap(err, "Failed reading .sql migrations file")
	}

	migrations, err := migrate.NewWithSourceInstance("iofs", driver, db.url)
	if err != nil {
		return errors.Wrap(err, "Failed init migration instance")
	}

	err = migrations.Up()
	if err != nil && err != migrate.ErrNoChange {
		return errors.Wrap(err, "Failed migrate database")
	}

	return nil
}

func (db *DB) Close() {
	logger.Log.Info().Msgf("Disconnecting from database %v", db.url)
	db.Pool.Close()
	logger.Log.Info().Msg("Disconnected from database")
}
