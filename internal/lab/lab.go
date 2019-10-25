package lab

import (
	"context"
	"github.com/hitman99/autograde/internal/lab/task"
	"sync"
	"time"
)

func NewLab(name string, cycle uint, duration time.Duration, students []*Student, defs []*task.Definition) *Lab {

	maker := task.NewTaskMaker()

	participants := make([]Assignment, 0, len(students))

	for _, stud := range students {
		tasks := make([]task.Task, 0, len(defs))
		for _, def := range defs {
			studentTask, err := maker.MakeTask(context.TODO(), def, stud)
			if err != nil {
				panic(err)
			}
			tasks = append(tasks, studentTask)
		}
		participants = append(participants, Assignment{
			Description: name,
			Student:     stud,
			Tasks:       tasks,
		})
	}

	return &Lab{
		Name:         name,
		Cycle:        cycle,
		IsRunning:    false,
		IsFinished:   false,
		start:        nil,
		duration:     duration,
		participants: participants,
	}
}

func (l *Lab) Run() error {
	l.wg = &sync.WaitGroup{}
	l.errChan = make(chan error)
	l.stopChan = make(chan bool)
	l.wg.Add(1)

	go func(stop chan<- bool, wg *sync.WaitGroup, a *Assignment) {
		for i := 0; i < 10; i++ {
			a.Tasks[0].Eval()
			time.Sleep(time.Second * 2)
		}
	}(l.stopChan, l.wg, &l.participants[0])

	l.wg.Wait()
	return nil
}

func (l *Lab) Stop() error {
	return nil
}
