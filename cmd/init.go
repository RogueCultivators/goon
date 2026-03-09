package cmd

import (
	"fmt"

	"github.com/RogueCultivators/goon/internal/generator"
	"github.com/RogueCultivators/goon/internal/interactive"
	"github.com/RogueCultivators/goon/internal/ui"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [项目名称]",
	Short: "初始化一个新的 Gin Web 项目",
	Long:  `创建一个包含基础功能的模块化 Gin Web 项目结构`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		interactiveMode, ok := mustGetBool(cmd, "interactive")
		if !ok {
			return
		}

		var opts generator.InitOptions

		// 交互式模式
		if interactiveMode {
			ui.Header("交互式项目初始化向导")

			var config *interactive.ProjectConfig
			var err error
			config, err = interactive.RunInitWizard()
			if err != nil {
				ui.Error(fmt.Sprintf("向导执行失败: %v", err))
				return
			}

			interactive.PrintSummary(config)

			opts = generator.InitOptions{
				ProjectName:   config.ProjectName,
				ModuleName:    config.ModuleName,
				Example:       config.UseExample,
				Database:      config.Database,
				UseAuth:       config.UseAuth,
				AuthMethod:    config.AuthMethod,
				UseDocker:     config.UseDocker,
				ExampleModule: config.ExampleModule,
			}
		} else {
			// 命令行模式
			if len(args) == 0 {
				ui.Error("请提供项目名称，或使用 --interactive 进入交互式模式")
				ui.Info("用法: goon init <项目名称> [选项]")
				ui.Info("      goon init --interactive")
				return
			}

			projectName := args[0]
			moduleName, ok := mustGetString(cmd, "module")
			if !ok {
				return
			}
			minimal, ok := mustGetBool(cmd, "minimal")
			if !ok {
				return
			}
			example, ok := mustGetBool(cmd, "example")
			if !ok {
				return
			}

			if moduleName == "" {
				moduleName = projectName
			}

			opts = generator.InitOptions{
				ProjectName: projectName,
				ModuleName:  moduleName,
				Minimal:     minimal,
				Example:     example,
				UseDocker:   true,
				UseAuth:     true,
				AuthMethod:  "JWT",
				Database:    "PostgreSQL",
			}
		}

		// 显示初始化信息
		if opts.Minimal {
			ui.Step(fmt.Sprintf("正在初始化项目 (minimal 模式): %s", opts.ProjectName))
		} else if opts.Example {
			ui.Step(fmt.Sprintf("正在初始化项目 (包含示例代码): %s", opts.ProjectName))
		} else {
			ui.Step(fmt.Sprintf("正在初始化项目: %s", opts.ProjectName))
		}

		// 执行项目初始化
		if err := generator.InitProject(opts); err != nil {
			ui.Error(fmt.Sprintf("初始化失败: %v", err))
			return
		}

		ui.Success("项目初始化成功!")
		if opts.Minimal {
			ui.Info("使用 minimal 模式，只生成了核心文件")
			ui.Info("如需添加额外功能包，请使用: goon add pkg <name>")
		} else if opts.Example {
			exampleModule := opts.ExampleModule
			if exampleModule == "" {
				exampleModule = "user"
			}
			ui.Info("已生成包含完整实现的示例代码")
			ui.Info(fmt.Sprintf("示例模块位于: internal/%s/", exampleModule))
			ui.Info("快速开始:")
			ui.Step("1. 启动数据库: docker-compose up -d")
			ui.Step("2. 运行迁移: make migrate-up")
			ui.Step("3. 启动服务: go run main.go")
		}
		ui.Info(fmt.Sprintf("进入项目目录: cd %s", opts.ProjectName))
		if !opts.Example {
			ui.Info("运行项目: go run main.go")
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
