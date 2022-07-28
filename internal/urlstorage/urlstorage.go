package urlstorage

import (
	"bufio"
	"errors"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/icyrogue/ya-sher/internal/config"
)

type storage struct {
	data    map[string]string
	mtx     sync.RWMutex
	writer  *bufio.Writer
	reader  *bufio.Reader
	file    *os.File
	Options *config.StrOpts
}

func New() *storage {
	return &storage{}
}

func (st *storage) Init() {
	flPath := ""
	if st.Options != nil {
		flPath = st.Options.Filepath
	}
	if flPath != "" {
		data, err := recoverData(flPath)
		if err != nil {
			log.Println("Couldnt recover data from file, creating new")
			data = make(map[string]string)
		}
		file, err := os.OpenFile(flPath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
		//defer file.Close() /*TODO: возможно стоит добавить отдельную функцию к структуре, котрая бы закрывала файл и вынести ее в мейн*/
		if err != nil {
			log.Println("Couldnt open storage file, runing in RAM mode")
			st.data = data
			return
		}

		st.data = data
		st.file = file
		st.writer = bufio.NewWriter(file)
		st.reader = bufio.NewReader(file)
		return
	}
	st.data = make(map[string]string)
}

func (st *storage) Close() {
	st.file.Close()
}

//AddToStorage: adds url to mock database
func (st *storage) Add(id string, long string) error {
	st.mtx.Lock()
	defer st.mtx.Unlock()

	if _, fd := st.data[id]; fd {
		return errors.New("url with that id already exists")
	}
	st.data[id] = long

	if st.file != nil {
		data := []byte(id + " " + long + "\n")
		if _, err := st.writer.Write(data); err != nil {
			return err
		}
		return st.writer.Flush()
	}

	return nil
}

//TODO: возможно стоит заставить и все остольные функции storage взаимодействовать с файлом, тогда можно не выгружать все в память, но, с другой стороны, у нас там дальше нужно будет базу данных добавить, поэтому я не знаю
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
