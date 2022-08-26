package api

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

type jsonURLTuple = jsonmodels.JSONURLTuple

type jsonBlkIn = jsonmodels.JSONBulkInput

//URLProcessor interface for creating short url using idgen business logic
type URLProcessor interface {
	CreateShortURL(long string) (shurl string, err error)
	BulkCreation(data []jsonBlkIn, baseURL string) ([]jsonBlkIn, error)
}

//Storage interface for interfacing with storage
type Storage interface {
	GetByLong(long string, ctx context.Context) (string, error)
	GetByID(ctx context.Context, id string) (string, error)
	Ping(ctx context.Context) error
}

type Mlt interface {
	GetInput() chan []string
}

//User manager methods for managing users based on cookie
type UserManager interface {
	AddUserURL(user string, long string, id string) error
	NewUser() (string, error)
	GetAllUserURLs(cookie string) map[string]string

}

type api struct {
	router  *gin.Engine
	logger  *zap.Logger
	opts    *Options
	urlProc URLProcessor
	st      Storage
	userManager UserManager
	mlt Mlt
	//wg *sync.WaitGroup
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

func New(logger *zap.Logger, opts *Options, urlProc URLProcessor, st Storage, userManager UserManager, mlt Mlt) *api {
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
		mlt: mlt,
		//wg: &sync.WaitGroup{},
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
	a.router.POST("/api/shorten/batch", a.convertBulk)
	a.router.DELETE("/api/user/urls", a.Delete)
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
		c.String(http.StatusConflict, a.opts.BaseURL+"/"+el)
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
	key, err := a.st.GetByID(c, id)
	if err != nil {
		c.String(http.StatusGone, err.Error())
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

	c.Header("Content-Type", "application/json")

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
	if el, err := a.st.GetByLong(url.URL, c); err == nil {
	resURL := jsonResult{
		Result: a.opts.BaseURL + "/" + el,
	}
		result, err := json.Marshal(&resURL)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		c.String(http.StatusConflict, string(result))
		return
	}

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
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.String(http.StatusCreated, string(result))
}

//mdwCompression: gzip compression middleware
func (a *api) mdwCompression(c *gin.Context) {
	if !strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip") {

		c.Next()
		return
	}
	gz, err := gzip.NewWriterLevel(c.Writer, gzip.BestCompression)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
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

	c.Writer.Header().Del("Content-Length") //<-- otherwise corruption occurs
	c.Request.Body = ioutil.NopCloser(gz)
	c.Next()
}

//middleware to sort out cookie related sruff
func (a *api) mdwCookie(c *gin.Context) {
	cookie, err := c.Request.Cookie("url_shortner")

	if err != nil && err.Error() != "http: named cookie not present" {
c.String(http.StatusBadRequest, err.Error())
		return
	}
	if cookie != nil {
		c.Set("cookie", cookie.Value)

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

//getAllUserURLs: GET all urls that user with this cookie shortened
func (a *api) getAllUserURLs(c *gin.Context) {
	//	var err error

	cookie := c.MustGet("cookie")

	tuples := []jsonURLTuple{}

	c.Header("Content-Type", "application/json")

	urls := a.userManager.GetAllUserURLs(fmt.Sprint(cookie))
	if len(urls) == 0 {
		c.String(http.StatusNoContent, "")
		return
	}
	for short, long := range(urls) {
		tuples = append(tuples, jsonURLTuple{
			Short: a.opts.BaseURL + `/` + short,
			Long: long,
		})
	}
	if res, err := json.Marshal(tuples); err == nil {
		c.String(http.StatusOK, string(res))
		return
	} else {
		c.String(http.StatusBadRequest, err.Error())
	}
/*Здесь вот я не понимаю, я проверяю, если ошибка == nil, если так, то return, если != nil, то проходит
  дальше и передает ее в респонс, но тогда go vet тест ругается, что я nil дереференсю, было как снизу else не было*/

//c.String(http.StatusBadRequest, err.Error())

}

//pingDB: GET database response
func (a *api) pingDB(c *gin.Context) {
	if err := a.st.Ping(c); err != nil {
		c.String(http.StatusInternalServerError, err.Error() )
		return
	}
	c.String(http.StatusOK, "" )
}

//convertBulk: POST multiple urls to shorten at once
func (a *api) convertBulk(c *gin.Context) {
	defer c.Request.Body.Close()

	body, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		a.logger.Error("couldnt read request body")
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	input := []jsonBlkIn{}
	var output []byte

	if err = json.Unmarshal(body, &input); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	if input, err = a.urlProc.BulkCreation(input, a.opts.BaseURL); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if output, err = json.Marshal(&input); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.Header("Content-Type", "application/json")
	c.String(http.StatusCreated, string(output))
}

func (a *api) Delete(c *gin.Context) {
	cookie, err := c.Request.Cookie("url_shortner")
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	defer c.Request.Body.Close()
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	data := []string{}
	if err = json.Unmarshal(body, &data); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.String(http.StatusAccepted, "")

	go func (){
		data = append(data, cookie.Value)
		a.mlt.GetInput() <- data
		log.Println(data)

	}()
}
