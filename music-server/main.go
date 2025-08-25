package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	initializers "github.com/franzego/music-server/Initializers"
	routers "github.com/franzego/music-server/Routers"
)

//The initializers simply run before the main function. See them as requirements for smooth sailing of the main function

//func init() {

//initializers.Syncdbandmodels()
//}

func main() {
	//routers.RouterFunc()
	//rou := mux.NewRouter() //main router with gorilla mux
	routers.RouterFunc()

	initializers.Connectiontodb()
	initializers.Syncdbandmodels()

	s := &http.Server{
		Addr:           ":8080",
		Handler:        routers.RouterFunc(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	//starting a goroutine to start the server. This is so that the server is run concurrently while other functions are still carried out in main

	go func() {
		log.Println("Starting server on :8080")

		//s.ListenAndServe()
		if err := s.ListenAndServe(); err != nil {
			log.Fatalf("%v error has occured", err)
		}

	}()

	// a channel carrying event streams. This channel will manage signals coming from user keyboard like a ctrl+c
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM) //send a stop signal
	<-stop

	//context to handle shutdown gracefully

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown failed, %v", err)
	}

	log.Print("Server Shutdown Successfully!")
}
