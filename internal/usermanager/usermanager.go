package usermanager

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
)

type user struct {
	cookie []byte
	urls map[string]string
}

type UserManager struct {
	users map[string]user
	mac hash.Hash
	key []byte

}

func (mu *UserManager) AddUserURL(cookie string, url string, id string) error {
	usr := user{}
	var ok bool
	fmt.Println(cookie)

	if usr, ok = mu.users[cookie]; !ok {
		fmt.Println(mu.users)
		return errors.New("no such user")
	}
	usr.urls[id] = url
	fmt.Println("hey", usr.urls)
	return nil
}

func (mu *UserManager) NewUser() (string, error) {
	ckeRaw := make([]byte, 8)
	_, err := rand.Read(ckeRaw)
	if err != nil {
		return "", err
	}
	cke := make([]byte, hex.EncodedLen(len(ckeRaw)))
	hex.Encode(cke, ckeRaw)
	mu.mac.Write(cke)

	cookie := string(cke)

	usr := user{
		cookie: mu.mac.Sum(nil),
		urls: make(map[string]string),
	}
	mu.users[cookie] = usr
	return  cookie, nil
}

func (mu *UserManager) CheckValid(cookie string) bool {
	cke := []byte(cookie)
	mu.mac.Write(cke)
	new := mu.mac.Sum(nil)
	org := mu.users[cookie].cookie
	if hmac.Equal(org, new) {
		return true
	}
	fmt.Println("awdawdawdawd")
	return false
}

func (mu *UserManager) GetAllUserURLs(cookie string) map[string]string {
	fmt.Println("hey" + cookie)
	res := mu.users[cookie].urls
	fmt.Print(res)
	return res
}

func New() (*UserManager, error) {
	key := make([]byte, 16)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	mac := hmac.New(sha256.New, key)
	return &UserManager{
		mac: mac,
		key: key,
	users: make(map[string]user),}, nil
}
