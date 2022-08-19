package mlt

import (
	"context"
	"errors"
	"log"

	"golang.org/x/sync/errgroup"
)

type Storage interface {
	GetByID(id string, ctx context.Context) (string, error)
}

//type storage struct {}

type Mlt struct {
	st Storage
	wg *errgroup.Group
	output chan string
	data chan []string
	state bool
}

type worker struct {
	id string
	st Storage
	otch chan string
	wg *errgroup.Group
}

func New(st Storage) *Mlt{
	return &Mlt {
		st: st,
			//	wg: &sync.WaitGroup{},
		output: make(chan string, 50),
		data: make(chan []string, 5),

	}
}

func (m *Mlt) Start (ctx context.Context) {
	g, _ := errgroup.WithContext(ctx)
	m.wg = g
	go func() {
	log.Println("started mlt")
	m.state = true
	loop:
	for {
	select {
	case data := <- m.data:
		for _, dt := range data {
			log.Println(data)
			wk := worker{id: dt, st: m.st, wg: m.wg, otch: m.output}
		go wk.Do()
		}
	case <- ctx.Done():
		break loop
	}
	}
		if err := m.wg.Wait(); err != nil {
			log.Println(err.Error())
		}
	m.state = false
	}()
}

func (w *worker) Do() error {

	if _, err := w.st.GetByID(w.id, context.Background()); err != nil {
		return errors.New("worket didnt find ID")
	}
	if w.id == "" {
		return errors.New("worker got a nonid string")
	}

	w.otch <- w.id
	return nil
}

func (m *Mlt) GetInput() chan []string {
	return m.data
}

func (m *Mlt) GetOutput() chan string {
	return m.output
}

func (m *Mlt) GetState() bool {
	return m.state
}
