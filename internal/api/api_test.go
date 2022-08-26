package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/icyrogue/ya-sher/internal/idgen"
	"github.com/icyrogue/ya-sher/internal/mlt"
	"github.com/icyrogue/ya-sher/internal/urlstorage"
	"github.com/icyrogue/ya-sher/internal/usermanager"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func Test_api_CrShort(t *testing.T) {

	tests := []struct {
		name       string
		want       string
		wantedCode int
	}{
		{name: "Simple POST test #1",
			want:       "smokybananas.com",
			wantedCode: http.StatusCreated},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//BOILERPLATE\\
			logger, err := zap.NewDevelopment()
			if err != nil {
				t.Fatal(err)
			}
			storage := urlstorage.New()
			storage.Init()
			usecase := idgen.New(storage)
			usermanager, err := usermanager.New()
			if err != nil {
				t.Fatal(err)
			}
			mlt := mlt.New(usermanager)
			api := New(logger, &Options{}, usecase, storage, usermanager, mlt)
			api.Init()
			//Testing POST itself
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(tt.want)))
			api.router.ServeHTTP(w, req)

			res := w.Result()

			if res.StatusCode != tt.wantedCode {
				t.Errorf("Expected %d got %d", tt.wantedCode, res.StatusCode)
			}
			defer res.Body.Close()
			body, err := ioutil.ReadAll(res.Body)
			assert.NoError(t, err)
			shurl, err := storage.GetByLong(tt.want, req.Context())
			assert.NoError(t, err)
			fmt.Println(string(body))
			cody := strings.SplitAfter(string(body), "/")[3]

			if cody[:8] != shurl {
				t.Errorf("Expected %v got %v", []byte(shurl), []byte(cody)[:8])
			}
		})
	}
}

func Test_api_ReLong(t *testing.T) {
	tests := []struct {
		name       string
		want       string
		wantedCode int
	}{{name: "Simple Test #1",
		want:       "google.com",
		wantedCode: http.StatusTemporaryRedirect},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//BOILERPLATE\\
			logger, err := zap.NewDevelopment()
			if err != nil {
				t.Fatal(err)
			}
			//	opts, _ := config.GetOpts()
			storage := urlstorage.New()
			storage.Init()
			usecase := idgen.New(storage)
			usermanager, err := usermanager.New()
			if err != nil {
				t.Fatal(err)
			}
			api := New(logger, &Options{}, usecase, storage, usermanager, mlt.New(usermanager))
			api.Init()
			//Creating mock short
			shurl, err1 := usecase.CreateShortURL(tt.want)
			if err1 != nil {
				t.Error(err1)

			}
			storage.Add(shurl, tt.want)
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/"+shurl, nil)
			api.router.ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.wantedCode {
				t.Errorf("Expected %d got %d", tt.wantedCode, res.StatusCode)
			}

			header := res.Header.Get("Location")
			if header != tt.want {
				t.Errorf("Expected %s got %s", tt.want, header)
			}

		})
	}
}

func Test_api_Shorten(t *testing.T) {

	tests := []struct {
		name       string
		want       string
		wantedCode int
	}{
		{name: "Simple POST with shorten #1",
			want:       "sosmosolonka.ru",
			wantedCode: http.StatusCreated},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//BOILERPLATE\\
			logger, err := zap.NewDevelopment()
			if err != nil {
				t.Fatal(err)
			}
			//	opts, _ := config.GetOpts()
			storage := urlstorage.New()
			storage.Init()
			usecase := idgen.New(storage)
			usermanager, err := usermanager.New()
			if err != nil {
				t.Fatal(err)
			}
			api := New(logger, &Options{}, usecase, storage, usermanager, mlt.New(usermanager))
			api.Init()
			//Testing POST itself
			w := httptest.NewRecorder()
			wantJSON := jsonURL{
				URL: tt.want,
			}

			reqJSON, errJ := json.Marshal(wantJSON)

			if errJ != nil {
				t.Fatal(errJ)
			}
			req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer(reqJSON))
			api.router.ServeHTTP(w, req)

			res := w.Result()

			if res.StatusCode != tt.wantedCode {
				t.Errorf("Expected %d got %d", tt.wantedCode, res.StatusCode)
			}
			defer res.Body.Close()
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Error(err)
			}
			bodyJSON := jsonResult{}

			if err := json.Unmarshal(body, &bodyJSON); err != nil {
				t.Fatal(err)
			}
			shurl, err := storage.GetByLong(tt.want, req.Context())
			if err != nil {
				t.Error(err.Error())
				return
			}
			body = []byte(bodyJSON.Result)
			body = body[len(body)-8:]
			if string(body) != shurl {
				t.Errorf("Expected %s got %s", shurl, body)
			}
		})
	}
}
