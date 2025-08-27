package cmd

import (
	"log/slog"

	"github.com/spf13/cobra"
)

var verbose bool

var rootCmd = &cobra.Command{
	Use:   "cc-manager",
	Short: "Claude 配置管理工具",
	Long:  "一个用于管理 Claude 配置文件的工具，支持切换不同的配置。",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "显示详细日志")
}

func initConfig() {
	if verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Debug logging enabled")
	}
}