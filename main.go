package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/aliSHAIKH001/api_go/internal/app"
	"github.com/aliSHAIKH001/api_go/internal/routes"
)

func main() {
	var port int
	// Parses the value of port flag form the terminal and stores in port integer var
	flag.IntVar(&port, "port", 8080, "go backend server port")
	flag.Parse()


	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}


	r := routes.SetupRoutes(app)
	server := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
		Handler: r,
		IdleTimeout: time.Minute,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.Logger.Printf("we are running our app on port %d\n", port)

	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Fatal(err)
	}
}
