package dbstorage

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v4"
)

type storage struct {
	db *pgx.Conn
	Options *Options
}

type Options struct {
	DBPath string
}

func New() *storage {
	return &storage{}
}

func (st *storage) Init() {
db, err := pgx.Connect(context.Background(), "postgres://postgres:postgres@postgres:5432/praktikum?sslmode=disable")
	if err != nil {
		log.Fatal(err)
		return
	}
	_, err = db.Exec(context.Background(), `CREATE TABLE urls("id" TEXT, "long" TEXT);`)
	if err != nil {
		log.Fatal(err)
		return
	}
st.db = db
}

func (st *storage) Ping(ctx context.Context) bool {
	if err := st.db.Ping(ctx); err != nil {
		return false
	}
	return true
}

func (st *storage) Close() {
	st.db.Close(context.Background())
}

func (st *storage) Add(id string, long string) error {
	_, err := st.db.Exec(context.Background(), `INSERT INTO urls(id, long) VALUES($1, $2)`, id, long )
	if err != nil {
		println(err)
		return err
	}
// 	rows, err := st.db.Query(`SELECT * FROM urls`)
// 	if err != nil {
// println(err)}
// 	var str string
// 	for rows.Next() {
// 		rows.Scan(&str)
// 		println(str)
// 	}
	return nil
}

func (st *storage) GetByID(id string, ctx context.Context) (string, error) {
	row, err := st.db.Query(ctx, `SELECT long FROM urls WHERE id = $1`, id)
	if err != nil {
		return "", err
	}
	var ot string
	for row.Next() {
		err = row.Scan(&ot)
	if err != nil {
		return "", err
	}
	if ot == "" {
		return "", errors.New("no such url")
	}
	}
	return ot, nil
}

func (st *storage) GetByLong(long string, ctx context.Context) (string, error) {
	row, err := st.db.Query(ctx, `SELECT id FROM urls WHERE long = $1`, long)
	if err != nil {
		return "", err
	}
	var ot string
	for row.Next() {
		err = row.Scan(&ot)
		if err != nil {
			return "", err
		}
	}
		if err = row.Err(); err != nil {
			return "", err
		}
	if ot == "" {
		return "", errors.New("no such url")
	}
	return ot, nil
}
