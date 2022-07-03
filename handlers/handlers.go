package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"

	idgen "github.com/icyrogue/ya-sher/IdGen"
)

var urls = make(map[string]string)
var hostname string = "http://localhost:8080/"

// Create short version from long one
func CrShort() gin.HandlerFunc {
	return func(c *gin.Context) {

		defer c.Request.Body.Close()
		req, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			fmt.Println(err)
		}
		r := regexp.MustCompile(`.*\..*`)
		if !r.MatchString(string(req)) {
			c.String(http.StatusBadRequest, "This isn't an URL!")
			return
		}
		if el, fd := urls[string(req)]; fd {
			c.String(http.StatusCreated, el)
			return
		}

		urls[string(req)] = idgen.GenID(string(req))
		c.String(http.StatusCreated, hostname+urls[string(req)])
	}

}

//Get original from id
func ReLong() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.String(http.StatusBadRequest, "This isn't an id")
			return
		}
		for key, val := range urls {
			if val == id {
				c.Header("Location", key)
				c.String(http.StatusTemporaryRedirect, key)
			}
		}
	}
}
