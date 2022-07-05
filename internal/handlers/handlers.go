package handlers

import (
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"

	"github.com/icyrogue/ya-sher/internal/idgen"
	"github.com/icyrogue/ya-sher/internal/urlstorage"
)

var hostname string = "http://localhost:8080/"

//CrShort: post short version from long one
func CrShort() gin.HandlerFunc {
	return func(c *gin.Context) {

		defer c.Request.Body.Close()
		req, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Fatal(err)
		}
		r := regexp.MustCompile(`.*\..*`)
		if !r.MatchString(string(req)) {
			c.String(http.StatusBadRequest, "This isn't an URL!")
			return
		}
		if el := urlstorage.GetByLong(string(req)); el != nil {
			c.String(http.StatusCreated, el.Short)
			return
		}

		url := urlstorage.NewUrl(string(req), hostname+idgen.GenID())
		c.String(http.StatusCreated, url.Short)
		url.AddToStorage()
	}

}

//Relong: get original from id
func ReLong() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.String(http.StatusBadRequest, "This isn't an id")
			return
		}
		key := urlstorage.GetByID(id)
		c.Header("Location", key.Long)
		c.String(http.StatusTemporaryRedirect, key.Long)

	}
}
