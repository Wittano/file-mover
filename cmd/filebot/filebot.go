package main

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/wittano/filebot/setting"
	"github.com/wittano/filebot/tasks"
	"github.com/wittano/filebot/watcher"
	"time"
)

func runMainCommand(_ *cobra.Command, _ []string) {
	ctx := context.Background()
	conf, err := setting.Flags.Config()
	if err != nil {
		setting.Logger().Fatal("Failed load configuration", err)
		return
	}

	w := watcher.NewWatcher(ctx)
	defer w.Close()
	w.AddFilesToObservable(conf)

	tasks.Run(ctx, 1*time.Hour, tasks.MoveToTrash)

	go w.UpdateObservableFileList()
	go w.ObserveFiles()

	w.WaitForEvents()
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		setting.Logger().Fatal("Failed to start FileBot", err)
	}
}
