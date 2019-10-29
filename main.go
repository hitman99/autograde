package main

import (
    "github.com/hitman99/autograde/internal/lab"
    "time"
)

func main() {

	labScenario := lab.NewLab("test lab", 1, time.Hour,
		[]*lab.Student{{
			FirstName:         "Test",
			LastName:          "The tester",
			DockerhubUsername: "hitman99",
			GithubUsername:    "hitman99",
		},
		},
		[]*lab.TaskDefinition{
		//    {
		//	Name: "checkFork",
		//	Kind: "github",
		//	Config: map[string]string{
		//		"repo": "autograde",
		//	},
		//	Description: "github fork checker",
		//	Score:       1,
		//},
		{
			Name: "checkBuildAction",
			Kind: "github",
			Config: map[string]string{
				"repo": "autograde",
				"buildAction": ".gitignore",
			},
			Description: "github fork checker",
			Score:       7,
		},
		},
	)

	labScenario.Run()

}
