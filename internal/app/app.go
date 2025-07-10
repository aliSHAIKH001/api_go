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
	
	// Sets up the structure of our database with migration files.
	err = store.MigrateFS(pgDB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// Bitwise OR | operator telling logger to combine both date and time
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// Stores will belong here
	// This gives us a struct with the capability to perform crud operations on the database.
	workoutStore := store.NewPostgresWorkoutStore(pgDB)

	// Handlers will belong here
	// The struct above is used by the handler depending on the requests we recieve
	workoutHandler := api.NewWorkoutHandler(workoutStore, logger)

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