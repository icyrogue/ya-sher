package config

import (
	"flag"
	"log"
	"os"

	"github.com/icyrogue/ya-sher/internal/api"
	"github.com/icyrogue/ya-sher/internal/dbstorage"
	"github.com/icyrogue/ya-sher/internal/musher"
	"github.com/icyrogue/ya-sher/internal/urlstorage"
)

type Cfg struct {
	URLOpts *api.Options
	StrOpts *urlstorage.Options
	DBOpts *dbstorage.Options
	MushOpts *musher.Options

}

//GetOpts: defines options for everyone!
func GetOpts() (*Cfg, error) {
	cfg := Cfg{
		StrOpts: &urlstorage.Options{},
		URLOpts: &api.Options{},
		DBOpts: &dbstorage.Options{},
		MushOpts: &musher.Options{},
	}
	flag.StringVar(&cfg.URLOpts.Hostname, "a", "http://localhost:8080", "Hostname URL")
	flag.StringVar(&cfg.URLOpts.BaseURL, "b", "http://localhost:8080", "Base URL")
	flag.StringVar(&cfg.StrOpts.Filepath, "f", "", "File path")
	flag.StringVar(&cfg.DBOpts.DBPath, "d", "", "db path")
	flag.IntVar(&cfg.MushOpts.MaxWaitTime, "mt", 10, "Max wait time before commiting to db")
	flag.IntVar(&cfg.MushOpts.MaxBufferLength, "mb", 60, "Max db buffer length")
	flag.Lookup("f").Value.Set(os.Getenv("FILE_STORAGE_PATH"))
	flag.Lookup("a").Value.Set(os.Getenv("SERVER_ADDRESS"))
	flag.Lookup("b").Value.Set(os.Getenv("BASE_URL"))
	flag.Lookup("d").Value.Set(os.Getenv("DATABASE_DSN"))
	flag.Lookup("mt").Value.Set(os.Getenv("MAX_WAIT_TIME_SHORTNER"))
	flag.Lookup("mb").Value.Set(os.Getenv("MAX_BUFFER_LENGTH_SHORTNER"))
	flag.Parse()
	log.Println(cfg.DBOpts)
	return &cfg, nil
	//Я перестал использовать env.Parse() -> error, но возвращение ошибки оста-
	//вил на будущее, если что то нужно будет сделать обязательным
}
