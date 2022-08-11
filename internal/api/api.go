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
	"github.com/icyrogue/ya-sher/internal/jsonmodels"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

//JSON model aliases
type jsonURL = jsonmodels.JSONURL

type jsonResult = jsonmodels.JSONResult

type jsonURLTouple = jsonmodels.JSONURLTouple

//URLProcessor interface for creating short url using idgen business logic
type URLProcessor interface {
	CreateShortURL(long string) (shurl string, err error)
}

//Storage interface for interfacing with storage
type Storage interface {
	GetByLong(long string, ctx context.Context) (string, error)
	GetByID(id string, ctx context.Context) (string, error)
	Ping(ctx context.Context) bool
}

type UserManager interface {
	AddUserURL(user string, long string, id string) error
	NewUser() (string, error)
	CheckValid(cookie string) bool
	GetAllUserURLs(cookie string) map[string]string
}

type UserManager interface {
	AddUserURL(user string, long string, id string) error
	NewUser() (string, error)
	CheckValid(cookie string) bool
	GetAllUserURLs(cookie string) map[string]string
}

type api struct {
	router  *gin.Engine
	logger  *zap.Logger
	opts    *Options
	urlProc URLProcessor
	st      Storage
	userManager UserManager
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

func New(logger *zap.Logger, opts *Options, urlProc URLProcessor, st Storage, userManager UserManager) *api {
	df := `http://localhost:8080`
	if opts == nil {
		opts = &Options{
			BaseURL:  df,
			Hostname: df,
		}
	}
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
		userManager: userManager,
	}
}

func (a *api) Init() {
	gin.SetMode(gin.DebugMode)
	a.router = gin.Default()
	a.router.Use(a.mdwDecompression, a.mdwCompression, a.mdwCookie)
	a.router.POST("/", a.CrShort)
	a.router.GET("/:id", a.ReLong)
	a.router.POST("/api/shorten", a.Shorten)
	a.router.GET("/api/user/urls", a.getAllUserURLs)
	a.router.GET("/ping", a.pingDB)
}
func (a *api) Run() {
	re := regexp.MustCompile(`:\d*$`)
	a.router.Run(re.FindString(a.opts.Hostname))

}

//CrShort: post short version from long one
func (a *api) CrShort(c *gin.Context) {
	cookie := c.MustGet("cookie")
	if cookie == "" {
	c.String(http.StatusBadRequest, "couldnt identify user")
	return
	}
	fmt.Println(cookie)
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
	if el, err := a.st.GetByLong(string(req), c); err == nil {
		a.userManager.AddUserURL(fmt.Sprint(cookie), string(req), el)
		fmt.Println(a.st.GetByLong(string(req), c))
		c.String(http.StatusCreated, a.opts.BaseURL+"/"+el)
		return
	}

		url, err := a.urlProc.CreateShortURL(string(req))
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	err  = a.userManager.AddUserURL(fmt.Sprint(cookie), string(req), url)
	if err != nil {
		println(err.Error())
	}
	c.String(http.StatusCreated, a.opts.BaseURL + "/" + url)


}

//Relong: get original from id
func (a *api) ReLong(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.String(http.StatusBadRequest, "This isn't an id")
		return
	}
	key, err := a.st.GetByID(id, c)
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

	cookie := fmt.Sprint(c.MustGet("cookie"))

	url := jsonURL{}
	res := c.Request.Body

	defer res.Close()
	body, err := ioutil.ReadAll(res)
	if err != nil {
		a.logger.Error("couldnt read request")
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	err = json.Unmarshal(body, &url)
	if err != nil {
		a.logger.Error(err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.Header("Content-Type", "application/json")
	shurl, err := a.urlProc.CreateShortURL(url.URL)

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	var result []byte
	resURL := jsonResult{
		Result: a.opts.BaseURL + "/" + shurl,
	}
	if result, err = json.Marshal(resURL); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	err = a.userManager.AddUserURL(cookie, url.URL, shurl)
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

func (a *api) mdwCookie(c *gin.Context) {
	cookie, err := c.Request.Cookie("url_shortner")

	if err != nil && err.Error() != "http: named cookie not present" {
c.String(http.StatusBadRequest, err.Error())
		return
	}
	if cookie != nil {
		c.Set("cookie", cookie.Value)
		fmt.Println(cookie.Value)

		c.Next()
		return
	}
	newCookie, err := a.userManager.NewUser()
	if err != nil {
		c.String(http.StatusBadRequest, err.Error() )
		return
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name: "url_shortner",
		Value: newCookie,
		MaxAge: 999,
	})
	c.Set("cookie", newCookie)
	c.Next()
	}

func (a *api) getAllUserURLs(c *gin.Context) {
	var err error

	cookie := c.MustGet("cookie")

	touples := []jsonURLTouple{}

	c.Header("Content-Type", "application/json")

	urls := a.userManager.GetAllUserURLs(fmt.Sprint(cookie))
	if len(urls) == 0 {
		c.String(http.StatusNoContent, "")
		return
	}
	for short, long := range(urls) {
		touples = append(touples, jsonURLTouple{
			Short: a.opts.BaseURL + `/` + short,
			Long: long,
		})
		fmt.Println(short, long)
	}
	if res, err := json.Marshal(touples); err == nil {
		c.String(http.StatusOK, string(res))
		return
	}
	c.String(http.StatusNoContent, err.Error())
}

func (a *api) pingDB(c *gin.Context) {
	if a.st.Ping(c) {
		c.String(http.StatusOK, "" )
		return
	}
	c.String(http.StatusInternalServerError, "" )
}
