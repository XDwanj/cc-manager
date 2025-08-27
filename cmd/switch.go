package cmd

import (
	"cc-manager/internal/config"
	"log/slog"

	"github.com/spf13/cobra"
)

var switchCmd = &cobra.Command{
	Use:   "switch [配置名]",
	Short: "切换 Claude 配置",
	Long: `切换到指定的 Claude 配置。
例如：cc-manager switch deepseek
这会将 settings.json 链接到 settings.deepseek.json`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		configName := args[0]
		
		manager := config.NewManager()
		
		slog.Info("Switching to configuration", "config", configName)
		
		if err := manager.SwitchConfig(configName); err != nil {
			slog.Error("Failed to switch configuration", "config", configName, "error", err)
			return
		}
		
		slog.Info("Successfully switched to configuration", "config", configName)
	},
}

func init() {
	rootCmd.AddCommand(switchCmd)
}