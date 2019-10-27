package state

type State struct {
	Students []Student `json:"students"`
}

type Student struct {
	FirstName string
	LastName  string
	UUID      string
	Scores    []TaskScore
}

type TaskScore struct {
	Description string
	UUID        string
	Value       int
}

func (s *Student) Total() int {
	total := 0
	for _, point := range s.Scores {
		total += point.Value
	}
	return total
}
