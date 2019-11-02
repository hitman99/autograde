package cmd

import (
	"github.com/gorilla/mux"
	"github.com/hitman99/autograde/internal/api"
	"github.com/hitman99/autograde/internal/config"
	"github.com/hitman99/autograde/internal/signup"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
	"time"
)

var signupCmd = &cobra.Command{
	Use: "signup",
	Run: func(cmd *cobra.Command, args []string) {
		runSignup()
		os.Exit(0)
	},
}

func runSignup() {
	logger := log.New(os.Stdout, "[signup api] ", log.Ltime)
	sig := signup.NewSignup(logger)
	r := mux.NewRouter()
	r.HandleFunc("/signup", sig.SignupHandler).Methods("POST")
	s := r.PathPrefix("/state").Subrouter()
	s.HandleFunc("/", sig.StateHandler).Methods("GET")
	amw := api.NewAuthMiddleware(config.GetConfig().AdminToken)
	s.Use(amw.Middleware)

	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("frontend/dist"))))

	srv := &http.Server{
		Addr:         "0.0.0.0:80",
		Handler:      r,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
	}
	logger.Println("started http server on port 80")
	log.Fatal(srv.ListenAndServe())
}
