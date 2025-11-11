package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Numeez/go-zenith/internal/api"
)

type Application struct {
	Logger         *log.Logger
	WorkOutHandler *api.WorkOutHandler
}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	workOutHandler := api.NewWorkOutHandler()
	return &Application{
		Logger:         logger,
		WorkOutHandler: workOutHandler,
	}, nil
}

func (app *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Server is running\n")
}
