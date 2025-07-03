package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aliSHAIKH001/api_go/internal/api"
	"github.com/aliSHAIKH001/api_go/internal/store"
	"github.com/aliSHAIKH001/api_go/migrations"
)

type Application struct {
	Logger *log.Logger
	WorkoutHandler *api.WorkoutHandler
	DB *sql.DB
}

func NewApplication() (*Application, error) {
	// Opens a connection to the docker database.
	pgDB, err := store.Open()
	if err != nil {
		return nil, err
	}

	err = store.MigrateFS(pgDB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// Bitwise OR | operator telling logger to combine both date and time
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// Stores will belong here

	// Handlers will belong here
	workoutHandler := api.NewWorkoutHandler()

	app := &Application{
		Logger: logger,
		WorkoutHandler: workoutHandler,
		DB: pgDB,
	}

	return app, nil
}


func (a *Application) HealthCheck(w http.ResponseWriter,  r *http.Request) {
	fmt.Fprintf(w, "Status is available\n")
}