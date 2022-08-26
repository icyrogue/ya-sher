package mlt

import (
	"context"
	"errors"
	"log"
	"regexp"

	"golang.org/x/sync/errgroup"
)

type Storage interface {
	GetAllUserURLs(cookie string) map[string]string
}


type Mlt struct {
	st Storage
	wg *errgroup.Group
	Output chan string
	data chan []string
}

type worker struct {
	id string
	otch chan string
	wg *errgroup.Group
	userurls map[string]string
}

func New(st Storage) *Mlt{
	return &Mlt {
		st: st,
			//	wg: &sync.WaitGroup{},
		data: make(chan []string, 5),

	}
}

func (m *Mlt) Start (ctx context.Context) {
	g, _ := errgroup.WithContext(ctx)
	m.wg = g
	var pop int
	go func() {
	log.Println("started mlt")
	loop:
	for {
	select {
	case data := <- m.data:
		log.Println("mlt got new input:", data)
		pop = len(data) - 1
		data, cookie := data[:pop], data[pop:][0]
		urls := m.st.GetAllUserURLs(cookie)

		for _, dt := range data {
			wk := worker{id: dt, wg: m.wg, otch: m.Output, userurls: urls}
		go wk.Do()
		}
	case <- ctx.Done():
		break loop
	}
	}
		if err := m.wg.Wait(); err != nil {
			log.Println(err.Error())
		}
	close(m.Output)
	}()
}

func (w *worker) Do() error {
	log.Println("worker started", w.id)
	re := regexp.MustCompile(`([A-Z]|[a-z]|[0-9]){8}`)
	if !re.MatchString(w.id) {
		log.Println("regex rejected")
		return errors.New("worker got a nonid string")
	}
	if _, fd := w.userurls[w.id]; !fd {
		err := errors.New("user dosent have such URL ID")
		log.Println(err.Error())
		return err
	}

	w.otch <- w.id
	log.Println("done", w.id)
	return nil
}

func (m *Mlt) GetInput() chan []string {
	return m.data
}


