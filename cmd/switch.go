package cmd

import (
	"cc-manager/internal/config"
	"log/slog"

	"github.com/spf13/cobra"
)

var switchClaude bool

var switchCmd = &cobra.Command{
	Use:   "switch [配置名]",
	Short: "切换 Claude 配置",
	Long: `切换到指定的 Claude 配置。
例如：cc-manager switch deepseek       (切换 settings 配置)
     cc-manager switch linus --claude  (切换 CLAUDE 配置)

默认情况下切换 settings.json 链接，使用 --claude 标志切换 CLAUDE.md 链接。`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		configName := args[0]
		
		manager := config.NewManager()
		
		var err error
		var configType string
		
		if switchClaude {
			slog.Info("Switching to CLAUDE configuration", "config", configName)
			err = manager.SwitchClaude(configName)
			configType = "CLAUDE"
		} else {
			slog.Info("Switching to settings configuration", "config", configName)
			err = manager.SwitchConfig(configName)
			configType = "settings"
		}
		
		if err != nil {
			slog.Error("Failed to switch configuration", "config", configName, "type", configType, "error", err)
			return
		}
		
		slog.Info("Successfully switched to configuration", "config", configName, "type", configType)
	},
}

func init() {
	switchCmd.Flags().BoolVar(&switchClaude, "claude", false, "切换 CLAUDE 配置而非 settings 配置")
	rootCmd.AddCommand(switchCmd)
}