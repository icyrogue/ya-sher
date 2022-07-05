package api

import (
	"github.com/gin-gonic/gin"
	"github.com/icyrogue/ya-sher/internal/handlers"
)

func ApiInit() *gin.Engine {

	r := gin.New()
	r.POST("/", handlers.CrShort())
	r.GET("/:id", handlers.ReLong())
	return r
}
