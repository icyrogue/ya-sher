package main

import (
	"github.com/icyrogue/ya-sher/internal/api"
	"github.com/icyrogue/ya-sher/internal/idgen"

	"github.com/gin-gonic/gin"
)

func main() {
	idgen.InitID()
	gin.SetMode(gin.ReleaseMode)
	r := api.ApiInit()
	r.Run()

}
