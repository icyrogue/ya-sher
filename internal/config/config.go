package config

import (
	"flag"
	"os"
)

type Cfg struct {
	URLOpts URLOpts
	StrOpts StrOpts
}
type URLOpts struct {
	Hostname string `env:"SERVER_ADDRESS,unset" envDefault:"http://localhost:8080"`
	BaseURL  string `env:"BASE_URL,unset" envDefault:"http://localhost:8080"`
}

type StrOpts struct {
	Filepath string `env:"FILE_STORAGE_PATH"`
}

//GetOpts gives defines options for everyone!
func GetOpts() (*Cfg, error) {
	cfg := Cfg{}
	flag.StringVar(&cfg.URLOpts.Hostname, "a", "http://localhost:8080", "Hostname URL")
	flag.StringVar(&cfg.URLOpts.BaseURL, "b", "http://localhost:8080", "Base URL")
	flag.StringVar(&cfg.StrOpts.Filepath, "f", "", "File path")
	flag.Lookup("f").Value.Set(os.Getenv("FILE_STORAGE_PATH"))
	flag.Lookup("a").Value.Set(os.Getenv("SERVER_ADDRESS"))
	flag.Lookup("b").Value.Set(os.Getenv("SERVER_ADDRESS"))
	flag.Parse()
	return &cfg, nil

}
