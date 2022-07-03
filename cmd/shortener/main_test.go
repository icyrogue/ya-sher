package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

type want struct {
	PostCode int
	GetCode  int
	header   []byte
}

var testurl []byte

func Test_main(t *testing.T) {
	tests := []struct {
		name string
		url  []byte
		want want
	}{
		{name: "Simple test #1",
			url: []byte("google.com"),
			want: want{
				PostCode: 201,
				GetCode:  http.StatusTemporaryRedirect,
				header:   []byte("google.com"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := apiInit()
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(tt.url))
			r.ServeHTTP(w, req)

			res := w.Result()
			//Testing Codes\\
			if res.StatusCode != tt.want.PostCode {
				t.Errorf("Expected %d got %d", tt.want.PostCode, res.StatusCode)
			}
			//Testing Responce\\
			defer res.Body.Close()
			body, err := ioutil.ReadAll(res.Body)

			if err != nil {
				t.Error(err)
			}

			re := regexp.MustCompile(`\/.{8}\b`)
			if !re.Match(body) {
				t.Errorf("Invalid ID signature %s", body)

			}
			//Testting GET\\
			w2 := httptest.NewRecorder()
			id := re.Find(body)
			req = httptest.NewRequest(http.MethodGet, string(id), nil)
			r.ServeHTTP(w2, req)
			res = w2.Result()
			//Testing Codes\\
			if res.StatusCode != tt.want.GetCode {
				t.Errorf("Expected %d got %d", tt.want.GetCode, res.StatusCode)
			}

			//Testing Headers\\
			header := res.Header.Get("Location")
			if header != string(tt.want.header) {
				t.Errorf("Expected %s got %s", tt.want.header, header)
			}

		})
	}
}
