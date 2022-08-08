package dbstorage

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type storage struct {
	db *sql.DB
	Options *Options
}

type Options struct {
	DBPath string
}

func New() *storage {
	return &storage{}
}

func (st *storage) Init() {
db, err := sql.Open("postgres", st.Options.DBPath)
	if err != nil {
		log.Fatal(err)
		return
	}
st.db = db
}

func (st *storage) Ping(ctx context.Context) bool {
	if err := st.db.PingContext(ctx); err != nil {
		return false
	}
	return true
}

func (st *storage) Close() {
	st.db.Close()
}

func (st *storage) Add(id string, long string) error {
	return nil
}

func (st *storage) GetByID(id string, ctx context.Context) (string, error) {
	return "", nil
}

func (st *storage) GetByLong(long string, ctx context.Context) (string, error) {
	return "", nil
}
