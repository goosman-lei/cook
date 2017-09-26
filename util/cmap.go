package util

import (
	"sync"
)

type CMap struct {
	sync.RWMutex
	mapping map[string]interface{}
}

func NewCMap() *CMap {
	return &CMap{
		mapping: make(map[string]interface{}),
	}
}

func (m *CMap) Set(k string, v interface{}) {
	m.Lock()
	defer m.Unlock()
	m.mapping[k] = v
}

func (m *CMap) Erase(k string) {
	m.Lock()
	defer m.Unlock()
	delete(m.mapping, k)
}

func (m *CMap) CheckAndErase(k string, fn func(interface{}) bool) bool {
	m.Lock()
	defer m.Unlock()
	if v, e := m.mapping[k]; e && fn(v) {
		delete(m.mapping, k)
		return true
	}
	return false
}

func (m *CMap) Copy(mapping map[string]interface{}) {
	m.Lock()
	defer m.Unlock()
	m.mapping = mapping
}

func (m *CMap) Get(k string) (interface{}, bool) {
	m.RLock()
	defer m.RUnlock()
	v, e := m.mapping[k]
	return v, e
}

func (m *CMap) Has(k string) bool {
	m.RLock()
	defer m.RUnlock()
	_, exists := m.mapping[k]
	return exists
}

func (m *CMap) MustGet(k string) interface{} {
	m.RLock()
	defer m.RUnlock()
	v, exists := m.mapping[k]
	if !exists {
		panic("key not exists")
	}
	return v
}

func (m *CMap) Len() int {
	m.RLock()
	defer m.RUnlock()
	return len(m.mapping)
}

func (m *CMap) Keys() []string {
	m.RLock()
	defer m.RUnlock()
	keys := make([]string, 0, len(m.mapping))
	for k, _ := range m.mapping {
		keys = append(keys, k)
	}
	return keys
}

func (m *CMap) Vals() []interface{} {
	m.RLock()
	defer m.RUnlock()
	vals := make([]interface{}, 0, len(m.mapping))
	for _, v := range m.mapping {
		vals = append(vals, v)
	}
	return vals
}

// unsafe: if you use other cmap method in fn
func (m *CMap) Iterate(fn func(string, interface{})) {
	m.RLock()
	defer m.RUnlock()
	for k, v := range m.mapping {
		fn(k, v)
	}
}

// unsafe: if you use other cmap method in fn
func (m *CMap) Apply(fn func(string, interface{}) interface{}) {
	m.Lock()
	defer m.Unlock()
	for k, v := range m.mapping {
		m.mapping[k] = fn(k, v)
	}
}
