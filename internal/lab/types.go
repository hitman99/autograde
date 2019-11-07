package lab

import (
	"github.com/go-redis/redis/v7"
	"github.com/hitman99/autograde/internal/kubernetes"
	"log"
	"sync"
	"time"
)

type Lab struct {
	Name         string
	Cycle        string
	enabled      bool
	IsRunning    bool
	IsFinished   bool
	start        *time.Time
	duration     time.Duration
	participants []Assignment
	logger       *log.Logger

	wg          *sync.WaitGroup
	wgErr       *sync.WaitGroup
	errChan     chan error
	stopChan    chan bool
	redisClient *redis.Client
	kubeClient  kubernetes.Client
}

type Assignment struct {
	Description string
	Student     *Student
	Tasks       []Task
	Score       int
}

type Student struct {
	FirstName         string `yaml:"firstName" json:"firstName"`
	LastName          string `yaml:"lastName" json:"lastName"`
	DockerhubUsername string `yaml:"dockerhubUsername" json:"dockerhubUsername"`
	GithubUsername    string `yaml:"githubUsername" json:"githubUsername"`
	K8sNamespace      string `yaml:"k8sNamespace" json:"k8sNamespace"`
}

type TaskDefinition struct {
	Name        string            `yaml:"name" json:"name"`
	Kind        string            `yaml:"kind" json:"kind"`
	Config      map[string]string `yaml:"config" json:"config"`
	Description string            `yaml:"description" json:"description"`
	Score       int               `yaml:"score" json:"score"`
}

type TaskState struct {
	Def       *TaskDefinition `json:"taskDefinition"`
	Completed bool            `json:"completed"`
	UUID      string          `json:"uuid"`
}

type AssignmentState struct {
	Description string       `json:"description"`
	Student     *Student     `json:"student"`
	Tasks       []*TaskState `json:"tasks"`
}

type LabState struct {
	Name        string             `json:"name"`
	Cycle       string               `json:"cycle"`
	Started     *time.Time         `json:"started"`
	Duration    time.Duration      `json:"duration"`
	Assignments []*AssignmentState `json:"assignments"`
}
