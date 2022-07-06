package idgen

import (
	"math/rand"
	"time"

	"github.com/icyrogue/ya-sher/internal/urlstorage"
)

//Get a seed so that ids are random every time
func InitID() {
	rand.Seed(time.Now().UnixMicro())
}

//GenID: generates a new ID until there is no such ID already in database
func GenID() string {
	chars := []byte("qwertyuiopasdfghklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890")
	for {
		output := []byte{}
		for i := 0; i != 8; i++ {
			output = append(output, chars[rand.Intn(len(chars))])
		}
		if urlstorage.GetByID(string(output)) == nil {
			return string(output)
		}
	}
}

type usecase struct {
	storage Storage
}

func New(storage Storage) *usecase {
	return &usecase{storage: storage}
}

func (u *usecase) CreateShortURL(long string) (shurl string, err error) {
	shurl = GenID()
	if err := u.storage.Add(shurl, long); err != nil {
		return "", err
	}
	return shurl, nil
}

type Storage interface {
	Add(id, long string) error
}
