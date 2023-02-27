package main

import "github.com/wittano/file-mover/src/watcher"

func main() {
	var configPath string

	parseFlags(&configPath)

	w := watcher.NewWatcher()

	go w.ObserveFiles()

	w.WaitForEvents()
}
