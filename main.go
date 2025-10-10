package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/makhammatovb/femProject/internal/routes"
	"github.com/makhammatovb/femProject/internal/app"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "Port to run the server on")
	flag.Parse()
	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}

	r := routes.SetupRoutes(app)
	server := &http.Server{
		Addr:       fmt.Sprintf(":%d", port),
		Handler:   r,
		IdleTimeout: time.Minute,
		ReadTimeout: time.Second * 10,
		WriteTimeout: time.Second * 30,
	}

	app.Logger.Printf("Listening on port %d", port)

	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Fatal(err)
	}
}

