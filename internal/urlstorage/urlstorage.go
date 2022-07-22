package urlstorage

import (
	"bufio"
	"errors"
	"os"
	"sync"
)

type storage struct {
	Data   map[string]string
	mtx    sync.RWMutex
	writer *bufio.Writer
	reader *bufio.Reader
	file   *os.File
}

func New(flPath string) *storage {
	if flPath != "" {
		data, err := recoverData(flPath)
		if err != nil {
			data = make(map[string]string)
		}
		file, err := os.OpenFile(flPath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
		if err != nil {
			return &storage{Data: data}
		}

		return &storage{
			Data:   data,
			writer: bufio.NewWriter(file),
			reader: bufio.NewReader(file),
			file:   file,
		}
	}

	return &storage{
		Data: make(map[string]string),
	}
}

//AddToStorage: adds url to mock database
func (st *storage) Add(id string, long string) error {
	st.mtx.Lock()
	defer st.mtx.Unlock()

	if _, fd := st.Data[id]; fd {
		return errors.New("url with that id already exists")
	}
	st.Data[id] = long
	return nil
}
func (st *storage) GetByID(id string) (string, error) {
	st.mtx.RLock()
	defer st.mtx.RUnlock()

	if _, fd := st.Data[id]; fd {
		var long = st.Data[id]
		return long, nil
	}
	return "", errors.New("no url with such ID")
}
func (st *storage) GetByLong(long string) (string, error) {
	for id, el := range st.Data {
		if el == long {
			return id, nil
		}
	}
	return "", errors.New("no url with such id")
}

func recoverData(flPath string) (map[string]string, error) {
	return nil, nil
}
