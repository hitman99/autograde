package cmd

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/hitman99/autograde/internal/api"
	"github.com/hitman99/autograde/internal/config"
	"github.com/hitman99/autograde/internal/lab"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
	"time"
)

var apiCmd = &cobra.Command{
	Use: "api",
	Run: func(cmd *cobra.Command, args []string) {
		runLabApi()
		os.Exit(0)
	},
}

func runLabApi() {
	logger := log.New(os.Stdout, "[api] ", log.Ltime)
	labControl := lab.NewLabController()
	r := mux.NewRouter()
	r.HandleFunc("/lab/scenario", labControl.LabScenarioHandler).Methods("POST", "PATCH")
	r.HandleFunc("/lab/deps/{resource}", labControl.LabDependencyHandler).Methods("POST", "DELETE")
	r.HandleFunc("/lab/scenario/state", labControl.LabStateHandler).Methods("GET")
	amw := api.NewAuthMiddleware(config.GetConfig().AdminToken)
	r.Use(amw.Middleware)

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", host, port),
		Handler:      r,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
	}
	logger.Printf("started http server on port %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
