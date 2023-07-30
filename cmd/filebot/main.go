package main

import (
	"github.com/wittano/fmanager/pkg/config"
	"github.com/wittano/fmanager/pkg/watcher"
	"log"
)

func main() {
	flags := parseFlags()

	conf, err := config.Load(flags.ConfigPath)
	if err != nil {
		log.Fatalf("Failed loaded configuration: %s", err)
	}

	w := watcher.NewWatcher()
	w.AddFilesToObservable(conf)

	go w.UpdateObservableFileList(flags)
	go w.ObserveFiles()

	w.WaitForEvents()
}
