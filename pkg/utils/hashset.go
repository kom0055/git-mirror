package utils

type (
	Empty             struct{}
	Any[K comparable] map[K]Empty
)

func NewAnySet[K comparable](items ...K) Any[K] {
	ss := Any[K]{}
	ss.Insert(items...)
	return ss
}

func AnyKeySet[K comparable, V any](theMap map[K]V) Any[K] {
	ret := NewAnySet[K]()

	for k := range theMap {
		ret.Insert(k)
	}
	return ret
}

func (s Any[K]) Insert(items ...K) Any[K] {
	for _, item := range items {
		s[item] = Empty{}
	}
	return s
}

func (s Any[K]) Delete(items ...K) Any[K] {
	for _, item := range items {
		delete(s, item)
	}
	return s
}

func (s Any[K]) Has(item K) bool {
	_, contained := s[item]
	return contained
}

func (s Any[K]) HasAll(items ...K) bool {
	for _, item := range items {
		if !s.Has(item) {
			return false
		}
	}
	return true
}

func (s Any[K]) HasAny(items ...K) bool {
	for _, item := range items {
		if s.Has(item) {
			return true
		}
	}
	return false
}

func (s Any[K]) Difference(s2 Any[K]) Any[K] {
	result := NewAnySet[K]()
	for key := range s {
		if !s2.Has(key) {
			result.Insert(key)
		}
	}
	return result
}

func (s Any[K]) Union(s2 Any[K]) Any[K] {
	result := NewAnySet[K]()
	for key := range s {
		result.Insert(key)
	}
	for key := range s2 {
		result.Insert(key)
	}
	return result
}

func (s Any[K]) Intersection(s2 Any[K]) Any[K] {
	var walk, other Any[K]
	result := NewAnySet[K]()
	if s.Len() < s2.Len() {
		walk = s
		other = s2
	} else {
		walk = s2
		other = s
	}
	for key := range walk {
		if other.Has(key) {
			result.Insert(key)
		}
	}
	return result
}

// IsSuperset returns true if and only if s1 is a superset of s2.
func (s Any[K]) IsSuperset(s2 Any[K]) bool {
	for item := range s2 {
		if !s.Has(item) {
			return false
		}
	}
	return true
}

func (s Any[K]) Equal(s2 Any[K]) bool {
	return len(s) == len(s2) && s.IsSuperset(s2)
}

func (s Any[K]) List() []K {
	res := make([]K, 0, len(s))
	for key := range s {
		res = append(res, key)
	}
	return res
}

func (s Any[K]) PopAny() (K, bool) {
	for key := range s {
		s.Delete(key)
		return key, true
	}
	var zeroValue K
	return zeroValue, false
}

func (s Any[K]) Len() int {
	return len(s)
}
