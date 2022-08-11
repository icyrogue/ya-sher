package idgen

import (
	"log"
	"math/rand"
	"time"
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
	log.Println("staeted generation")
	shurl = genID()
	err = u.st.Add(shurl, long)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return shurl, nil
}
func New(storage Storage) *usecase {
	return &usecase{st: storage}
}

type Storage interface {
	Add(id, long string) error
}
