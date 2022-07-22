package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

type Cfg struct {
	URLOpts *URLOpts
}
type URLOpts struct {
	Hostname string `env:"SERVER_ADDRESS" envDefault:"http://localhost:8080"`
	BaseURL  string `env:"BASE_URL" envDefault:"http://localhost:8080"`
}

//GetOpts gives defines options for everyone!
func GetOpts() (*Cfg, error) {
	urlOpts := URLOpts{}
	if err := env.Parse(&urlOpts); err != nil {
		return nil, err
	}
	fmt.Println(urlOpts)
	return &Cfg{
		URLOpts: &urlOpts,
	}, nil

}
