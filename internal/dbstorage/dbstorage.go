package dbstorage

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

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
	_, err = db.Exec(`CREATE TABLE urls("id" TEXT, "long" TEXT, "token" TEXT, "deleted" BOOL DEFAULT FALSE);`) //TODO возможно нужна какая то проверка, если таблица
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
	row, err := st.db.QueryContext(ctx, `SELECT long, deleted FROM urls WHERE id = $1`, id)
	if err != nil {
		return "", err
	}
	var ot string
	var del bool
	for row.Next() {
		err = row.Scan(&ot, &del)
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
		if del {
			return ot, errors.New("url gone")
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

func (st *storage) BulkDelete(ctx context.Context, cancel context.CancelFunc, otch chan string) {
	bch := make([]interface{}, 0, 5)
	var args = "("
	log.Println("started storage")

	timer := time.AfterFunc(time.Duration(10)*time.Second, cancel)
	defer timer.Stop()
	go func (){
		loop:
		for {
		select {
		case v := <- otch:
				log.Println("Storage got ", v)
				timer.Reset(time.Duration(10)*time.Second)
				bch = append(bch, v)
				l := len(bch)
			args = args + " $" + strconv.Itoa(l) + ","
			log.Printf("storage is %d/%d", l, 60)
			if len(bch) > 60 {
				timer.Stop()
				break loop
			}

		case <- ctx.Done():
			log.Println("time ran out")

			break loop
		}}

		if l := len(otch); l != 0 {
		log.Printf("Buffer has %d elems left!", l)
		check:
			for v := range otch {
				bch = append(bch, v)
				args = args + " $" + strconv.Itoa(len(bch)) + ","
				if len(otch) == 0 {
					break check
				}
			}
		}

		log.Println("Commit to db pending")


		stmt, err := st.db.Prepare("UPDATE urls SET deleted = TRUE WHERE id IN" + args[:len(args)-1] + ")")
		if err != nil {
			log.Println(err.Error())
			return
		}
		_, err = stmt.Exec(bch...)
		if err != nil {
			log.Println(err.Error())
			return
		}
		log.Printf("Deleted %d elems in db", len(bch))
	}()
}
