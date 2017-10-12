package pool

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrClosed     = errors.New("pool already closed")
	ErrTimeout    = errors.New("pool get timeout")
	ErrNotWrapper = errors.New("is not wrapper for pool object")
)

type Pool struct {
	len int
	cap int

	objs  chan interface{}
	mutex *sync.Mutex

	factory    func() (interface{}, error)
	destructor func(interface{})
}

func new_pool(c int, factory func() (interface{}, error), destructor func(interface{})) *Pool {
	return &Pool{
		len: 0,
		cap: c,

		mutex: new(sync.Mutex),
		objs:  make(chan interface{}, c),

		factory:    factory,
		destructor: destructor,
	}
}

func (p *Pool) get_timeout(duration time.Duration) (interface{}, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	var (
		obj      interface{}
		err      error
		deadline time.Time
	)
	if p.objs == nil {
		return nil, ErrClosed
	}

	if duration > 0 {
		deadline = time.Now().Add(duration)
	}

	select {
	case obj = <-p.objs:
		if obj == nil {
			return nil, ErrClosed
		}
		return obj, nil
	default:
		if p.len < p.cap {
			goto Product_obj
		} else {
			goto Wait_obj
		}
	}

Product_obj:
	if obj, err = p.factory(); err != nil {
		return nil, err
	} else {
		p.len++
		// who product, who first own it
		return obj, nil
	}

Wait_obj:
	if deadline.IsZero() {
		select {
		case obj = <-p.objs:
			if obj == nil {
				return nil, ErrClosed
			}
		}
	} else {
		select {
		case obj = <-p.objs:
			if obj == nil {
				return nil, ErrClosed
			}
		case <-time.After(deadline.Sub(time.Now())):
			return nil, ErrTimeout
		}
	}
	return obj, nil
}

func (p *Pool) get() (interface{}, error) {
	return p.get_timeout(0)
}

func (p *Pool) put(obj interface{}) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.objs == nil {
		p.destructor(obj)
		return
	}

	select {
	case p.objs <- obj:
	default:
		// pool is full, destruct it
		p.destructor(obj)
	}
}

func (p *Pool) Close() {
	p.mutex.Lock()

	objs := p.objs
	p.objs = nil

	p.mutex.Unlock()

	if objs == nil {
		return
	}

	close(objs)
	for obj := range objs {
		p.destructor(obj)
		p.len--
	}
}

func (p *Pool) Len() int {
	return p.len
}

func (p *Pool) Cap() int {
	return p.cap
}

func (p *Pool) IsClosed() bool {
	return p.objs == nil
}
