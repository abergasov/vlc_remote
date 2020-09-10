package main

import (
	"fmt"
	"log"
	"net/http"
	"vlc_remote/pkg/config"
)

func main() {
	conf := config.InitConf()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World!")
	})
	err := http.ListenAndServe(":"+conf.AppPort, nil)
	if err != nil {
		log.Fatal("error start server", err.Error())
	}
}
