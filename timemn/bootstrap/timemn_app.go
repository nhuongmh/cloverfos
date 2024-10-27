package bootstrap

import (
	"context"

	"github.com/nhuongmh/cvfs/timemn/internal/database/postgresdb"
	"github.com/nhuongmh/cvfs/timemn/internal/logger"
)

type TimeMnApplication struct {
	Env *Env
	DB  *postgresdb.DB
}

func Init() TimeMnApplication {
	app := &TimeMnApplication{}
	app.Env = NewEnv()

	ctx := context.Background()
	db, err := postgresdb.ConnectDB(ctx, app.Env.DBUser, app.Env.DBPassword, app.Env.DBHost, app.Env.DBPort, app.Env.DBName, app.Env.DBSchema)

	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed connect database")
	}

	app.DB = db

	err = db.Migrate()
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed migrate database")
	}
	logger.Log.Info().Msg("Successfully migrated database")

	return *app
}

func (app *TimeMnApplication) CloseDB() {
	app.DB.Close()
}
