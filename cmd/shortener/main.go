package main

import (
	"net/http"

	"github.com/icyrogue/ya-sher/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "link shortner")
	})

	r.GET("/popa/", handlers.Popa())
	r.POST("/", handlers.CrShort())
	r.GET("/:id", handlers.ReLong())
	r.Run()
}
