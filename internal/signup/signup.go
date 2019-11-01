package signup

import (
	"encoding/json"
	"github.com/go-redis/redis/v7"
	"github.com/hitman99/autograde/internal/config"
	"github.com/hitman99/autograde/internal/lab"
	"github.com/hitman99/autograde/internal/utils"
	"net/http"
)

type Signup struct {
	redisClient *redis.Client
}

func NewSignup() *Signup {
	rcli := redis.NewClient(&redis.Options{
		Addr:     config.GetConfig().RedisAddress,
		Password: "",
		DB:       0,
	})
	_, err := rcli.Ping().Result()
	if err != nil {
		panic(err)
	}
	return &Signup{redisClient: rcli}
}

func (s *Signup) SignupHandler(w http.ResponseWriter, r *http.Request) {
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
	students, err := s.redisClient.LRange("students", 0, -1).Result()
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

func (s *Signup) persistStudent(stud *lab.Student) error {
	sBlob, err := json.Marshal(stud)
	if err != nil {
		return err
	}
	err = s.redisClient.LPush("students", sBlob).Err()
	if err != nil {
		return err
	}
	return nil
}
