package main

import (
	"github.com/wittano/file-mover/src/config"
	"github.com/wittano/file-mover/src/watcher"
	"log"
)

func main() {
	var configPath string

	parseFlags(&configPath)

	conf, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed loaded configuration: %s", err)
	}

	w := watcher.NewWatcher()
	w.AddFileToObservable(conf)

	go w.ObserveFiles()

	w.WaitForEvents()
}
