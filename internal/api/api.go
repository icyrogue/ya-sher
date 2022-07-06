package api

import (
	"github.com/gin-gonic/gin"
	"github.com/icyrogue/ya-sher/internal/urlstorage"
	"go.uber.org/zap"
	"io/ioutil"

	"net/http"
	"regexp"
)

type Options struct {
	Hostname string
}

type URLProcessor interface {
	CreateShortURL(long string) (shurl string, err error)
}

type api struct {
	gengine *gin.Engine
	logger  *zap.Logger
	opts    *Options

	urlProc URLProcessor
}

func New(logger *zap.Logger, opts *Options, urlProc URLProcessor) *api {
	return &api{
		opts:    opts,
		logger:  logger,
		urlProc: urlProc,
	}
}

func (a *api) Init() {
	a.gengine = gin.New()
	a.gengine.POST("/", a.CrShort())
	a.gengine.GET("/:id", a.ReLong())

}

func (a *api) Run() {
	a.gengine.Run()
}

//CrShort: post short version from long one
func (a *api) CrShort() gin.HandlerFunc {
	return func(c *gin.Context) {

		defer c.Request.Body.Close()
		req, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			a.logger.Error("unable err")
			return
		}
		r := regexp.MustCompile(`.*\..*`)
		if !r.MatchString(string(req)) {
			c.String(http.StatusBadRequest, "This isn't an URL!")
			return
		}
		//

		shortUrl, err := a.urlProc.CreateShortURL(string(req))
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		c.String(http.StatusCreated, shortUrl)

		//if el := urlstorage.GetByLong(string(req)); el != nil {
		//	c.String(http.StatusCreated, el.Short)
		//	return
		//}
		//
		//url := urlstorage.NewUrl(string(req), a.opts.Hostname+idgen.GenID())
		//c.String(http.StatusCreated, url.Short)

	}

}

//Relong: get original from id
func (a *api) ReLong() gin.HandlerFunc {
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
