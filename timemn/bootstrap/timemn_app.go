package bootstrap

import (
	"github.com/nhuongmh/cvfs/timemn/internal/database/postgresdb"
)

type T2EApplication struct {
	Env *Env
	DB  *postgresdb.DB
}

func Init() T2EApplication {
	app := &T2EApplication{}
	app.Env = NewEnv()

	// ctx := context.Background()
	// db, err := postgresdb.ConnectDB(ctx, app.Env.DBUser, app.Env.DBPassword, app.Env.DBHost, app.Env.DBPort, app.Env.DBName, app.Env.DBSchema)

	// if err != nil {
	// 	logger.Log.Fatal().Err(err).Msg("Failed connect database")
	// }

	// app.DB = db

	// err = db.Migrate()
	// if err != nil {
	// 	logger.Log.Fatal().Err(err).Msg("Failed migrate database")
	// }
	// logger.Log.Info().Msg("Successfully migrated database")

	return *app
}

func (app *T2EApplication) CloseDB() {
	app.DB.Close()
}
