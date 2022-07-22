package urlstorage

import (
	"bufio"
	"errors"
	"log"
	"os"
	"strings"
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
			log.Println("Couldnt recover data from file, creating new")
			data = make(map[string]string)
		}
		file, err := os.OpenFile(flPath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
		if err != nil {
			log.Println("Couldnt open storage file, runing in RAM mode")
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

	if st.file != nil {
		data := []byte(id + " " + long + "\n")
		if _, err := st.writer.Write(data); err != nil {
			return err
		}
		return st.writer.Flush()
	}

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

//recoverData: tries to recover urls from previus session
func recoverData(flPath string) (map[string]string, error) {
	file, err := os.Open(flPath)
	if err != nil {
		return nil, err
	}
	scaner := bufio.NewScanner(file)
	data := make(map[string]string)

	for scaner.Scan() {
		el := strings.Split(string(scaner.Bytes()), " ")
		data[el[0]] = el[1]
	}
	if len(data) == 0 {
		return nil, errors.New("file is empty")
	}

	return data, nil
}
