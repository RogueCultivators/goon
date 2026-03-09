package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/RogueCultivators/goon/internal/generator"
	"github.com/RogueCultivators/goon/internal/ui"
	"github.com/RogueCultivators/goon/internal/utils"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove [模块名称]",
	Short: "删除项目中的模块",
	Long:  `删除指定的业务模块及其所有文件`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		moduleName := utils.SanitizeInput(args[0])
		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			ui.Error(fmt.Sprintf("获取参数失败: %v", err))
			return
		}

		// 检查是否在项目根目录
		if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
			ui.Error("当前目录不是一个 Go 项目")
			return
		}

		moduleDir := filepath.Join("internal", moduleName)

		if err := utils.ValidatePath(".", moduleDir); err != nil {
			ui.Error(fmt.Sprintf("不安全的模块路径: %v", err))
			return
		}

		// 检查模块是否存在
		if _, err := os.Stat(moduleDir); os.IsNotExist(err) {
			ui.Error(fmt.Sprintf("模块 %s 不存在", moduleName))
			return
		}

		// 如果没有 --force 标志，要求确认
		if !force {
			ui.Warning(fmt.Sprintf("即将删除模块 %s 及其所有文件", moduleName))
			if !ui.Confirm("确认删除?") {
				ui.Info("已取消删除")
				return
			}
		}

		// 删除模块目录
		if err := os.RemoveAll(moduleDir); err != nil {
			ui.Error(fmt.Sprintf("删除模块失败: %v", err))
			return
		}

		ui.Success(fmt.Sprintf("模块 %s 已删除", moduleName))

		// 自动清理路由注册
		ui.Step("正在清理路由注册...")
		if err := generator.UnregisterModuleRoute(moduleName); err != nil {
			ui.Warning(fmt.Sprintf("路由清理失败: %v", err))
			ui.Info("请手动从 internal/router/router.go 中移除相关路由注册代码")
		} else {
			ui.Success("路由注册已自动清理")
		}
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
	removeCmd.Flags().BoolP("force", "f", false, "强制删除，不要求确认")
}
