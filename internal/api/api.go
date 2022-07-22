package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//URLProcessor interface for creating short url using idgen business logic
type URLProcessor interface {
	CreateShortURL(long string) (shurl string, err error)
}

//Storage interface for interfacing with storage
type Storage interface {
	GetByLong(long string) (string, error)
	GetByID(id string) (string, error)
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

	gin.SetMode(gin.DebugMode)
	a.router = gin.Default()
	a.router.POST("/", a.CrShort)
	a.router.GET("/:id", a.ReLong)
	a.router.POST("/api/shorten", a.Shorten)
}
func (a *api) Run() {
	a.router.Run()

}

//CrShort: post short version from long one
func (a *api) CrShort(c *gin.Context) {

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
	if el, errEl := a.st.GetByLong(string(req)); errEl == nil {
		c.String(http.StatusCreated, a.opts.Hostname+"/"+el)
		return
	}

	url, err := a.urlProc.CreateShortURL(string(req))
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusCreated, a.opts.Hostname+"/"+url) //<-┐
	//Если использовать  Path.Join, то автотест ставит ///  --┘
}

//Relong: get original from id
func (a *api) ReLong(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.String(http.StatusBadRequest, "This isn't an id")
		return
	}
	key, err := a.st.GetByID(id)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}
	c.Header("Location", key)
	c.String(http.StatusTemporaryRedirect, key)

}

//Shorten: gives back json short link
func (a *api) Shorten(c *gin.Context) {
	type tmp struct {
		Url string `json:"url"`
	}

	url := tmp{}
	res := c.Request.Body

	defer res.Close()
	body, err := ioutil.ReadAll(res)
	if err != nil {
		a.logger.Error("couldnt read request")
		return
	}
	json.Unmarshal(body, &url)

	c.Header("Content-Type", "application/json")
	shurl, err2 := a.urlProc.CreateShortURL(url.Url)

	if err2 != nil {
		c.String(http.StatusInternalServerError, err2.Error())
		return
	}
	var result []byte
	var err3 error
	resUrl := struct {
		Result string `json:"result"`
	}{
		Result: a.opts.Hostname + "/" + shurl,
	}
	if result, err3 = json.Marshal(resUrl); err3 != nil {
		c.String(http.StatusInternalServerError, err3.Error())
		return
	}
	c.String(http.StatusCreated, string(result))
}
