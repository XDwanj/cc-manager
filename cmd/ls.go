package cmd

import (
	"cc-manager/internal/config"
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
)

var showDetail bool

var lsCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "列出所有可用的 Claude 配置",
	Long: `列出 ~/.claude 目录下所有可用的配置文件。
当前激活的配置会用 * 标记。

例如：cc-manager ls
     cc-manager ls -d  (显示详细信息)`,
	Run: func(cmd *cobra.Command, args []string) {
		manager := config.NewManager()

		configs, err := manager.ListConfigs()
		if err != nil {
			slog.Error("Failed to list configurations", "error", err)
			return
		}

		if len(configs) == 0 {
			fmt.Println("未找到任何配置文件")
			return
		}

		fmt.Println("可用的 Claude 配置:")
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
	rootCmd.AddCommand(lsCmd)
}