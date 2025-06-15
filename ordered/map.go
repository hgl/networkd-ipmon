package ordered

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"iter"
)

type Map[K comparable, V any] struct {
	m    map[K]*entry[K, V]
	head *entry[K, V]
	tail *entry[K, V]
}

type entry[K comparable, V any] struct {
	next, prev *entry[K, V]
	key        K
	value      V
}

func (m *Map[K, V]) Set(k K, v V) bool {
	if m.m == nil {
		m.m = make(map[K]*entry[K, V])
	} else {
		e, ok := m.m[k]
		if ok {
			e.value = v
			return false
		}
	}
	e := &entry[K, V]{key: k, value: v}
	m.m[k] = e
	if m.head == nil {
		m.head = e
		m.tail = e
	} else {
		e.prev = m.tail
		m.tail.next = e
		m.tail = e
	}
	return true
}

func (m *Map[K, V]) Get(k K) (V, bool) {
	e, ok := m.m[k]
	if ok {
		return e.value, true
	}
	var v V
	return v, false
}

func (m *Map[K, V]) Index(k K) int {
	e := m.m[k]
	if e == nil {
		return -1
	}
	i := 0
	needle := e
	for needle != m.head {
		needle = needle.prev
		i += 1
	}
	return i
}

func (m *Map[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for e := m.head; e != nil; e = e.next {
			if !yield(e.key, e.value) {
				return
			}
		}
	}
}

func (m *Map[K, V]) Len() int {
	return len(m.m)
}

func (m *Map[K, V]) UnmarshalJSON(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	token, err := dec.Token()
	if err == io.EOF {
		return errors.New("unexpected EOF")
	}
	if t, ok := token.(json.Delim); !ok || t != '{' {
		return fmt.Errorf("json: expect { got: %v", token)
	}
	for dec.More() {
		var k K
		err = dec.Decode(&k)
		if err != nil {
			return err
		}
		var v V
		err = dec.Decode(&v)
		if err != nil {
			return err
		}
		m.Set(k, v)
	}
	return nil
}
