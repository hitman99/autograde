package signup

import (
    "encoding/json"
    "github.com/go-redis/redis/v7"
	"github.com/gorilla/mux"
	"github.com/hitman99/autograde/internal/config"
    "github.com/hitman99/autograde/internal/kubernetes"
    "github.com/hitman99/autograde/internal/lab"
    "github.com/hitman99/autograde/internal/utils"
    "log"
    "net/http"
)

type Signup struct {
    redisClient *redis.Client
    kubeClient  kubernetes.Client
    logger      *log.Logger
}

func NewSignup(logger *log.Logger) *Signup {
    rcli := redis.NewClient(&redis.Options{
        Addr:     config.GetConfig().RedisAddress,
        Password: "",
        DB:       0,
    })
    _, err := rcli.Ping().Result()
    if err != nil {
        panic(err)
    }
    return &Signup{redisClient: rcli, logger: logger, kubeClient: kubernetes.MustNewClient()}
}

func (s *Signup) SignupHandler(w http.ResponseWriter, r *http.Request) {
    s.logger.Printf("%s to /signup", r.Method)
    if r.Body != nil {
        defer r.Body.Close()
        var stud lab.Student
        err := utils.UnmarshalBody(r.Body, &stud)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
        }
        err = s.persistStudent(&stud)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        } else {
            w.WriteHeader(http.StatusOK)
        }

    } else {
        http.Error(w, "request body is nil", http.StatusBadRequest)
    }
}

func (s *Signup) StateHandler(w http.ResponseWriter, r *http.Request) {
    students, err := s.redisClient.LRange(config.GetConfig().RedisStudKey, 0, -1).Result()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    } else {
        studs, err := json.Marshal(students)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        } else {
            w.WriteHeader(http.StatusOK)
            w.Write(studs)
        }
    }
}

func (s *Signup) KubeconfigHandler(w http.ResponseWriter, r *http.Request) {
	githubUsername := mux.Vars(r)["githubUsername"]
	kubeconfig, err := s.kubeClient.GetKubeconfig("ktu-stud-" + githubUsername)
    if err != nil {
       http.Error(w, err.Error(), http.StatusInternalServerError)
    } else {
           w.WriteHeader(http.StatusOK)
           w.Write([]byte(kubeconfig))
    }
}

func (s *Signup) persistStudent(stud *lab.Student) error {
    sBlob, err := json.Marshal(stud)
    if err != nil {
        return err
    }
    err = s.redisClient.LPush(config.GetConfig().RedisStudKey, sBlob).Err()
    if err != nil {
        return err
    }
    return nil
}
