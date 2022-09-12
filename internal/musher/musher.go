package musher

import (
	"context"
	"time"
	"strconv"
	"log"
)

type Musher struct {
	Input chan string
	Options *Options
	st Storage
}

type Storage interface {
	BulkDelete(bch []interface{}, args string)
}

type Options struct {
	MaxWaitTime int
	MaxBufferLength int
}


func New(opts *Options, st Storage) *Musher {
	const (
	deffMaxWaitTime = 15
	deffMaxBufferLength = 60
	)
	if opts.MaxWaitTime <= 0 {
		opts.MaxWaitTime = deffMaxWaitTime
	}
	if opts.MaxBufferLength <= 0 {
		opts.MaxBufferLength = deffMaxBufferLength
	}
	inpt := make(chan string, opts.MaxBufferLength)
	return &Musher{Options: opts, Input: inpt, st: st}
}

func (m *Musher) Start(ctx context.Context) {
	bch := make([]interface{}, 0, 5)
	var args = "(ARRAY["
	bl := m.Options.MaxBufferLength
	t := make(chan struct{})
	wt := time.Duration(m.Options.MaxWaitTime) * time.Second
	f := func (){
		t<-struct{}{}
		log.Println("Storage timed out")
	}
	timer := time.AfterFunc(wt, f)
	defer timer.Stop()

		go func (){
		loop:
		for {
		log.Println("Storage is alive")
		select {
		case v := <- m.Input:
				log.Println("Storage got ", v)
				timer.Reset(wt)
				bch = append(bch, v)
				l := len(bch)
			args = args + " $" + strconv.Itoa(l) + ","
			log.Printf("storage is %d/%d", l, bl)
			//Stop if buffer is certain size
			if len(bch) > bl {
				m.st.BulkDelete(bch, args)
			}

		case <- ctx.Done():
			log.Println("Shutting down musher")

			break loop
		case <- t:
			m.st.BulkDelete(bch, args)
		}}

		if l := len(m.Input); l != 0 {
		log.Printf("Buffer has %d elems left!", l)
			for v := range m.Input {
				bch = append(bch, v)
				args = args + " $" + strconv.Itoa(len(bch)) + ","
			}
			m.st.BulkDelete(bch, args)
		}

}()

}
