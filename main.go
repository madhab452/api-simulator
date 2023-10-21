package main

import "github.com/madhab452/api-simulator/internal"

func main() {
	opts := internal.Option{
		ListenAddrHTTP: ":1949", // todo: read from env var
	}
	srv, err := internal.New(opts)
	if err != nil {
		panic(err)
	}

	srv.Start()

	defer srv.ShutDown()
}
