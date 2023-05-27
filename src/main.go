package main

import (
	"github.com/wittano/file-mover/src/config"
	"github.com/wittano/file-mover/src/watcher"
	"log"
)

func main() {
	flags := config.ParseFlags()

	conf, err := config.LoadConfig(flags.ConfigPath)
	if err != nil {
		log.Fatalf("Failed loaded configuration: %s", err)
	}

	w := watcher.NewWatcher()
	w.AddFilesToObservable(conf)

	go w.UpdateObservableFileList(flags)
	go w.ObserveFiles()

	w.WaitForEvents()
}
