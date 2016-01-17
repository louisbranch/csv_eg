package worker

import (
	"errors"
	"fmt"
	"mime/multipart"
	"sync"
	"time"

	"github.com/luizbranco/csv_eg/internal/parser/csv"
	"github.com/luizbranco/csv_eg/internal/transaction"

	"code.google.com/p/go-uuid/uuid"
)

type Pool struct {
	mutex sync.Mutex
	in    chan struct{}
	queue map[string]*multipart.FileHeader
}

func NewPool(size int) Pool {
	return Pool{
		in:    make(chan struct{}, size),
		queue: make(map[string]*multipart.FileHeader),
		mutex: sync.Mutex{},
	}
}

func (p Pool) Enqueue(h *multipart.FileHeader) string {
	id := uuid.New()
	p.mutex.Lock()
	p.queue[id] = h
	p.mutex.Unlock()
	return id
}

func (p Pool) Retrieve(id string, m csv.Mapping) ([]transaction.Transaction, error) {
	select {
	case p.in <- struct{}{}:
		defer func() {
			<-p.in
		}()

		p.mutex.Lock()
		h, ok := p.queue[id]
		delete(p.queue, id)
		p.mutex.Unlock()

		if !ok || h == nil {
			return nil, fmt.Errorf("File not found, please try again")
		}

		f, err := h.Open()
		if err != nil {
			return nil, fmt.Errorf("File couldn't be opened, please try again")
		}

		return csv.Parse(f, m)
	case <-time.After(1 * time.Minute):
		return nil, errors.New("Request timeout")
	}
}
