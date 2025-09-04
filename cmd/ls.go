package cmd

import (
	"cc-manager/internal/config"
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
)

var showDetail bool
var claudeConfig bool

var lsCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "列出所有可用的 Claude 配置",
	Long: `列出 ~/.claude 目录下所有可用的配置文件。
当前激活的配置会用 * 标记。

例如：cc-manager ls              (列出 settings 配置)
     cc-manager ls -d           (显示详细信息)
     cc-manager ls --claude     (列出 CLAUDE 配置)`,
	Run: func(cmd *cobra.Command, args []string) {
		manager := config.NewManager()

		var configs []config.ConfigInfo
		var err error
		var configType string

		if claudeConfig {
			configs, err = manager.ListClaudes()
			configType = "CLAUDE"
		} else {
			configs, err = manager.ListConfigs()
			configType = "settings"
		}

		if err != nil {
			slog.Error("Failed to list configurations", "error", err, "type", configType)
			return
		}

		if len(configs) == 0 {
			fmt.Printf("未找到任何 %s 配置文件\n", configType)
			return
		}

		fmt.Printf("可用的 %s 配置:\n", configType)
		for _, cfg := range configs {
			if cfg.IsCurrent {
				if showDetail {
					fmt.Printf("  * %s (当前) - %s\n", cfg.Name, cfg.FullPath)
				} else {
					fmt.Printf("  * %s (当前)\n", cfg.Name)
				}
			} else {
				if showDetail {
					fmt.Printf("    %s - %s\n", cfg.Name, cfg.FullPath)
				} else {
					fmt.Printf("    %s\n", cfg.Name)
				}
			}
		}
	},
}

func init() {
	lsCmd.Flags().BoolVarP(&showDetail, "detail", "d", false, "显示详细信息（包括文件路径）")
	lsCmd.Flags().BoolVar(&claudeConfig, "claude", false, "列出 CLAUDE 配置而非 settings 配置")
	rootCmd.AddCommand(lsCmd)
}