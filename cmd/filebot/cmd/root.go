package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wittano/filebot/pkg/config"
	"github.com/wittano/filebot/pkg/watcher"
	"log"
)

var (
	flags   config.Flags
	rootCmd = &cobra.Command{
		Use:   "filebot",
		Short: "Automatically manager your files",
		Long:  "FileBot is simple file manager to automation your boring action with files e.g. moving to trash, another directory or renaming files",
		Run:   runMainCommand,
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&flags.ConfigPath, "config", "c", config.GetDefaultConfigPath(), "Specific path for filebot configuration")
	rootCmd.PersistentFlags().DurationVarP(&flags.UpdateInterval, "updateInterval", "u", config.GetDefaultUpdateInterval(), "Set time after filebot should be refresh watched file state")
}

func runMainCommand(_ *cobra.Command, _ []string) {
	conf, err := config.Get(flags.ConfigPath)
	if err != nil {
		log.Fatalf("Failed loaded configuration: %s", err)
	}

	w := watcher.NewWatcher()
	w.AddFilesToObservable(conf)

	go w.UpdateObservableFileList(flags)
	go w.ObserveFiles()

	w.WaitForEvents()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
