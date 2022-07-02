package main

import (
	"github.com/icyrogue/ya-sher/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/", handlers.CrShort())
	r.GET("/:id", handlers.ReLong())
	r.Run()
}
