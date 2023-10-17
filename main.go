package main

import "github.com/madhab452/api-simulator/internal"

func main() {
	opts := internal.Options{}
	srv := internal.New(opts)
	srv.Start()
	defer srv.ShutDown()
}
