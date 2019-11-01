package main

import "github.com/hitman99/autograde/cmd"

func main() {

	cmd.Execute()
	//logger := log.New(os.Stdout, "[Lab Scenario]", log.Ltime)
	//kClient := kubernetes.MustNewClient()
	//labConfig, err := kClient.GetConfigMap(config.GetConfig().Configmap, "tcentric")
	//if err != nil {
	//	logger.Fatalf("cannot load lab scenario from configmap %s", config.GetConfig().Configmap)
	//}
	//students, tasks, err := lab.GetLabScenarioFromConfig(labConfig["students"], labConfig["tasks"])
	//if err == nil {
	//	labScenario := lab.NewLab("test lab", 1, time.Hour, students, tasks)
	//	st, _ := json.Marshal(labScenario.GetState())
	//	logger.Printf(string(st))
	//	//labScenario.Run()
	//} else {
	//	logger.Fatalf("cannot run scenario, config load failed: %s", err.Error())
	//}
}
