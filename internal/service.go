package internal

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"golang.org/x/exp/slog"
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

func New(opt Options) (*Service, error) {
	entries, err := os.ReadDir("./resources")
	if err != nil {
		slog.Error("os.ReadDir(): ./resources", "error", err)
	}
	s := &Service{}
	for _, e := range entries {
		// read all index.yaml inside resources/*/index.yaml
		fname := fmt.Sprintf("resources/%s/index.yaml", e.Name())
		file, err := ioutil.ReadFile(fname)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return nil, fmt.Errorf("%q doesn't exist: %w", fname, err)
			}
			return nil, fmt.Errorf("internal error: %w", err)
		}

		var conf Conf
		if err := yaml.Unmarshal(file, &conf); err != nil {
			return nil, fmt.Errorf("unmarsal failed: %w", err)
		}
		s.confs = append(s.confs, conf)
	}

	return s, nil
}

func (s *Service) handle() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("receviced req", "uri", r.RequestURI)
		for _, config := range s.confs {
			for _, route := range config.Routes {
				if route.Path == r.RequestURI { // match
					jsonPath := config.AssetBasePath + route.OnSuccess
					fContents, err := os.ReadFile(jsonPath)
					if err != nil {
						slog.Error(err.Error())
						return
					}
					for _, h := range route.ResHeader {
						for k, v := range h {
							w.Header().Add(k, v)
						}
					}

					w.Write(fContents)
					return
				}
			}
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
		slog.Error("route not registered", "req uri", r.RequestURI)
	})
}

func (s *Service) Start() {
	for _, c := range s.confs {
		go func(c Conf) {
			slog.Info(fmt.Sprintf("listening and serveing %s at %s", c.Name, c.ListenAddrHTTP))
			http.ListenAndServe(c.ListenAddrHTTP, nil)
		}(c)
	}

	s.handle()
	select {}
}

func (s *Service) ShutDown() {
	slog.Info("//")
}
