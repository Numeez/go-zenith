package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	app "github.com/Numeez/go-zenith/internal/app"
	router "github.com/Numeez/go-zenith/internal/routes"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "This is the port used to host the server")
	flag.Parse()
	application, err := app.NewApplication()
	if err != nil {
		panic(err)
	}
	defer application.DB.Close()
	http.HandleFunc("/health", application.HealthCheck)
	r := router.SetupRoutes(application)
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
	}
	application.Logger.Printf("Server is running on port: %d\n", port)

	if err := server.ListenAndServe(); err != nil {
		application.Logger.Fatal(err)
	}

}
