package set

type Set struct {
	data map[string]interface{}
}

func NewSet() *Set {
	s := Set{
		data: make(map[string]interface{}),
	}
	return &s
}

func (s *Set) Add(toAdd string) bool {
	if _, exists := s.data[toAdd]; !exists {
		s.data[toAdd] = true
		return true
	}
	return false
}

func (s *Set) AddAll(toAdd []string) {
	for _, str := range toAdd {
		s.Add(str)
	}
}

func (s *Set) AsSlice() []string {
	toReturn := make([]string, len(s.data))
	i := 0
	for key, _ := range s.data {
		toReturn[i] = key
		i++
	}
	return toReturn
}
