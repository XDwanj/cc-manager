package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "生成自动补全脚本",
	Long: `生成指定 shell 的自动补全脚本。

支持的 shell:
  bash        为 bash 生成补全脚本
  zsh         为 zsh 生成补全脚本  
  fish        为 fish 生成补全脚本
  powershell  为 PowerShell 生成补全脚本

使用方法:

Bash:
  $ cc-manager completion bash > /etc/bash_completion.d/cc-manager
  $ source ~/.bashrc

Zsh:
  # 如果启用了 zsh 补全，将脚本写入补全目录之一
  $ cc-manager completion zsh > "${fpath[1]}/_cc-manager"

Fish:
  $ cc-manager completion fish > ~/.config/fish/completions/cc-manager.fish

PowerShell:
  PS> cc-manager completion powershell | Out-String | Invoke-Expression
  # 将上述命令添加到 PowerShell 配置文件中使其永久生效`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}