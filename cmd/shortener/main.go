package main

import (
	"log"

	"github.com/icyrogue/ya-sher/internal/api"
	"github.com/icyrogue/ya-sher/internal/idgen"
	"github.com/icyrogue/ya-sher/internal/urlstorage"
	"go.uber.org/zap"
)

func main() {
	idgen.InitID()
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	storage := urlstorage.New()
	usecase := idgen.New(storage)
	api := api.New(logger, &api.Options{Hostname: "http://localhost:8080"}, usecase, storage)
	api.Init()
	api.Run()
}
