package state

type Student struct {
	FirstName string
	LastName  string
	Points    []Point
}

type Point struct {
	Description string
	Value       int
}

func (s *Student) Total() int {
	total := 0
	for _, point := range s.Points {
		total += point.Value
	}
	return total
}
