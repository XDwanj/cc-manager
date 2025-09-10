package cmd

import (
	"cc-manager/internal/config"
	"log/slog"

	"github.com/spf13/cobra"
)

var verbose bool
var globalClientName string

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
	rootCmd.PersistentFlags().StringVar(&globalClientName, "client", "", "指定客户端 (claude/codex/gemini，默认: claude)")
	
	// 为全局 --client 标志添加补全功能
	rootCmd.RegisterFlagCompletionFunc("client", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return config.GetClientNames(), cobra.ShellCompDirectiveNoFileComp
	})
}

func initConfig() {
	if verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Debug logging enabled")
	}
}

// GetGlobalClientName 返回全局 --client 标志的值，如果未设置则返回默认值
func GetGlobalClientName() string {
	if globalClientName == "" {
		return config.GetDefaultClient().Name
	}
	return globalClientName
}