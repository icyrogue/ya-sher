package config

import (
	"flag"
	"os"

	"github.com/caarlos0/env/v6"
)

type Cfg struct {
	URLOpts URLOpts
	StrOpts StrOpts
}
type URLOpts struct {
	Hostname string `env:"SERVER_ADDRESS" envDefault:"http://localhost:8080"`
	BaseURL  string `env:"BASE_URL" envDefault:"http://localhost:8080"`
}

type StrOpts struct {
	Filepath string `env:"FILE_STORAGE_PATH"`
}

//GetOpts gives defines options for everyone!
func GetOpts() (*Cfg, error) {
	opts := Cfg{}
	if err := env.Parse(&opts); err != nil {
		return nil, err
	}
	if len(os.Args) < 2 {
		return &opts, nil
	}

	flag.StringVar(&opts.URLOpts.Hostname, "a", "http://localhost:8080", "Hostname URL, default is localhost:8080")
	flag.StringVar(&opts.URLOpts.BaseURL, "b", "http://localhost:8080", "Hostname URL, default is localhost:8080")
	flag.StringVar(&opts.StrOpts.Filepath, "f", "", "Hostname URL, default is localhost:8080")
	flag.Parse()
	return &opts, nil

}
