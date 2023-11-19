package main

import (
	"github.com/spf13/cobra"
	"github.com/wittano/filebot/cron"
	"github.com/wittano/filebot/setting"
	"github.com/wittano/filebot/watcher"
	"log"
)

func runMainCommand(_ *cobra.Command, _ []string) {
	conf := setting.Flags.GetConfig()

	w := watcher.NewWatcher()
	w.AddFilesToObservable(conf)

	s := cron.NewScheduler()
	s.StartAsync()
	defer s.Stop()

	go w.UpdateObservableFileList()
	go w.ObserveFiles()

	w.WaitForEvents()
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
