package main

import (
	"github.com/spf13/cobra"
	"github.com/wittano/filebot/setting"
)

var (
	rootCmd = &cobra.Command{
		Use:   "filebot",
		Short: "Automatically manager your files",
		Long:  "FileBot is simple file manager to automation your boring action with files e.g. moving to trash, another directory or renaming files",
		Run:   runMainCommand,
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&setting.Flags.ConfigPath, "setting", "c", setting.DefaultConfigPath(), "Specific path for filebot configuration")
	rootCmd.PersistentFlags().StringVarP(&setting.Flags.LogFilePath, "log", "l", "", "Path to log file")
	rootCmd.PersistentFlags().StringVarP(&setting.Flags.LogLevelName, "logLevel", "ll", "INFO", "Log level")
	rootCmd.PersistentFlags().DurationVarP(&setting.Flags.UpdateInterval, "updateInterval", "u", setting.DefaultUpdateInterval(), "Set time after filebot should be refresh watched file state")
}
