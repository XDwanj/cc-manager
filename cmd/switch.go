package cmd

import (
	"cc-manager/internal/config"
	"log/slog"

	"github.com/spf13/cobra"
)

var switchClientName string
var switchTypeName string

var switchCmd = &cobra.Command{
	Use:   "switch <配置名>",
	Short: "切换配置",
	Long: `切换到指定客户端和类型的配置。

例如：cc-manager switch yescode                                (默认: --client=claude --type=config)
     cc-manager switch linus --client=claude --type=agents     (切换 claude agents 配置)
     cc-manager switch yescode --client=codex --type=config    (切换 codex config 配置)
     cc-manager switch custom --client=gemini --type=agents    (切换 gemini agents 配置)`,
	Args: cobra.ExactArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) >= 1 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		
		manager := config.NewManager()
		
		// 获取客户端和类型，使用默认值
		clientName := switchClientName
		if clientName == "" {
			clientName = config.GetDefaultClient().Name
		}
		
		typeName := switchTypeName
		if typeName == "" {
			typeName = config.GetDefaultConfigType().Name
		}
		
		configs, err := manager.List(clientName, typeName)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		
		var names []string
		for _, cfg := range configs {
			names = append(names, cfg.Name)
		}
		return names, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		configName := args[0]
		manager := config.NewManager()
		
		// 使用标志值或默认值
		clientName := switchClientName
		if clientName == "" {
			clientName = config.GetDefaultClient().Name
		}
		
		typeName := switchTypeName
		if typeName == "" {
			typeName = config.GetDefaultConfigType().Name
		}
		
		slog.Info("Switching configuration", "client", clientName, "type", typeName, "config", configName)
		err := manager.Switch(clientName, typeName, configName)
		
		if err != nil {
			slog.Error("Failed to switch configuration", "client", clientName, "type", typeName, "config", configName, "error", err)
			return
		}
		
		slog.Info("Successfully switched to configuration", "client", clientName, "type", typeName, "config", configName)
	},
}

func init() {
	switchCmd.Flags().StringVar(&switchClientName, "client", "", "指定客户端 (claude/codex/gemini，默认: claude)")
	switchCmd.Flags().StringVar(&switchTypeName, "type", "", "指定配置类型 (config/agents，默认: config)")
	
	// 为 --client 标志添加补全功能
	switchCmd.RegisterFlagCompletionFunc("client", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return config.GetClientNames(), cobra.ShellCompDirectiveNoFileComp
	})
	
	// 为 --type 标志添加补全功能
	switchCmd.RegisterFlagCompletionFunc("type", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return config.GetConfigTypeNames(), cobra.ShellCompDirectiveNoFileComp
	})
	
	rootCmd.AddCommand(switchCmd)
}