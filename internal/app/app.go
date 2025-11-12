package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Numeez/go-zenith/internal/api"
	"github.com/Numeez/go-zenith/internal/store"
	"github.com/Numeez/go-zenith/migrations"
)

type Application struct {
	Logger         *log.Logger
	WorkOutHandler *api.WorkOutHandler
	DB             *sql.DB
}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	db, err := store.Open()
	if err != nil {
		return nil, err
	}

	err = store.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		return nil, err
	}
	workoutStore := store.NewPostgresWorkoutStore(db)
	workOutHandler := api.NewWorkOutHandler(workoutStore,logger)
	return &Application{
		Logger:         logger,
		WorkOutHandler: workOutHandler,
		DB:             db, 
	}, nil
}

func (app *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Server is running\n")
}
