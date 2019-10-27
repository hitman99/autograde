package lab

import (
	"github.com/hitman99/autograde/internal/state"
	"log"
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
	state        state.State
	logger       *log.Logger

	wg        *sync.WaitGroup
	wgErr     *sync.WaitGroup
	errChan   chan error
	stopChan  <-chan bool
	stateChan <-chan *state.TaskScore
}

type Assignment struct {
	Description string
	Student     *Student
	Tasks       []Task
}

type Student struct {
	FirstName         string
	LastName          string
	DockerhubUsername string
	GithubUsername    string
}
