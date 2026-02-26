package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove [模块名称]",
	Short: "删除项目中的模块",
	Long:  `删除指定的业务模块及其所有文件`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		moduleName := args[0]
		force, _ := cmd.Flags().GetBool("force")

		// 检查是否在项目根目录
		if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
			fmt.Println("错误: 当前目录不是一个 Go 项目")
			return
		}

		moduleDir := filepath.Join("internal", moduleName)

		// 检查模块是否存在
		if _, err := os.Stat(moduleDir); os.IsNotExist(err) {
			fmt.Printf("错误: 模块 %s 不存在\n", moduleName)
			return
		}

		// 如果没有 --force 标志，要求确认
		if !force {
			fmt.Printf("警告: 即将删除模块 %s 及其所有文件\n", moduleName)
			fmt.Print("确认删除? (y/N): ")

			var response string
			fmt.Scanln(&response)

			if response != "y" && response != "Y" {
				fmt.Println("已取消删除")
				return
			}
		}

		// 删除模块目录
		if err := os.RemoveAll(moduleDir); err != nil {
			fmt.Printf("删除模块失败: %v\n", err)
			return
		}

		fmt.Printf("✓ 模块 %s 已删除\n", moduleName)
		fmt.Println("\n注意: 请手动从 internal/router/router.go 中移除相关路由注册代码")
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
	removeCmd.Flags().BoolP("force", "f", false, "强制删除，不要求确认")
}
