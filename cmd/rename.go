package cmd

import (
	"fmt"

	"github.com/RogueCultivators/goon/internal/generator"
	"github.com/RogueCultivators/goon/internal/ui"

	"github.com/spf13/cobra"
)

var renameCmd = &cobra.Command{
	Use:   "rename [旧模块名] [新模块名]",
	Short: "重命名一个模块",
	Long: `rename 命令用于重命名项目中的模块，并自动更新所有引用。

示例:
  goon rename user account           # 将 user 模块重命名为 account
  goon rename user-profile profile   # 支持各种命名格式
  goon rename old new --dry-run      # 预览将要进行的更改`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		oldName := args[0]
		newName := args[1]
		dryRun, err := cmd.Flags().GetBool("dry-run")
		if err != nil {
			ui.Error(fmt.Sprintf("获取参数失败: %v", err))
			return
		}
		noColor, err := cmd.Flags().GetBool("no-color")
		if err != nil {
			ui.Error(fmt.Sprintf("获取参数失败: %v", err))
			return
		}

		ui.NoColor = noColor

		if dryRun {
			ui.Info("Dry-run 模式：将显示要进行的更改但不实际执行")
		}

		ui.Header(fmt.Sprintf("重命名模块: %s → %s", oldName, newName))

		if err := generator.RenameModule(oldName, newName, dryRun); err != nil {
			ui.Error(fmt.Sprintf("重命名失败: %v", err))
			return
		}

		if dryRun {
			ui.Success("Dry-run 完成！以上是将要进行的更改")
			ui.Info("运行不带 --dry-run 标志的命令来实际执行重命名")
		} else {
			ui.Success(fmt.Sprintf("模块已成功从 %s 重命名为 %s", oldName, newName))
		}
	},
}

func init() {
	rootCmd.AddCommand(renameCmd)
	renameCmd.Flags().Bool("dry-run", false, "预览更改但不实际执行")
	renameCmd.Flags().Bool("no-color", false, "禁用彩色输出")
}
