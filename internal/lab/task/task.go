package task

import (
	"context"
	"errors"
	"fmt"
	"github.com/hitman99/autograde/internal/config"
	"github.com/hitman99/autograde/internal/dockerhub"
	"github.com/hitman99/autograde/internal/github"
	"github.com/hitman99/autograde/internal/kubernetes"
	"github.com/hitman99/autograde/internal/lab"
	"github.com/hitman99/autograde/internal/rest"
	"log"
	"os"
)

type Maker interface {
	MakeTask(ctx context.Context, def *Definition, s *lab.Student) (Task, error)
}

type maker struct {
	logger          *log.Logger
	githubClient    github.Client
	restClient      rest.Client
	dockerhubClient dockerhub.Client
	kubeClient      kubernetes.Client
}

func NewTaskMaker() Maker {
	logger := log.New(os.Stderr, "[Task Maker] ", log.Ltime)
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
	Eval() (score int, err error)
	GetKind() string
	GetName() string
}

type task struct {
	evaluator func() (bool, error)
	completed bool
	def       *Definition
}

func (t *task) Eval() (int, error) {
	if !t.completed {
		done, err := t.evaluator()
		if err != nil {
			return 0, err
		}
		t.completed = done
		return t.def.Score, nil
	} else {
		return t.def.Score, nil
	}
}

func (t *task) GetKind() string {
	return t.def.Kind
}

func (t *task) GetName() string {
	return t.def.Name
}

type Definition struct {
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

func (m *maker) MakeTask(ctx context.Context, def *Definition, s *lab.Student) (Task, error) {
	switch def.Kind {
	case "github":
		switch def.Name {
		case "checkFork":
			return &task{
				evaluator: func() (bool, error) {
					log.Println("evaluating: " + def.Name + "." + def.Kind)
					return m.githubClient.CheckFork(ctx, s.GithubUsername, def.Config["repo"])
				},
				completed: false,
				def:       def,
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
}
