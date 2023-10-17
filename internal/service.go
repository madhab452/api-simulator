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

type Options struct {
}

type Route struct {
	Path      string              `yaml:"path"`
	Method    string              `yaml:"method"`
	OnSuccess string              `yaml:"onSuccess"`
	ResHeader []map[string]string `yaml:"resHeader"`
}

type Conf struct {
	Name           string `yaml:"name"`
	ListenAddrHTTP string `yaml:"listenAddrHTTP"`
	AssetBasePath  string `yaml:"assetBasePath"`
	Routes         []Route
}

type Service struct {
	confs []Conf
}

func New(opt Options) *Service {
	entries, err := os.ReadDir("./resources")
	if err != nil {
		log.Fatal(fmt.Errorf("os.ReadDir(): ./resources: %w", err))
	}
	serv := &Service{}
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
		if err := yaml.Unmarshal(file, &conf); err != nil {
			panic(err)
		}
		serv.confs = append(serv.confs, conf)
	}

	return serv
}

func (s *Service) handle() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Default().Println("received req: ", r.RequestURI)
		for _, config := range s.confs {
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
}

func (s *Service) Start() {
	for _, c := range s.confs {
		go func(c Conf) {
			log.Default().Printf("listening and serving %s at %s \n", c.Name, c.ListenAddrHTTP)
			http.ListenAndServe(c.ListenAddrHTTP, nil)
		}(c)
	}

	s.handle()
	select {}
}

func (s *Service) ShutDown() {
	log.Println("shutting down...")
}
