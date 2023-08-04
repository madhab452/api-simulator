package internal

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"
)

type ServiceMap struct {
	Address string
	Name    string
}

type Options struct {
	ServiceMap []ServiceMap
}

type Route struct {
	Path      string              `yaml:"path"`
	Method    string              `yaml:"method"`
	OnSuccess string              `yaml:"onSuccess"`
	ResHeader []map[string]string `yaml:"resHeader"`
}

type Conf struct {
	Name          string `yaml:"name"`
	AssetBasePath string `yaml:"assetBasePath"`
	Routes        []Route
}

type Service struct {
}

func New(opt Options) *Service {
	var configs []Conf
	entries, err := os.ReadDir("./resources")
	if err != nil {
		log.Fatal(fmt.Errorf("os.ReadDir(): ./resources: %w", err))
	}

	for _, e := range entries {
		// read all index.yaml inside resources/*/index.yaml
		fname := fmt.Sprintf("resources/%s/index.yaml", e.Name())
		file, err := ioutil.ReadFile(fname)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				log.Fatalf("file doesn't exist: %q", fname)
				os.Exit(1)
			}
			log.Fatalf(err.Error())
		}

		var conf Conf
		err2 := yaml.Unmarshal(file, &conf)
		if err2 != nil {
			panic(err)
		}

		configs = append(configs, conf)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Default().Println("received req: ", r.RequestURI)
		for _, config := range configs {
			for _, route := range config.Routes {
				if route.Path == r.RequestURI { // match
					jsonPath := config.AssetBasePath + route.OnSuccess
					fContents, err := os.ReadFile(jsonPath)
					if err != nil {
						log.Println(err.Error())
						return
					}
					for _, h := range route.ResHeader {
						for k, v := range h {
							// w.Header().Add("content-type", "application/json")
							w.Header().Add(k, v)
						}
					}

					w.Write(fContents)
					return
				}
			}
		}
		log.Printf("route not registered: %s \n", r.RequestURI)
	})

	for _, v := range opt.ServiceMap {
		log.Println("starting http server at: ", v.Address)
		if err := http.ListenAndServe(v.Address, nil); err != nil {
			log.Fatalf("http.ListenAndServe(%v): %v", v, err)
		}
	}

	return &Service{}
}

func (s *Service) ShutDown() {
	log.Println("shutting down...")
}
