package cmd

import (
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/wittano/filebot/pkg/config"
	"github.com/wittano/filebot/pkg/cron"
	"github.com/wittano/filebot/pkg/watcher"
	"log"
	"path/filepath"
	"time"
)

var (
	Flags   config.Flags
	rootCmd = &cobra.Command{
		Use:   "filebot",
		Short: "Automatically manager your files",
		Long:  "FileBot is simple file manager to automation your boring action with files e.g. moving to trash, another directory or renaming files",
		Run:   runMainCommand,
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&Flags.ConfigPath, "config", "c", getDefaultConfigPath(), "Specific path for filebot configuration")
	rootCmd.PersistentFlags().DurationVarP(&Flags.UpdateInterval, "updateInterval", "u", getDefaultUpdateInterval(), "Set time after filebot should be refresh watched file state")
}

func runMainCommand(_ *cobra.Command, _ []string) {
	conf, err := config.Get(Flags.ConfigPath)
	if err != nil {
		log.Fatalf("Failed loaded configuration: %s", err)
	}

	w := watcher.NewWatcher()
	w.AddFilesToObservable(conf)

	s := cron.NewScheduler()
	s.StartAsync()

	go w.UpdateObservableFileList(Flags)
	go w.ObserveFiles()

	w.WaitForEvents()
}

func getDefaultConfigPath() string {
	homeDir, err := homedir.Dir()
	if err != nil {
		panic(err)
	}

	return filepath.Join(homeDir, ".config", "filebot", "config.toml")
}

func getDefaultUpdateInterval() time.Duration {
	duration, _ := time.ParseDuration("10m")
	return duration
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
