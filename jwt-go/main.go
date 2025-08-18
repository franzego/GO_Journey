package main

import (
	initializers "github.com/franzego/jwt-go/Initializers"
	"github.com/franzego/jwt-go/router"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.SyncDatabase()
}

func main() {

	router.RouterFunc()

	/*s := &http{
		Addr:           ":8080",
		Handler:        router.RouterFunc(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())*/

}
