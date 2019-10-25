package main

import (
	"github.com/hitman99/autograde/internal/lab"
	"github.com/hitman99/autograde/internal/lab/task"
	"github.com/spf13/viper"
	"time"
)

func main() {

	viper.Set("GIHUB_TOKEN", "token")
	labScenario := lab.NewLab("test lab", 1, time.Hour,
		[]*lab.Student{{
			FirstName:         "Test",
			LastName:          "The tester",
			DockerhubUsername: "hitman99",
			GithubUsername:    "hitman99",
		},
		},
		[]*task.Definition{{
			Name: "checkFork",
			Kind: "github",
			Config: map[string]string{
				"repo": "hlds-docker",
			},
			Description: "github fork checker",
			Score:       1,
		},
		},
	)

	labScenario.Run()

}
