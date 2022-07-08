package api

import (
	"io/ioutil"
	"net/http"
	"path"
	"regexp"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type URLProcessor interface {
	CreateShortURL(long string) (shurl string, err error)
}
type Storage interface {
	GetByLong(long string) *string
	GetByID(id string) *string
}

type api struct {
	router  *gin.Engine
	logger  *zap.Logger
	opts    *Options
	urlProc URLProcessor
	st      Storage
}

type Options struct {
	Hostname string
}

func New(logger *zap.Logger, opts *Options, urlProc URLProcessor, st Storage) *api {
	return &api{
		opts:    opts,
		logger:  logger,
		urlProc: urlProc,
		st:      st,
	}
}

func (a *api) Init() {
	gin.SetMode(gin.ReleaseMode)
	a.router = gin.New()
	a.router.POST("/", a.CrShort())
	a.router.GET("/:id", a.ReLong())
}
func (a *api) Run() {
	a.router.Run()

}

//CrShort: post short version from long one
func (a *api) CrShort() gin.HandlerFunc {
	return func(c *gin.Context) {

		defer c.Request.Body.Close()
		req, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			a.logger.Error("couldnt read file")
			return
		}
		r := regexp.MustCompile(`.*\..*`)
		if !r.MatchString(string(req)) {
			c.String(http.StatusBadRequest, "This isn't an URL!")
			return
		}
		if el := a.st.GetByLong(string(req)); el != nil {
			c.String(http.StatusCreated, *el)
			return
		}

		url, err := a.urlProc.CreateShortURL(string(req))
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		c.String(http.StatusCreated, path.Join(a.opts.Hostname, url))
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
		key := a.st.GetByID(id)
		if key == nil {
			c.String(http.StatusNotFound, "There isnt a url for this id")
			return
		}
		var long = *key
		c.Header("Location", long)
		c.String(http.StatusTemporaryRedirect, long)

	}
}
