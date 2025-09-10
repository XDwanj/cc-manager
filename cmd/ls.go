package cmd

import (
	"cc-manager/internal/config"
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
)

var showDetail bool
var clientName string
var typeName string

var lsCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "列出所有可用的配置",
	Long: `列出指定客户端和类型的所有可用配置文件。
当前激活的配置会用 * 标记。

例如：cc-manager ls                                (默认: --client=claude --type=config)
     cc-manager ls --client=claude --type=agents   (列出 claude agents 配置)
     cc-manager ls --client=codex --type=config    (列出 codex config 配置)
     cc-manager ls --client=gemini --type=agents   (列出 gemini agents 配置)
     cc-manager ls -d                              (显示详细信息)`,
	Run: func(cmd *cobra.Command, args []string) {
		manager := config.NewManager()

		// 使用标志值或默认值
		if clientName == "" {
			clientName = config.GetDefaultClient().Name
		}
		if typeName == "" {
			typeName = config.GetDefaultConfigType().Name
		}

		configs, err := manager.List(clientName, typeName)
		if err != nil {
			slog.Error("Failed to list configurations", "error", err, "client", clientName, "type", typeName)
			return
		}

		if len(configs) == 0 {
			fmt.Printf("未找到任何 %s %s 配置文件\n", clientName, typeName)
			return
		}

		fmt.Printf("可用的 %s %s 配置:\n", clientName, typeName)
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
	lsCmd.Flags().StringVar(&clientName, "client", "", "指定客户端 (claude/codex/gemini，默认: claude)")
	lsCmd.Flags().StringVar(&typeName, "type", "", "指定配置类型 (config/agents，默认: config)")
	
	// 为 --client 标志添加补全功能
	lsCmd.RegisterFlagCompletionFunc("client", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return config.GetClientNames(), cobra.ShellCompDirectiveNoFileComp
	})
	
	// 为 --type 标志添加补全功能
	lsCmd.RegisterFlagCompletionFunc("type", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return config.GetConfigTypeNames(), cobra.ShellCompDirectiveNoFileComp
	})
	
	rootCmd.AddCommand(lsCmd)
}