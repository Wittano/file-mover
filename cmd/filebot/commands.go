package main

import (
	"github.com/go-playground/validator/v10"
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
	rootCmd.PersistentFlags().StringVarP(&setting.Flags.ConfigPath, "config", "c", setting.DefaultConfigPath(), "specific path for filebot configuration")
	rootCmd.PersistentFlags().StringVarP(&setting.Flags.LogFilePath, "log", "l", "", "path to log file")
	rootCmd.PersistentFlags().StringVarP(&setting.Flags.LogLevelName, "logLevel", "", "INFO", "log level")
	rootCmd.PersistentFlags().DurationVarP(&setting.Flags.UpdateInterval, "updateInterval", "u", setting.DefaultUpdateInterval(), "set time after filebot should be refresh watched file state")

	v := validator.New(validator.WithRequiredStructEnabled())
	if err := v.Struct(setting.Flags); err != nil {
		setting.Logger().Fatal("Failed to parse flags", err)
	}
}
