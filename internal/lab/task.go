package lab

import (
	"context"
	"errors"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/hitman99/autograde/internal/config"
	"github.com/hitman99/autograde/internal/dockerhub"
	"github.com/hitman99/autograde/internal/github"
	"github.com/hitman99/autograde/internal/hashes"
	"github.com/hitman99/autograde/internal/kubernetes"
	"github.com/hitman99/autograde/internal/rest"
	"log"
	"os"
)

type Maker interface {
	NewTask(ctx context.Context, def *TaskDefinition, s *Student) (Task, error)
	MakeTaskFromState(ctx context.Context, state *TaskState, s *Student) (Task, error)
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
	GetState() *TaskState
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

func (t *task) GetState() *TaskState {
	return &TaskState{
		Def:       t.def,
		Completed: t.completed,
		UUID:      t.uuid,
	}
}

func unknownKindError(kind string) error {
	return errors.New(fmt.Sprintf("undefined task kind: %s", kind))
}

func unknownNameError(name string) error {
	return errors.New(fmt.Sprintf("undefined task name: %s", name))
}

func (m *maker) makeTask(ctx context.Context, def *TaskDefinition, s *Student, taskUUID *string, completed *bool) (Task, error) {
	var (
		newUUID uuid.UUID
		err     error
	)
	if taskUUID == nil {
		newUUID, err = uuid.NewV4()
		if err != nil {
			return nil, err
		}
	} else {
		newUUID = uuid.FromStringOrNil(*taskUUID)
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
				uuid:      newUUID.String(),
				logger:    m.logger,
			}, nil
		case "checkBuildAction":
			return &task{
				evaluator: func() (bool, error) {
					return m.githubClient.CheckBuildAction(ctx, s.GithubUsername, def.Config["repo"], def.Config["actionName"])
				},
				def:    def,
				uuid:   newUUID.String(),
				logger: m.logger,
			}, nil
		default:
			return nil, unknownNameError(def.Name)
		}
	case "dockerhub":
		switch def.Name {
		case "checkRepo":
			return &task{
				evaluator: func() (bool, error) {
					return m.dockerhubClient.CheckRepo(fmt.Sprintf("%s/%s", s.DockerhubUsername, def.Config["repo"]))
				},
				def:    def,
				uuid:   newUUID.String(),
				logger: m.logger,
			}, nil
		case "checkTags":
			return &task{
				evaluator: func() (bool, error) {
					return m.dockerhubClient.CheckTags(fmt.Sprintf("%s/%s", s.DockerhubUsername, def.Config["repo"]))
				},
				def:    def,
				uuid:   newUUID.String(),
				logger: m.logger,
			}, nil
		default:
			return nil, unknownNameError(def.Name)
		}
	case "kubernetes":
		switch def.Name {
		case "checkContainerImage":
			return &task{
				evaluator: func() (bool, error) {
					return m.kubeClient.CheckContainerImage(s.K8sNamespace, def.Config["deploymentLabelSelector"], fmt.Sprintf("%s/%s", s.DockerhubUsername, def.Config["imageName"]))
				},
				def:    def,
				uuid:   newUUID.String(),
				logger: m.logger,
			}, nil
		case "checkEndpointExists":
			return &task{
				evaluator: func() (bool, error) {
					return m.restClient.CheckEndpointExists(fmt.Sprintf("http://%s.%s/%s", def.Config["serviceName"], s.K8sNamespace, s.GithubUsername))
				},
				def:    def,
				uuid:   newUUID.String(),
				logger: m.logger,
			}, nil
		case "checkEndpointResult":
			return &task{
				evaluator: func() (bool, error) {
					sha256Username, err := hashes.GetHash("sha256", s.GithubUsername)
					if err != nil {
						return false, err
					}
					return m.restClient.CheckEndpointResult(fmt.Sprintf("http://%s.%s/%s", def.Config["serviceName"], s.K8sNamespace, s.GithubUsername), *sha256Username)
				},
				def:    def,
				uuid:   newUUID.String(),
				logger: m.logger,
			}, nil
		default:
			return nil, unknownNameError(def.Name)
		}
	default:
		return nil, unknownKindError(def.Kind)
	}
}

func (m *maker) NewTask(ctx context.Context, def *TaskDefinition, s *Student) (Task, error) {
	return m.makeTask(ctx, def, s, nil, nil)
}

func (m *maker) MakeTaskFromState(ctx context.Context, t *TaskState, s *Student) (Task, error) {
	return m.makeTask(ctx, t.Def, s, &t.UUID, &t.Completed)
}
