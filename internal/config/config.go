package config

import (
	"flag"
	"os"

	"github.com/icyrogue/ya-sher/internal/api"
	"github.com/icyrogue/ya-sher/internal/urlstorage"
)

type Cfg struct {
	URLOpts *api.Options
	StrOpts *urlstorage.Options
}

//GetOpts: defines options for everyone!
func GetOpts() (*Cfg, error) {
	cfg := Cfg{
		StrOpts: &urlstorage.Options{},
		URLOpts: &api.Options{},
	}
	flag.StringVar(&cfg.URLOpts.Hostname, "a", "http://localhost:8080", "Hostname URL")
	flag.StringVar(&cfg.URLOpts.BaseURL, "b", "http://localhost:8080", "Base URL")
	flag.StringVar(&cfg.StrOpts.Filepath, "f", "", "File path")
	flag.Lookup("f").Value.Set(os.Getenv("FILE_STORAGE_PATH"))
	flag.Lookup("a").Value.Set(os.Getenv("SERVER_ADDRESS"))
	flag.Lookup("b").Value.Set(os.Getenv("BASE_URL"))
	flag.Parse()
	return &cfg, nil
	//Я перестал использовать env.Parse() -> error, но возвращение ошибки оста-
	//вил на будущее, если что то нужно будет сделать обязательным
}
