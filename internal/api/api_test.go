package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/icyrogue/ya-sher/internal/idgen"
	"github.com/icyrogue/ya-sher/internal/urlstorage"
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
				log.Fatalln(err)
			}
			storage := urlstorage.New("")
			usecase := idgen.New(storage)
			api := New(logger, &Options{Hostname: "http://localhost:8080"}, usecase, storage)
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
			body, err1 := ioutil.ReadAll(res.Body)
			if err1 != nil {
				t.Error(err1)
			}
			shurl, err2 := storage.GetByLong(tt.want)
			if err2 != nil {
				t.Error(err.Error())
				return
			}
			body = body[len(body)-8:]
			if string(body) != shurl {
				t.Errorf("Expected %s got %s", shurl, body)
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
				log.Fatalln(err)
			}
			storage := urlstorage.New("")
			usecase := idgen.New(storage)
			api := New(logger, &Options{Hostname: "http://localhost:8080"}, usecase, storage)
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
				log.Fatalln(err)
			}
			storage := urlstorage.New("")
			usecase := idgen.New(storage)
			api := New(logger, &Options{Hostname: "http://localhost:8080"}, usecase, storage)
			api.Init()
			//Testing POST itself
			w := httptest.NewRecorder()
			wantJSON := struct {
				Want string `json:"url"`
			}{
				Want: tt.want,
			}

			reqJSON, errJ := json.Marshal(wantJSON)

			if errJ != nil {
				log.Fatal(errJ)
			}
			req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer(reqJSON))
			api.router.ServeHTTP(w, req)

			res := w.Result()

			if res.StatusCode != tt.wantedCode {
				t.Errorf("Expected %d got %d", tt.wantedCode, res.StatusCode)
			}
			defer res.Body.Close()
			body, err1 := ioutil.ReadAll(res.Body)
			if err1 != nil {
				t.Error(err1)
			}
			bodyJSON := struct {
				URL string `json:"result"`
			}{}

			if errJ2 := json.Unmarshal(body, &bodyJSON); errJ2 != nil {
				log.Fatal(errJ2)
			}
			shurl, err2 := storage.GetByLong(tt.want)
			if err2 != nil {
				t.Error(err.Error())
				return
			}
			body = []byte(bodyJSON.URL)
			body = body[len(body)-8:]
			if string(body) != shurl {
				t.Errorf("Expected %s got %s", shurl, body)
			}
		})
	}
}
