package main

import (
	"encoding/json"
	"fmt"
	"github.com/hitman99/autograde/internal/lab"
	"time"
)

func main() {

	labScenario := lab.NewLab("test lab", 1, time.Hour,
		[]*lab.Student{{
			FirstName:         "Test",
			LastName:          "The tester",
			DockerhubUsername: "hitman99",
			GithubUsername:    "cloudtr",
		},{
			FirstName:         "Test1",
			LastName:          "The tester1",
			DockerhubUsername: "xxxx",
			GithubUsername:    "autograde",
		},
		},
		[]*lab.TaskDefinition{{
			Name: "checkFork",
			Kind: "github",
			Config: map[string]string{
				"repo": "autograde",
			},
			Description: "github fork checker",
			Score:       1,
		}, {
			Name: "checkFork",
			Kind: "github",
			Config: map[string]string{
				"repo": "autograde",
			},
			Description: "github fork checker",
			Score:       7,
		},
		},
	)

	state := labScenario.GetState()
	bytes, err := json.Marshal(state)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))
	return
	labScenario.Run()

}
