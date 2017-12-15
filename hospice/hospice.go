package hospice

import (
	"reflect"
	"sync"
)

var (
	Wg *sync.WaitGroup = new(sync.WaitGroup)
)

func Go(fn interface{}, args ...interface{}) {
	r_fn := reflect.Indirect(reflect.ValueOf(fn))
	if r_fn.Kind() != reflect.Func {
		return
	}

	r_in := []reflect.Value{}
	for _, arg := range args {
		r_in = append(r_in, reflect.ValueOf(arg))
	}

	Wg.Add(1)
	go func() {
		defer Wg.Done()
		r_fn.Call(r_in)
	}()
}

func Add(delta int) {
	Wg.Add(delta)
}

func Done() {
	Wg.Done()
}

func Wait() {
	Wg.Wait()
}
