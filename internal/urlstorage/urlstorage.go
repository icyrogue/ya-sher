package urlstorage

import (
	"errors"
	"sync"
)

type storage struct {
	data map[string]string
	mtx  sync.RWMutex
}

func New() *storage {
	return &storage{
		data: make(map[string]string),
	}
}

//AddToStorage: adds url to mock database
func (st *storage) Add(id string, long string) error {
	st.mtx.Lock()
	defer st.mtx.Unlock()
	if _, fd := st.data[id]; fd {
		return errors.New("url with that id already exists")
	}
	st.data[id] = long
	return nil
}
func (st *storage) GetByID(id string) (string, error) {
	st.mtx.RLock()
	defer st.mtx.RUnlock()

	if _, fd := st.data[id]; fd {
		var long = st.data[id]
		return long, nil
	}
	return "", errors.New("no url with such ID")
}
func (st *storage) GetByLong(long string) (string, error) {
	for id, el := range st.data {
		if el == long {
			return id, nil
		}
	}
	return "", errors.New("no url with such id")
}
