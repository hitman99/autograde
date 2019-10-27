package lab

import (
	"context"
	"errors"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/hitman99/autograde/internal/config"
	"github.com/hitman99/autograde/internal/dockerhub"
	"github.com/hitman99/autograde/internal/github"
	"github.com/hitman99/autograde/internal/kubernetes"
	"github.com/hitman99/autograde/internal/rest"
	"log"
	"os"
)

type Maker interface {
	MakeTask(ctx context.Context, def *TaskDefinition, s *Student) (Task, error)
}

type maker struct {
	logger          *log.Logger
	githubClient    github.Client
	restClient      rest.Client
	dockerhubClient dockerhub.Client
	kubeClient      kubernetes.Client
}

func NewTaskMaker() Maker {
	logger := log.New(os.Stderr, "[Task] ", log.Ltime)
	cfg := config.GetConfig()
	return &maker{
		logger:          logger,
		githubClient:    github.MustNewClient(cfg.GithubToken),
		restClient:      rest.MustNewClient(),
		dockerhubClient: dockerhub.MustNewClient(),
		kubeClient:      kubernetes.MustNewClient(),
	}
}

type Task interface {
	Eval() error
	GetKind() string
	GetName() string
	GetUUID() string
	GetScore() int
	IsCompleted() bool
}

type task struct {
	evaluator func() (bool, error)
	completed bool
	uuid      string
	def       *TaskDefinition
	logger    *log.Logger
}

func (t *task) Eval() error {
	if !t.completed {
		t.logger.Printf("evaluating task: %s.%s", t.def.Kind, t.def.Name)
		done, err := t.evaluator()
		if err != nil {
			return err
		}
		if done {
			t.logger.Printf("eval done: %s.%s, score: %d", t.def.Kind, t.def.Name, t.def.Score)
			t.completed = done
		} else {
			t.logger.Printf("eval not done: %s.%s", t.def.Kind, t.def.Name)
		}
		return nil
	} else {
		t.logger.Printf("task completed: %s.%s, score: %d", t.def.Kind, t.def.Name, t.def.Score)
		return nil
	}
}

func (t *task) GetKind() string {
	return t.def.Kind
}

func (t *task) GetName() string {
	return t.def.Name
}

func (t *task) GetUUID() string {
	return t.uuid
}

func (t *task) GetScore() int {
	if t.completed {
		return t.def.Score
	} else {
		return 0
	}
}

func (t *task) IsCompleted() bool {
	return t.completed
}

type TaskDefinition struct {
	Name        string
	Kind        string
	Config      map[string]string
	Description string
	Score       int
}

func unknownKindError(kind string) error {
	return errors.New(fmt.Sprintf("undefined task kind: %s", kind))
}

func unknownNameError(name string) error {
	return errors.New(fmt.Sprintf("undefined task name: %s", name))
}

func (m *maker) MakeTask(ctx context.Context, def *TaskDefinition, s *Student) (Task, error) {
	taskUUID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	switch def.Kind {
	case "github":
		switch def.Name {
		case "checkFork":
			return &task{
				evaluator: func() (bool, error) {
					return m.githubClient.CheckFork(ctx, s.GithubUsername, def.Config["repo"])
				},
				completed: false,
				def:       def,
				uuid:      taskUUID.String(),
				logger:    m.logger,
			}, nil
		case "checkBuildAction":
		default:
			return nil, unknownNameError(def.Name)
		}
	case "dockerhub":
		switch def.Name {
		case "checkRepo":
		case "checkTags":
		default:
			return nil, unknownNameError(def.Name)
		}
	case "kubernetes":
		switch def.Name {
		case "checkPodImage":
		default:
			return nil, unknownNameError(def.Name)
		}
	case "rest":
		switch def.Name {
		case "checkEndpointExists":
		case "checkEndpointResult":
		default:
			return nil, unknownNameError(def.Name)
		}
	default:
		return nil, unknownKindError(def.Kind)
	}
	return nil, nil
}
