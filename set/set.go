package set

import (
	"encoding/json"
	"iter"
	"maps"
)

type Set[E comparable] struct {
	m map[E]struct{}
}

func (s *Set[E]) Add(el E) bool {
	if s.m == nil {
		s.m = make(map[E]struct{})
	}
	ln := len(s.m)
	s.m[el] = struct{}{}
	return len(s.m) > ln
}

func (s *Set[E]) AddSlice(els ...E) {
	if s.m == nil {
		s.m = make(map[E]struct{})
	}
	for _, el := range els {
		s.m[el] = struct{}{}
	}
}

func (s *Set[E]) Contains(v E) bool {
	_, ok := s.m[v]
	return ok
}

func (s *Set[E]) Len() int {
	return len(s.m)
}

func (s *Set[E]) Equal(s2 Set[E]) bool {
	if len(s.m) != len(s2.m) {
		return false
	}
	for v := range s2.m {
		if _, ok := s.m[v]; !ok {
			return false
		}
	}
	return true
}

func (s *Set[E]) All() iter.Seq[E] {
	return func(yield func(E) bool) {
		for v := range s.m {
			if !yield(v) {
				return
			}
		}
	}
}

func (s *Set[E]) Clone() Set[E] {
	var r Set[E]
	r.m = maps.Clone(s.m)
	return r
}

func (s *Set[E]) Union(sets ...Set[E]) Set[E] {
	r := s.Clone()
	for _, s := range sets {
		for v := range s.m {
			r.Add(v)
		}
	}
	return r
}

func (s *Set[E]) UnmarshalJSON(data []byte) error {
	var els []E
	err := json.Unmarshal(data, &els)
	if err != nil {
		return err
	}
	s.AddSlice(els...)
	return nil
}
