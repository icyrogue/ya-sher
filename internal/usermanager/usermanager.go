package usermanager

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"hash"
	"sync"
)

type user struct {
	cookie []byte
	urls map[string]string
	//появлется ошибка go assignment copies lock value to usr
}

type UserManager struct {
	users map[string]user
	mac hash.Hash
	key []byte
	mtx sync.RWMutex
}

func (mu *UserManager) AddUserURL(cookie string, url string, id string) error {
	usr := user{}
	var ok bool
	mu.mtx.RLock()
		defer mu.mtx.RUnlock()
	if usr, ok = mu.users[cookie]; !ok {
		return errors.New("no such user")
	}
	usr.urls[id] = url
	return nil
}

func (mu *UserManager) NewUser() (string, error) {
	mu.mtx.Lock()
	defer mu.mtx.Unlock()

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
	return hmac.Equal(org, new)
}

func (mu *UserManager) GetAllUserURLs(cookie string) map[string]string {
	usr := mu.users[cookie]


	res := usr.urls
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
