package main

import "github.com/madhab452/api-simulator/internal"

func main() {
	opts := internal.Options{
		ServiceMap: []internal.ServiceMap{
			{
				Address: ":1949",
				Name:    "blog",
			},
		},
	}
	srv := internal.New(opts)
	defer srv.ShutDown()
}
