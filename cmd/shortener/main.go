package main

import (
	idgen "github.com/icyrogue/ya-sher/IdGen"
	"github.com/icyrogue/ya-sher/handlers"

	"github.com/gin-gonic/gin"
)

func apiInit() *gin.Engine {

	r := gin.New()
	r.POST("/", handlers.CrShort())
	r.GET("/:id", handlers.ReLong())
	return r
}
func main() {
	idgen.InitID()
	gin.SetMode(gin.ReleaseMode)
	r := apiInit()
	r.Run()

}
