package main

import (
	"log"

	"github.com/icyrogue/ya-sher/internal/api"
	"github.com/icyrogue/ya-sher/internal/config"
	"github.com/icyrogue/ya-sher/internal/dbstorage"
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
	usermanager, err := usermanager.New()
	if err != nil {
		log.Fatal(err)
	}
	if opts.DBOpts.DBPath == "" {
	storage := urlstorage.New()
	storage.Options = opts.StrOpts
	usecase := idgen.New(storage)
<<<<<<< HEAD
	usermanager, err := usermanager.New()
	if err != nil {
		log.Fatal(err)
	}
	api := api.New(logger, opts.URLOpts, usecase, storage, usermanager)
=======
	api := api.New(logger, opts.URLOpts, usecase, storage, usermanager)
	storage.Init()
>>>>>>> inc10
	api.Init()
	api.Run()
	defer storage.Close()
	return
	} else {
	storage := dbstorage.New()
	storage.Options = opts.DBOpts
	usecase := idgen.New(storage)
	api := api.New(logger, opts.URLOpts, usecase, storage, usermanager)
	storage.Init()
	api.Init()
	api.Run()
	defer storage.Close()
	}
}
