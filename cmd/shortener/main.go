package main

import (
	"context"
	"log"

	"github.com/icyrogue/ya-sher/internal/api"
	"github.com/icyrogue/ya-sher/internal/config"
	"github.com/icyrogue/ya-sher/internal/idgen"
	"github.com/icyrogue/ya-sher/internal/multithreaddeleteurlprocessor"
	"github.com/icyrogue/ya-sher/internal/musher"
	"github.com/icyrogue/ya-sher/internal/storager"
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
	storage := storager.Get(*opts)
	usecase := idgen.New(storage)

	musher := musher.New(opts.MushOpts, storage)
	multiThreadDeleteURLProcessor := mlt.New(usermanager)
	multiThreadDeleteURLProcessor.Output = musher.Input
	ctx := context.Background()

	api := api.New(logger, opts.URLOpts, usecase, storage, usermanager, multiThreadDeleteURLProcessor)
	storage.Init()
	api.Init()

	musher.Start(ctx)
	multiThreadDeleteURLProcessor.Start(ctx)

	api.Run()

	defer storage.Close()
}
