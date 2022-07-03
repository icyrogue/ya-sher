package idgen

import (
	"math/rand"
	"time"
)

//Get a seed so that ids are random every time
func InitID() {
	rand.Seed(time.Now().UnixMicro())
}
func GenID(data string) string {
	chars := []byte("qwertyuiopasdfghklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM1234567890")

	output := []byte{}
	for i := 0; i != 8; i++ {
		output = append(output, chars[rand.Intn(len(chars))])
	}
	return string(output)
}
