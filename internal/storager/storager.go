package storager

import (
	"github.com/icyrogue/ya-sher/internal/api"
	"github.com/icyrogue/ya-sher/internal/config"
	"github.com/icyrogue/ya-sher/internal/dbstorage"
	"github.com/icyrogue/ya-sher/internal/idgen"
	"github.com/icyrogue/ya-sher/internal/musher"
	"github.com/icyrogue/ya-sher/internal/urlstorage"
)

type Storage interface {
api.Storage
idgen.Storage
musher.Storage
Init()
Close()
}

//type Options struct



func Get(cfg config.Cfg) Storage {
if cfg.DBOpts.DBPath == "" {
	str := urlstorage.New()

	str.Options = cfg.StrOpts
	return str
}
	str := dbstorage.New()
	str.Options = cfg.DBOpts
	return str

}

