package ordered

import (
	"iter"
)

type Set[E comparable] Map[E, struct{}]

func (s *Set[E]) Add(el E) bool {
	added := (*Map[E, struct{}])(s).Set(el, struct{}{})
	return added
}

func (s *Set[E]) All() iter.Seq[E] {
	seq := (*Map[E, struct{}])(s).All()
	return func(yield func(el E) bool) {
		seq(func(key E, val struct{}) bool {
			return yield(key)
		})
	}
}

func (s1 *Set[E]) Equal(s2 *Set[E]) bool {
	m1 := (*Map[E, struct{}])(s1)
	m2 := (*Map[E, struct{}])(s2)
	if m1.Len() != m2.Len() {
		return false
	}
	for key := range m1.All() {
		_, ok := m2.Get(key)
		if !ok {
			return false
		}
	}
	return true
}
