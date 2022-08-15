package dbstorage

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strings"

	"github.com/icyrogue/ya-sher/internal/jsonmodels"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/stdlib"
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
	driverConfig := stdlib.DriverConfig{
			ConnConfig: pgx.ConnConfig{
				//				PreferSimpleProtocol: true, //стейтменты не работают с этой опцией
			},
		}
	stdlib.RegisterDriverConfig(&driverConfig)
	/*Я вот эту вот всю штуку взял из кода автотестов, потому что у меня два дня не открывалась бд при тестах,
	  проблема была в парсинге пути, но код автотестов же умные дяденьки написали, они знают, как лучше сделать*/

	db, err := sql.Open("pgx", driverConfig.ConnectionString(st.Options.DBPath))
	if err != nil {
		log.Fatal(err)
		return
	}
	_, err = db.Exec(`CREATE TABLE urls("id" TEXT, "long" TEXT, "token" TEXT);`) //TODO возможно нужна какая то проверка, если таблица
	if err != nil && !strings.Contains(err.Error(), "already exists") {												 // уже существует
		log.Println(err)
		return
	}
st.db = db
}

//Ping: returns true if db is avilible
func (st *storage) Ping(ctx context.Context) bool {
	return st.db.PingContext(ctx) == nil
}

func (st *storage) Close() {
	st.db.Close()
}

func (st *storage) Add(id string, long string) error {
	_, err := st.db.Exec(`INSERT INTO urls(id, long) VALUES($1, $2)`, id, long ) //TODO возможно сделать тоже самое с транзакциями,
	if err != nil {																// которые заготавливаются в Init()
		log.Println(err)
		return err
	}
return nil
}

func (st *storage) GetByID(id string, ctx context.Context) (string, error) {
	row, err := st.db.QueryContext(ctx, `SELECT long FROM urls WHERE id = $1`, id)
	if err != nil {
		return "", err
	}
	var ot string
	for row.Next() {
		err = row.Scan(&ot)
	if err != nil {
		return "", err
	}
	if err = row.Err(); err != nil {
		return "", err
	}
	if ot == "" {
		return "", errors.New("no such url") //возможно лучше возвращать еще count и если он 0, то кидать эту ошибку
											//но у нас всегда одна строчка, поэтому стоит ли оно того
	}
	}
	return ot, nil
}

func (st *storage) GetByLong(long string, ctx context.Context) (string, error) {
	row, err := st.db.QueryContext(ctx, `SELECT id FROM urls WHERE long = $1`, long)
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

//BulkAdd: add multimple urls to db
func (st *storage) BulkAdd(data []jsonmodels.JSONBulkInput) error {
	tx, err := st.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare("INSERT INTO urls(id, long, token) VALUES($1, $2, $3);")
	if err != nil {
		log.Println(err)
		return err
	}
	defer stmt.Close()
	for _, el := range(data) {
		_, err := stmt.Exec(el.Short[len(el.Short)-8:], el.URL, el.CrlID)
		if err != nil {
		log.Println(err)
		return err
	}
	}
	if err := tx.Commit(); err != nil {
    log.Println("update drivers: unable to commit: ", err)
    return err
    }
	return nil
}
