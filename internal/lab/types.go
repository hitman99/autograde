package lab

import (
	"github.com/hitman99/autograde/internal/lab/task"
	"sync"
	"time"
)

type Lab struct {
	Name         string
	Cycle        uint
	enabled      bool
	IsRunning    bool
	IsFinished   bool
	start        *time.Time
	duration     time.Duration
	participants []Assignment

	wg       *sync.WaitGroup
	errChan  <-chan error
	stopChan chan<- bool
}

type Assignment struct {
	Description string
	Student     *Student
	Tasks       []task.Task
}

type Student struct {
	FirstName         string
	LastName          string
	DockerhubUsername string
	GithubUsername    string
}

type Definition struct {
	Name string
}
