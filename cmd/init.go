package cmd

import (
	"fmt"

	"github.com/RogueCultivators/goon/internal/generator"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [项目名称]",
	Short: "初始化一个新的 Gin Web 项目",
	Long:  `创建一个包含基础功能的模块化 Gin Web 项目结构`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		moduleName, _ := cmd.Flags().GetString("module")
		minimal, _ := cmd.Flags().GetBool("minimal")

		if moduleName == "" {
			moduleName = projectName
		}

		if minimal {
			fmt.Printf("正在初始化项目 (minimal 模式): %s\n", projectName)
		} else {
			fmt.Printf("正在初始化项目: %s\n", projectName)
		}

		if err := generator.InitProject(projectName, moduleName, minimal); err != nil {
			fmt.Printf("初始化失败: %v\n", err)
			return
		}

		fmt.Println("✓ 项目初始化成功!")
		if minimal {
			fmt.Println("\n提示: 使用 minimal 模式，只生成了核心文件")
			fmt.Println("如需添加额外功能包，请使用: goon add pkg <name>")
		}
		fmt.Printf("\n进入项目目录: cd %s\n", projectName)
		fmt.Println("运行项目: go run main.go")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringP("module", "m", "", "Go module 名称 (默认使用项目名称)")
	initCmd.Flags().Bool("minimal", false, "只生成核心文件和目录")
}
