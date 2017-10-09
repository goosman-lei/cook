package stats

import (
	"bytes"
	"fmt"
	cook_log "gitlab.niceprivate.com/golang/cook/log"
	"sort"
	"strings"
	"sync"
	"time"
)

type t_msg struct {
	o uint8
	k string
	v interface{}
}

const (
	MSG_QUEUE_LEN = 8192
	KEY_SEPERATOR = "."
)

const (
	OP_SET = iota + 1
	OP_INCR
	OP_DECR
)

var (
	mapping          map[string]interface{} = make(map[string]interface{})
	mq               chan *t_msg            = make(chan *t_msg, MSG_QUEUE_LEN)
	Eye              string
	refresh_interval time.Duration = time.Second
	map_mutex        *sync.RWMutex = new(sync.RWMutex)

	wg     *sync.WaitGroup = new(sync.WaitGroup)
	doneCh chan bool       = make(chan bool)
)

func Init_stats() {
	update()
	wg.Add(1)
	go func() {
		defer wg.Done()
		var (
			m  *t_msg
			ok bool
		)
	SelectMsgLoop:
		for {
			select {
			case m, ok = <-mq:
				if !ok {
					break SelectMsgLoop
				}

				switch m.o {
				case OP_SET:
					set(m.k, m.v)
				case OP_INCR:
					incr(m.k, m.v)
				case OP_DECR:
					decr(m.k, m.v)
				}
			}
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-doneCh:
				return
			case <-time.After(refresh_interval):
				update()
			}
		}
	}()
}

func Done_stats() {
	close(mq)
	close(doneCh)
	wg.Wait()
}

func SetRefreshInterval(interval time.Duration) {
	refresh_interval = interval
}

func update() {
	map_mutex.RLock()
	defer map_mutex.RUnlock()
	buf := new(bytes.Buffer)
	dump_map_to_buffer(mapping, buf)
	Eye = buf.String()
}

func dump_map_to_buffer(m map[string]interface{}, buf *bytes.Buffer) {
	var (
		isFirst, iIsFirst bool             = true, true
		keys              sort.StringSlice = make([]string, 0, len(m))
		key               string
		val               interface{}
	)

	buf.WriteString("{")
	for key, _ = range m {
		keys = append(keys, key)
	}
	keys.Sort()
	for _, key = range keys {
		val = m[key]
		if isFirst {
			isFirst = false
			fmt.Fprintf(buf, "%s: ", key)
		} else {
			fmt.Fprintf(buf, ", %s: ", key)
		}

		switch v := val.(type) {
		case int, uint32, int32, uint64, int64:
			fmt.Fprintf(buf, "%d", v)
		case float32, float64:
			fmt.Fprintf(buf, "%0.3f", v)
		case bool:
			fmt.Fprintf(buf, "%t", v)
		case string:
			fmt.Fprintf(buf, "\"%s\"", v)
		case func() string:
			fmt.Fprintf(buf, "\"%s\"", v())
		case []int:
			buf.WriteString("[")
			iIsFirst = true
			for _, e := range v {
				if iIsFirst {
					iIsFirst = false
					fmt.Fprintf(buf, "%d", e)
				} else {
					fmt.Fprintf(buf, ", %d", e)
				}
			}
			buf.WriteString("]")
		case []uint32:
			buf.WriteString("[")
			iIsFirst = true
			for _, e := range v {
				if iIsFirst {
					iIsFirst = false
					fmt.Fprintf(buf, "%d", e)
				} else {
					fmt.Fprintf(buf, ", %d", e)
				}
			}
			buf.WriteString("]")
		case []int32:
			buf.WriteString("[")
			iIsFirst = true
			for _, e := range v {
				if iIsFirst {
					iIsFirst = false
					fmt.Fprintf(buf, "%d", e)
				} else {
					fmt.Fprintf(buf, ", %d", e)
				}
			}
			buf.WriteString("]")
		case []uint64:
			buf.WriteString("[")
			iIsFirst = true
			for _, e := range v {
				if iIsFirst {
					iIsFirst = false
					fmt.Fprintf(buf, "%d", e)
				} else {
					fmt.Fprintf(buf, ", %d", e)
				}
			}
			buf.WriteString("]")
		case []int64:
			buf.WriteString("[")
			iIsFirst = true
			for _, e := range v {
				if iIsFirst {
					iIsFirst = false
					fmt.Fprintf(buf, "%d", e)
				} else {
					fmt.Fprintf(buf, ", %d", e)
				}
			}
			buf.WriteString("]")
		case []float32:
			buf.WriteString("[")
			iIsFirst = true
			for _, e := range v {
				if iIsFirst {
					iIsFirst = false
					fmt.Fprintf(buf, "%0.3f", e)
				} else {
					fmt.Fprintf(buf, ", %0.3f", e)
				}
			}
			buf.WriteString("]")
		case []float64:
			buf.WriteString("[")
			iIsFirst = true
			for _, e := range v {
				if iIsFirst {
					iIsFirst = false
					fmt.Fprintf(buf, "%0.3f", e)
				} else {
					fmt.Fprintf(buf, ", %0.3f", e)
				}
			}
			buf.WriteString("]")
		case []bool:
			buf.WriteString("[")
			iIsFirst = true
			for _, e := range v {
				if iIsFirst {
					iIsFirst = false
					fmt.Fprintf(buf, "%t", e)
				} else {
					fmt.Fprintf(buf, ", %t", e)
				}
			}
			buf.WriteString("]")
		case []string:
			buf.WriteString("[")
			iIsFirst = true
			for _, e := range v {
				if iIsFirst {
					iIsFirst = false
					fmt.Fprintf(buf, "\"%s\"", e)
				} else {
					fmt.Fprintf(buf, ", \"%s\"", e)
				}
			}
			buf.WriteString("]")
		case map[string]interface{}:
			dump_map_to_buffer(v, buf)
		default:
			fmt.Fprintf(buf, "<unknown>")
		}
	}
	buf.WriteString("}")
}

func get_last_map_and_key(key string) (wrapper map[string]interface{}, data_key string, err error) {
	var (
		i        int
		k        string
		ok       bool
		key_eles []string = strings.Split(key, KEY_SEPERATOR)
		last_key          = len(key_eles) - 1
	)

	wrapper = mapping
	data_key = key_eles[last_key]

	// root key
	if last_key <= 0 {
		return wrapper, key, nil
	}

	map_mutex.RLock()
	// find sub map util its value is nil
	for i, k = range key_eles[:last_key] {
		// have no this value, go to fill
		if _, ok = wrapper[k]; !ok {
			map_mutex.RUnlock()
			goto FillNilSubKeys
		}
		// wrong type
		if wrapper, ok = wrapper[k].(map[string]interface{}); !ok {
			map_mutex.RUnlock()
			return nil, "", fmt.Errorf("%s is not map[string]interface{}. but want key: %s", strings.Join(key_eles[:i+1], KEY_SEPERATOR), key)
		}
	}
	map_mutex.RUnlock()
	return wrapper, data_key, nil

FillNilSubKeys:
	map_mutex.Lock()
	defer map_mutex.Unlock()
	// fill remain sub key
	for i, k = range key_eles[i:last_key] {
		wrapper[k] = make(map[string]interface{})
		wrapper = wrapper[k].(map[string]interface{})
	}

	return wrapper, data_key, nil
}

// ignore type compare of new value and old value
func set(key string, val interface{}) {
	if wrapper, data_key, err := get_last_map_and_key(key); err == nil {
		map_mutex.Lock()
		defer map_mutex.Unlock()
		wrapper[data_key] = val
	}
}

func incr(key string, val interface{}) {
	defer func() {
		if r := recover(); r != nil {
			cook_log.Infof("stats decr panic: %q", r)
		}
	}()
	var (
		data_key string
		wrapper  map[string]interface{}
		err      error
	)

	if wrapper, data_key, err = get_last_map_and_key(key); err != nil {
		return
	}

	map_mutex.Lock()
	defer map_mutex.Unlock()
	if wrapper[data_key] == nil {
		wrapper[data_key] = val
		return
	}

	switch v := val.(type) {
	case int:
		wrapper[data_key] = wrapper[data_key].(int) + v
	case uint32:
		wrapper[data_key] = wrapper[data_key].(uint32) + v
	case int32:
		wrapper[data_key] = wrapper[data_key].(int32) + v
	case uint64:
		wrapper[data_key] = wrapper[data_key].(uint64) + v
	case int64:
		wrapper[data_key] = wrapper[data_key].(int64) + v
	case float32:
		wrapper[data_key] = wrapper[data_key].(float32) + v
	case float64:
		wrapper[data_key] = wrapper[data_key].(float64) + v
	}
	// ignore other type
}

func decr(key string, val interface{}) {
	defer func() {
		if r := recover(); r != nil {
			cook_log.Infof("stats decr panic: %q", r)
		}
	}()
	var (
		data_key string
		wrapper  map[string]interface{}
		err      error
	)

	if wrapper, data_key, err = get_last_map_and_key(key); err != nil {
		return
	}

	map_mutex.Lock()
	defer map_mutex.Unlock()
	if wrapper[data_key] == nil {
		wrapper[data_key] = val
		return
	}

	switch v := val.(type) {
	case int:
		wrapper[data_key] = wrapper[data_key].(int) - v
	case uint32:
		wrapper[data_key] = wrapper[data_key].(uint32) - v
	case int32:
		wrapper[data_key] = wrapper[data_key].(int32) - v
	case uint64:
		wrapper[data_key] = wrapper[data_key].(uint64) - v
	case int64:
		wrapper[data_key] = wrapper[data_key].(int64) - v
	case float32:
		wrapper[data_key] = wrapper[data_key].(float32) - v
	case float64:
		wrapper[data_key] = wrapper[data_key].(float64) - v
	}
	// ignore other type
}

func Set_int(key string, val int) {
	mq <- &t_msg{o: OP_SET, k: key, v: val}
}

func Set_uint32(key string, val uint32) {
	mq <- &t_msg{o: OP_SET, k: key, v: val}
}

func Set_int32(key string, val int32) {
	mq <- &t_msg{o: OP_SET, k: key, v: val}
}

func Set_uint64(key string, val uint64) {
	mq <- &t_msg{o: OP_SET, k: key, v: val}
}

func Set_int64(key string, val int64) {
	mq <- &t_msg{o: OP_SET, k: key, v: val}
}

func Set_float32(key string, val float32) {
	mq <- &t_msg{o: OP_SET, k: key, v: val}
}

func Set_float64(key string, val float64) {
	mq <- &t_msg{o: OP_SET, k: key, v: val}
}

func Set_string(key string, val string) {
	mq <- &t_msg{o: OP_SET, k: key, v: val}
}

func Set_bool(key string, val bool) {
	mq <- &t_msg{o: OP_SET, k: key, v: val}
}

func Set_fn(key string, val func() string) {
	mq <- &t_msg{o: OP_SET, k: key, v: val}
}

func Set_ints(key string, val []int) {
	mq <- &t_msg{o: OP_SET, k: key, v: val}
}

func Set_uint32s(key string, val []uint32) {
	mq <- &t_msg{o: OP_SET, k: key, v: val}
}

func Set_int32s(key string, val []int32) {
	mq <- &t_msg{o: OP_SET, k: key, v: val}
}

func Set_uint64s(key string, val []uint64) {
	mq <- &t_msg{o: OP_SET, k: key, v: val}
}

func Set_int64s(key string, val []int64) {
	mq <- &t_msg{o: OP_SET, k: key, v: val}
}

func Set_float32s(key string, val []float32) {
	mq <- &t_msg{o: OP_SET, k: key, v: val}
}

func Set_float64s(key string, val []float64) {
	mq <- &t_msg{o: OP_SET, k: key, v: val}
}

func Set_strings(key string, val []string) {
	mq <- &t_msg{o: OP_SET, k: key, v: val}
}

func Incr_int(key string, val int) {
	mq <- &t_msg{o: OP_INCR, k: key, v: val}
}

func Incr_uint32(key string, val uint32) {
	mq <- &t_msg{o: OP_INCR, k: key, v: val}
}

func Incr_int32(key string, val int32) {
	mq <- &t_msg{o: OP_INCR, k: key, v: val}
}

func Incr_uint64(key string, val uint64) {
	mq <- &t_msg{o: OP_INCR, k: key, v: val}
}

func Incr_int64(key string, val int64) {
	mq <- &t_msg{o: OP_INCR, k: key, v: val}
}

func Incr_float32(key string, val float32) {
	mq <- &t_msg{o: OP_INCR, k: key, v: val}
}

func Incr_float64(key string, val float64) {
	mq <- &t_msg{o: OP_INCR, k: key, v: val}
}

func Decr_int(key string, val int) {
	mq <- &t_msg{o: OP_DECR, k: key, v: val}
}

func Decr_uint32(key string, val uint32) {
	mq <- &t_msg{o: OP_DECR, k: key, v: val}
}

func Decr_int32(key string, val int32) {
	mq <- &t_msg{o: OP_DECR, k: key, v: val}
}

func Decr_uint64(key string, val uint64) {
	mq <- &t_msg{o: OP_DECR, k: key, v: val}
}

func Decr_int64(key string, val int64) {
	mq <- &t_msg{o: OP_DECR, k: key, v: val}
}

func Decr_float32(key string, val float32) {
	mq <- &t_msg{o: OP_DECR, k: key, v: val}
}

func Decr_float64(key string, val float64) {
	mq <- &t_msg{o: OP_DECR, k: key, v: val}
}
