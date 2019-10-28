package lab

import (
	"context"
	"log"
	"os"
	"sync"
	"time"
)

func NewLab(name string, cycle uint, duration time.Duration, students []*Student, defs []*TaskDefinition) *Lab {

	maker := NewTaskMaker()

	participants := make([]Assignment, 0, len(students))

	for _, stud := range students {
		tasks := make([]Task, 0, len(defs))
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
		logger:       log.New(os.Stderr, "[Lab]", log.Ltime),
	}
}

func (l *Lab) Run() error {
	l.wg = &sync.WaitGroup{}
	l.wgErr = &sync.WaitGroup{}
	l.errChan = make(chan error)
	l.stopChan = make(chan bool)
	l.wg.Add(1)

	// one go routine per student
	l.wgErr.Add(1)
	go func(errs <-chan error, wg *sync.WaitGroup) {
		defer wg.Done()
		for err := range errs {
			log.Printf("error: %s", err.Error())
			select {
			case <-errs:
				break
			}
		}
	}(l.errChan, l.wgErr)

	go func(stop <-chan bool, errs chan<- error, wg *sync.WaitGroup, a *Assignment) {
		defer wg.Done()

		for {
			select {
			case <-stop:
				return
			default:
				//do nothing
			}

			allFinished := true
			scores := 0
			for _, t := range a.Tasks {
				scores += t.GetScore()
				if !t.IsCompleted() {
					if err := t.Eval(); err != nil {
						errs <- err
					}
					allFinished = false
				}
			}
			a.Score = scores
			// no need to evaluate further, all finished
			if allFinished {
				l.logger.Printf("lab completed by: %s %s, score: %d", a.Student.FirstName, a.Student.LastName, scores)
				break
			}
			time.Sleep(time.Second * 2)
		}

	}(l.stopChan, l.errChan, l.wg, &l.participants[0])

	l.wg.Wait()
	close(l.errChan)
	l.wgErr.Wait()
	return nil
}

func (l *Lab) Stop() error {
	return nil
}

func (l *Lab) GetState() *LabState {
	assignments := make([]*AssignmentState, 0, len(l.participants))
	for _, a := range l.participants {
		assignments = append(assignments, a.GetState())
	}
	return &LabState{
		Name:        l.Name,
		Cycle:       l.Cycle,
		Started:     l.start,
		Duration:    l.duration,
		Assignments: assignments,
	}
}

func (a *Assignment) GetState() *AssignmentState {
	tasks := make([]*TaskState, 0, len(a.Tasks))
	for _, task := range a.Tasks {
		tasks = append(tasks, task.GetState())
	}
	return &AssignmentState{
		Student: a.Student,
		Tasks:   tasks,
	}
}
