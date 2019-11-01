package lab

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v7"
	"github.com/hitman99/autograde/internal/config"
	"github.com/hitman99/autograde/internal/kubernetes"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

func NewLabController() *Lab {
	rcli := redis.NewClient(&redis.Options{
		Addr:     config.GetConfig().RedisAddress,
		Password: "",
		DB:       0,
	})
	return &Lab{
		logger:      log.New(os.Stderr, "[Lab]", log.Ltime),
		redisClient: rcli,
		kubeClient:  kubernetes.MustNewClient(),
	}
}

func NewLabFromState(state []byte) (*Lab, error) {
	var labState LabState
	err := json.Unmarshal(state, &labState)
	if err != nil {
		return nil, err
	}
	maker := NewTaskMaker()
	participants := make([]Assignment, 0, len(labState.Assignments))

	for _, assig := range labState.Assignments {
		tasks := make([]Task, 0, len(assig.Tasks))
		for _, tsk := range assig.Tasks {
			studentTask, err := maker.MakeTaskFromState(context.TODO(), tsk, assig.Student)
			if err != nil {
				panic(err)
			}
			tasks = append(tasks, studentTask)
		}
		participants = append(participants, Assignment{
			Description: assig.Description,
			Student:     assig.Student,
			Tasks:       tasks,
		})
	}

	isRunning := (labState.Started != nil && labState.Started.Before(time.Now().Add(labState.Duration)))
	return &Lab{
		Name:         labState.Name,
		Cycle:        labState.Cycle,
		IsRunning:    isRunning,
		IsFinished:   !isRunning,
		start:        labState.Started,
		duration:     labState.Duration,
		participants: participants,
		logger:       log.New(os.Stderr, "[Lab]", log.Ltime),
	}, nil
}

func (l *Lab) Run() error {
	now := time.Now()
	l.wg = &sync.WaitGroup{}
	l.wgErr = &sync.WaitGroup{}
	l.errChan = make(chan error)
	l.stopChan = make(chan bool)
	l.IsRunning = true
	l.IsFinished = false
	if l.start == nil {
		l.start = &now
	}

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
	for _, a := range l.participants {
		l.wg.Add(1)
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
				time.Sleep(config.GetConfig().CheckInterval)
			}

		}(l.stopChan, l.errChan, l.wg, &a)
	}

	go func(stop <-chan bool) {
		for {
			select {
			case <-stop:
				return
			default:
				//do nothing

			}
			time.Sleep(time.Minute)
			l.saveStateToRedis()
			if l.start.Add(l.duration).After(time.Now()) {
				l.Stop()
			}
		}

	}(l.stopChan)
	return nil
}

func (l *Lab) Stop() error {
	if l.IsRunning {
		close(l.stopChan)
		l.wg.Wait()
		close(l.errChan)
		l.wgErr.Wait()
		l.IsRunning = false
		l.IsFinished = true
		l.start = nil
	}
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

func (l *Lab) setState(state string) error {
	var labState LabState
	err := json.Unmarshal([]byte(state), &labState)
	if err != nil {
		return err
	}
	maker := NewTaskMaker()
	participants := make([]Assignment, 0, len(labState.Assignments))

	for _, assig := range labState.Assignments {
		tasks := make([]Task, 0, len(assig.Tasks))
		for _, tsk := range assig.Tasks {
			studentTask, err := maker.MakeTaskFromState(context.TODO(), tsk, assig.Student)
			if err != nil {
				panic(err)
			}
			tasks = append(tasks, studentTask)
		}
		participants = append(participants, Assignment{
			Description: assig.Description,
			Student:     assig.Student,
			Tasks:       tasks,
		})
	}

	isRunning := (labState.Started != nil && labState.Started.Before(time.Now().Add(labState.Duration)))

	l.Name = labState.Name
	l.Cycle = labState.Cycle
	l.start = labState.Started
	l.duration = labState.Duration
	l.participants = participants
	if isRunning {
		return l.Run()
	}
	return nil
}

func (a *Assignment) GetState() *AssignmentState {
	tasks := make([]*TaskState, 0, len(a.Tasks))
	for _, task := range a.Tasks {
		tasks = append(tasks, task.GetState())
	}
	return &AssignmentState{
		Description: a.Description,
		Student:     a.Student,
		Tasks:       tasks,
	}
}

func (l *Lab) SetupScenario(name string, cycle uint, duration time.Duration, students []*Student, defs []*TaskDefinition) {
	maker := NewTaskMaker()

	participants := make([]Assignment, 0, len(students))

	for _, stud := range students {
		tasks := make([]Task, 0, len(defs))
		for _, def := range defs {
			studentTask, err := maker.NewTask(context.TODO(), def, stud)
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

	l.participants = participants
	l.Name = name
	l.Cycle = cycle
	l.duration = duration
}

func GetLabScenarioFromConfig(students, tasks string) ([]*Student, []*TaskDefinition, error) {
	var (
		stud     = []*Student{}
		taskDefs = []*TaskDefinition{}
	)
	if students != "" {
		s, _ := strconv.Unquote(string(students))
		err := json.Unmarshal([]byte(s), &stud)
		if err != nil {
			return nil, nil, err
		}
	}
	if tasks != "" {
		err := yaml.Unmarshal([]byte(tasks), &taskDefs)
		if err != nil {
			return nil, nil, err
		}
	}
	return stud, taskDefs, nil
}

func (l *Lab) saveStateToRedis() {
	state := l.GetState()
	stateStr, err := json.Marshal(state)
	if err != nil {
		l.logger.Printf("failed to save lab state: %s", err.Error())
	} else {
		err := l.redisClient.Set("labState", stateStr, time.Hour*12).Err()
		if err != nil {
			l.logger.Printf("failed to save lab state: %s", err.Error())
		}
	}
}

func (l *Lab) loadStateFromRedis() error {
	state, err := l.redisClient.Get("labState").Result()
	if err != nil {
		return err
	}
	return l.setState(state)
}
