package main

import "github.com/madhab452/api-simulator/internal"

func main() {
	opts := internal.Options{}
	srv, err := internal.New(opts)
	if err != nil {
		panic(err)
	}

	srv.Start()

	defer srv.ShutDown()
}
