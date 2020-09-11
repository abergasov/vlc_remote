package main

import (
	"fmt"
	"log"
	"net/http"
	"vlc_remote/pkg/config"
	"vlc_remote/pkg/logger"
	"vlc_remote/pkg/vlc"
)

func main() {
	logger.InitLogger()
	conf := config.InitConf()
	//fp := files.InitFileParser(conf.MediaRoot)
	//fp.BuildFilesTree()

	player := vlc.Init(conf.VlcHost, conf.VlcPort, conf.VlcLogin, conf.VlcPassword)
	print(player.IsAlive())
	print(player.Pause())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World!")
	})
	err := http.ListenAndServe(":"+conf.AppPort, nil)
	if err != nil {
		log.Fatal("error start server", err.Error())
	}
}
