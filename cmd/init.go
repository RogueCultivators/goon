package cmd

import (
	"fmt"

	"github.com/RogueCultivators/goon/internal/generator"
	"github.com/RogueCultivators/goon/internal/interactive"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [项目名称]",
	Short: "初始化一个新的 Gin Web 项目",
	Long:  `创建一个包含基础功能的模块化 Gin Web 项目结构`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		interactiveMode, err := cmd.Flags().GetBool("interactive")
		if err != nil {
			fmt.Printf("获取参数失败: %v\n", err)
			return
		}

		var projectName, moduleName string
		var minimal, example bool

		// 交互式模式
		if interactiveMode {
			fmt.Println("🎯 交互式项目初始化向导")

			config, err := interactive.RunInitWizard()
			if err != nil {
				fmt.Printf("向导执行失败: %v\n", err)
				return
			}

			interactive.PrintSummary(config)

			projectName = config.ProjectName
			moduleName = config.ModuleName
			example = config.UseExample
			// minimal 模式在交互式中不使用，因为我们有更细粒度的控制
		} else {
			// 命令行模式
			if len(args) == 0 {
				fmt.Println("错误: 请提供项目名称，或使用 --interactive 进入交互式模式")
				fmt.Println("用法: goon init <项目名称> [选项]")
				fmt.Println("      goon init --interactive")
				return
			}

			projectName = args[0]
			moduleName, err = cmd.Flags().GetString("module")
			if err != nil {
				fmt.Printf("获取参数失败: %v\n", err)
				return
			}
			minimal, err = cmd.Flags().GetBool("minimal")
			if err != nil {
				fmt.Printf("获取参数失败: %v\n", err)
				return
			}
			example, err = cmd.Flags().GetBool("example")
			if err != nil {
				fmt.Printf("获取参数失败: %v\n", err)
				return
			}

			if moduleName == "" {
				moduleName = projectName
			}
		}

		// 显示初始化信息
		if minimal {
			fmt.Printf("正在初始化项目 (minimal 模式): %s\n", projectName)
		} else if example {
			fmt.Printf("正在初始化项目 (包含示例代码): %s\n", projectName)
		} else {
			fmt.Printf("正在初始化项目: %s\n", projectName)
		}

		// 执行项目初始化
		if err := generator.InitProject(projectName, moduleName, minimal, example); err != nil {
			fmt.Printf("初始化失败: %v\n", err)
			return
		}

		fmt.Println("✓ 项目初始化成功!")
		if minimal {
			fmt.Println("\n提示: 使用 minimal 模式，只生成了核心文件")
			fmt.Println("如需添加额外功能包，请使用: goon add pkg <name>")
		} else if example {
			fmt.Println("\n提示: 已生成包含完整实现的示例代码")
			fmt.Println("示例模块位于: internal/user/")
			fmt.Println("\n快速开始:")
			fmt.Println("1. 启动数据库: docker-compose up -d")
			fmt.Println("2. 运行迁移: make migrate-up")
			fmt.Println("3. 启动服务: go run main.go")
		}
		fmt.Printf("\n进入项目目录: cd %s\n", projectName)
		if !example {
			fmt.Println("运行项目: go run main.go")
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringP("module", "m", "", "Go module 名称 (默认使用项目名称)")
	initCmd.Flags().Bool("minimal", false, "只生成核心文件和目录")
	initCmd.Flags().Bool("example", false, "生成包含完整实现的示例代码（user 模块）")
	initCmd.Flags().BoolP("interactive", "i", false, "交互式向导模式")
}
