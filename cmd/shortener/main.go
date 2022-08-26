package main

import (
	"context"
	"log"

	"github.com/icyrogue/ya-sher/internal/api"
	"github.com/icyrogue/ya-sher/internal/config"
	"github.com/icyrogue/ya-sher/internal/idgen"
	"github.com/icyrogue/ya-sher/internal/mlt"
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

	msh := musher.New(opts.MushOpts, storage)
	mlt := mlt.New(usermanager)
	mlt.Output = msh.Input
	ctx := context.Background()

	api := api.New(logger, opts.URLOpts, usecase, storage, usermanager, mlt)
	storage.Init()
	api.Init()

	msh.Start(ctx)
	mlt.Start(ctx)

	api.Run()

	defer storage.Close()
}
