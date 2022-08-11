package main

import (
	"log"

	"github.com/icyrogue/ya-sher/internal/api"
	"github.com/icyrogue/ya-sher/internal/config"
	"github.com/icyrogue/ya-sher/internal/idgen"
	"github.com/icyrogue/ya-sher/internal/urlstorage"
	"github.com/icyrogue/ya-sher/internal/usermanager"
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
	storage := urlstorage.New()
	storage.Options = opts.StrOpts
	storage.Init()
	usecase := idgen.New(storage)
	usermanager, err := usermanager.New()
	if err != nil {
		log.Fatal(err)
	}
	api := api.New(logger, opts.URLOpts, usecase, storage, usermanager)
	api.Init()
	api.Run()
	defer storage.Close()
}
