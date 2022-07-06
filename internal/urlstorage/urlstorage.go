package urlstorage

import (
	"errors"
	"sync"
)

type Url struct {
	Short string
	Long  string
	ID    string
}

var mockStorage = []Url{}

type storage struct {
	data map[string]string
	mtx  sync.RWMutex
}

func New() *storage {
	return &storage{data: make(map[string]string)}
}

func NewUrl(long string, short string) Url {
	return Url{
		Short: short,
		Long:  long,
		ID:    short[len(short)-8:],
	}
}

//Add: adds url to mock database
func (st *storage) Add(id string, long string) error {
	st.mtx.Lock()
	defer st.mtx.Unlock()
	if _, ok := st.data[id]; ok {
		return errors.New("already exist")
	}
	st.data[id] = long

	return nil
}

//GetLongByID: returns long version from id
func GetByID(id string) *Url {
	for _, u := range mockStorage {
		if u.ID == id {
			return &u
		}
	}
	return nil
}

//GetByLong: retruns short version by long version
func GetByLong(long string) *Url {
	for _, u := range mockStorage {
		if u.Long == long {
			return &u
		}
	}
	return nil
}
