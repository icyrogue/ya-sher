package main

import (
	"github.com/icyrogue/ya-sher/internal/api"
	"github.com/icyrogue/ya-sher/internal/idgen"
	"github.com/icyrogue/ya-sher/internal/urlstorage"
	"go.uber.org/zap"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	idgen.InitID()
	gin.SetMode(gin.ReleaseMode)

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalln(err)
	}

	storage := urlstorage.New()

	usecase := idgen.New(storage)

	api := api.New(logger, &api.Options{Hostname: "http://localhost:8080/"}, usecase)

	api.Init()

	api.Run()

}
