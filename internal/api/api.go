package api

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//JSON models
type jsonURL struct {
	URL string `json:"url"`
}

type jsonResult struct {
	Result string `json:"result"`
}

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
	BaseURL  string
}

//Struct for new gzip writer
type gzipWriter struct {
	gin.ResponseWriter
	writer *gzip.Writer
}

//Write string method needed for gin.conetxt.String
func (g *gzipWriter) WriteString(s string) (int, error) {
	g.Header().Del("Content-Length") //<-- deleting Content-Length header since
	return g.writer.Write([]byte(s)) //new length dosent match after compression
}

func (g *gzipWriter) Write(data []byte) (int, error) {
	g.Header().Del("Content-Length")
	return g.writer.Write(data)
}

func New(logger *zap.Logger, opts *Options, urlProc URLProcessor, st Storage) *api {
	df := `http://localhost:8080`
	if opts.BaseURL == "" {
		opts.BaseURL = df
	}
	if opts.Hostname == "" {
		opts.Hostname = df
	}
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
	a.router.Use(a.mdwDecompression, a.mdwCompression)
	a.router.POST("/", a.CrShort)
	a.router.GET("/:id", a.ReLong)
	a.router.POST("/api/shorten", a.Shorten)
}
func (a *api) Run() {
	re := regexp.MustCompile(`:\d*$`)
	a.router.Run(re.FindString(a.opts.Hostname))

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
		c.String(http.StatusCreated, a.opts.BaseURL+"/"+el)
		return
	}

	url, err := a.urlProc.CreateShortURL(string(req))
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusCreated, a.opts.BaseURL+"/"+url) //<-┐
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
	fmt.Println("key is ", key)
	c.String(http.StatusTemporaryRedirect, key)

}

//Shorten: gives back json short link
func (a *api) Shorten(c *gin.Context) {

	url := jsonURL{}
	res := c.Request.Body

	defer res.Close()
	body, err := ioutil.ReadAll(res)
	if err != nil {
		a.logger.Error("couldnt read request")
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	json.Unmarshal(body, &url)

	c.Header("Content-Type", "application/json")
	shurl, err := a.urlProc.CreateShortURL(url.URL)

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	var result []byte
	var err3 error
	resURL := jsonResult{
		Result: a.opts.BaseURL + "/" + shurl,
	}
	if result, err3 = json.Marshal(resURL); err3 != nil {
		c.String(http.StatusInternalServerError, err3.Error())
		return
	}
	c.String(http.StatusCreated, string(result))
}

//mdwCompression: gzip compression middleware
func (a *api) mdwCompression(c *gin.Context) {
	if !strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip") {
		fmt.Println("normal mode")

		c.Next()
		return
	}
	gz, err := gzip.NewWriterLevel(c.Writer, gzip.BestCompression)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		//	return
	}

	defer gz.Close()
	c.Writer = &gzipWriter{c.Writer, gz}

	c.Header("Content-Encoding", "gzip")

	c.Next()
	//	return
}

//mdwDecompression: gzip decompression middleware
func (a *api) mdwDecompression(c *gin.Context) {
	if !strings.Contains(c.Request.Header.Get("Content-Encoding"), "gzip") {
		c.Next()
		return
	}
	gz, err := gzip.NewReader(c.Request.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	defer gz.Close()
	/* newBody, err := io.ReadAll(gz)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	} */
	c.Writer.Header().Del("Content-Length") //<-- otherwise corruption occurs
	c.Request.Body = ioutil.NopCloser(gz)
	c.Next()
}
