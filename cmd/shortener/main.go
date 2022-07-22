package main

import (
	"log"

	"github.com/icyrogue/ya-sher/internal/api"
	"github.com/icyrogue/ya-sher/internal/config"
	"github.com/icyrogue/ya-sher/internal/idgen"
	"github.com/icyrogue/ya-sher/internal/urlstorage"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	opts, err := config.GetOpts()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(opts)
	storage := urlstorage.New()
	usecase := idgen.New(storage)
	api := api.New(logger, &api.Options{Hostname: opts.URLOpts.Hostname, BaseURL: opts.URLOpts.BaseURL}, usecase, storage)
	api.Init()
	api.Run()
}
