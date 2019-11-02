package lab

import (
	"encoding/json"
	"fmt"
	"github.com/hitman99/autograde/internal/utils"
	"net/http"
	"strings"
	"time"
)

type duration struct {
	time.Duration
}

type labScenario struct {
	Name        string   `json:"name"`
	Cycle       uint     `json:"cycle"`
	Duration    duration `json:"duration"`
	StudentsKey string   `json:"studentsKey"`
	TasksKey    string   `json:"tasksKey"`
}

func (d *duration) UnmarshalJSON(b []byte) (err error) {
	d.Duration, err = time.ParseDuration(strings.Trim(string(b), `"`))
	return
}

type labCtrl struct {
	Action string `json:"action"`
}

func (l *Lab) LabScenarioHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// configure new scenario
	case "POST":
		if r.Body != nil {
			defer r.Body.Close()
			var s labScenario
			err := utils.UnmarshalBody(r.Body, &s)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			labConfig, err := l.kubeClient.GetConfigMap(s.TasksKey, "autograde")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			studs, err := l.redisClient.LRange(s.StudentsKey, 0, -1).Result()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			students := make([]*Student, 0, len(studs))

			for _, st := range studs {
				stud := Student{}
				err := json.Unmarshal([]byte(st), &stud)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				students = append(students, &stud)
			}

			_, tasks, err := GetLabScenarioFromConfig("", labConfig["tasks"])
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			} else {
				l.SetupScenario(s.Name, s.Cycle, s.Duration.Duration, students, tasks)
			}
		} else {
			http.Error(w, "request body is nil", http.StatusBadRequest)
		}
	// start/stop checker
	case "PATCH":
		if r.Body != nil {
			defer r.Body.Close()
			var c labCtrl
			err := utils.UnmarshalBody(r.Body, &c)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			switch c.Action {
			case "start":
				err := l.Run()
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			case "startFromState":
				err := l.loadStateFromRedis()
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			case "stop":
				err := l.Stop()
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "OK")
		} else {
			http.Error(w, "request body is nil", http.StatusBadRequest)
		}
	}
}

func (l *Lab) LabStateHandler(w http.ResponseWriter, r *http.Request) {
	state, err := json.Marshal(l.GetState())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(state)
	}
}
