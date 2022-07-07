package api

import (
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Test_api_CrShort(t *testing.T) {
	type fields struct {
		router  *gin.Engine
		logger  *zap.Logger
		opts    *Options
		urlProc URLProcessor
		st      Storage
	}
	tests := []struct {
		name   string
		fields fields
		want   gin.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &api{
				router:  tt.fields.router,
				logger:  tt.fields.logger,
				opts:    tt.fields.opts,
				urlProc: tt.fields.urlProc,
				st:      tt.fields.st,
			}
			if got := a.CrShort(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("api.CrShort() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_api_ReLong(t *testing.T) {
	type fields struct {
		router  *gin.Engine
		logger  *zap.Logger
		opts    *Options
		urlProc URLProcessor
		st      Storage
	}
	tests := []struct {
		name   string
		fields fields
		want  string
	}{ {name: "Simple Test #1",
			}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &api{
				router:  tt.fields.router,
				logger:  tt.fields.logger,
				opts:    tt.fields.opts,
				urlProc: tt.fields.urlProc,
				st:      tt.fields.st,
			}
			if got := a.ReLong(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("api.ReLong() = %v, want %v", got, tt.want)
			}
		})
	}
}
