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
            K8sNamespace:      "sandbox",
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
                    "repo":        "autograde",
                    "buildAction": ".gitignore",
                },
                Description: "github fork checker",
                Score:       7,
            }, {
                Name: "checkContainerImage",
                Kind: "kubernetes",
                Config: map[string]string{
                    "imageName":               "k8s101",
                    "deploymentLabelSelector": "lab=microservices",
                },
                Description: "kubernetes container image check",
                Score:       3,
            },
        },
    )

    labScenario.Run()

}
