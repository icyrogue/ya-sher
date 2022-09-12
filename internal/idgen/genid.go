package idgen

import (
	"log"
	"math/rand"
	"time"

	"github.com/icyrogue/ya-sher/internal/jsonmodels"
)

//Get a seed so that ids are random every time
func init() {
	rand.Seed(time.Now().UnixMicro())
}

//GenID: generates a new ID until there is no such ID already in database
func genID() string {
	chars := []byte("qwertyuiopasdfghklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890")
	output := []byte{}
	for i := 0; i != 8; i++ {
		output = append(output, chars[rand.Intn(len(chars))])
	}
	return string(output)
}

type usecase struct {
	st Storage
}

func (u *usecase) CreateShortURL(long string) (shurl string, err error) {
	shurl = genID()
	err = u.st.Add(shurl, long)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return shurl, nil
}

func (u usecase) BulkCreation(data []jsonmodels.JSONBulkInput, baseURL string) ([]jsonmodels.JSONBulkInput, error) {
	for i := range data{
		el := &data[i]
		el.URL = el.Long
		el.Long = ""
		el.Short = (baseURL + "/" + genID())
	}
	/*Cкорей всго правильнее это сделать, добавляя новый эле-
	  мент в слайс, но тогда нужно создавать слайс и
	  каждый раз делать append, не знаю, что эффективнее
	  А еще так можно сделать все с одной структурой!*/
	log.Println("adding urls to storage:", data)
	if err := u.st.BulkAdd(data); err != nil {
		return []jsonmodels.JSONBulkInput{}, err
	}
	return data, nil
}

func New(storage Storage) *usecase {
	return &usecase{st: storage}
}

type Storage interface {
	Add(id, long string) error
	BulkAdd(data []jsonmodels.JSONBulkInput) error
}
