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
