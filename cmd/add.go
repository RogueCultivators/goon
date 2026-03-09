package cmd

import (
	"fmt"

	"github.com/RogueCultivators/goon/internal/generator"
	"github.com/RogueCultivators/goon/internal/ui"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [模块名称...]",
	Short: "添加一个或多个新模块到项目",
	Long: `add 命令用于向项目中添加新的模块，支持批量添加多个模块。

示例:
  goon add user                           # 生成所有层（handler, service, model, repository, schema, routes）
  goon add user product order             # 批量添加多个模块
  goon add product -l handler,service     # 只生成指定的层
  goon add order --layers handler         # 只生成 handler 层
  goon add post --no-register             # 不自动注册路由
  goon add user --example                 # 生成包含完整实现的示例代码

可用的层:
  - handler: HTTP 处理器
  - service: 业务逻辑层
  - model: 数据模型
  - repository: 数据访问层
  - schema: 请求/响应结构体
  - routes: 路由注册文件`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		moduleNames := args
		layers, ok := mustGetStringSlice(cmd, "layers")
		if !ok {
			return
		}
		register, ok := mustGetBool(cmd, "register")
		if !ok {
			return
		}
		verbose, ok := mustGetBool(cmd, "verbose")
		if !ok {
			return
		}
		example, ok := mustGetBool(cmd, "example")
		if !ok {
			return
		}
		dryRun, ok := mustGetBool(cmd, "dry-run")
		if !ok {
			return
		}

		if dryRun {
			ui.Info("预览模式 (--dry-run)")
			ui.Info("以下文件将会被生成:")
		}

		// 批量处理多个模块
		totalModules := len(moduleNames)
		successCount := 0
		failedModules := []string{}

		for i, moduleName := range moduleNames {
			if totalModules > 1 {
				ui.Header(fmt.Sprintf("处理模块 %d/%d: %s", i+1, totalModules, moduleName))
			}

			if verbose {
				ui.Step(fmt.Sprintf("正在添加模块: %s", moduleName))
				if len(layers) > 0 {
					ui.Info(fmt.Sprintf("指定的层: %v", layers))
				}
				if example {
					ui.Info("使用示例模板（包含完整实现）")
				}
			} else if !dryRun {
				ui.Step(fmt.Sprintf("正在添加模块: %s", moduleName))
			}

			if err := generator.AddModule(moduleName, layers, example, dryRun); err != nil {
				ui.Error(fmt.Sprintf("添加模块 %s 失败: %v", moduleName, err))
				failedModules = append(failedModules, moduleName)
				continue
			}

			if dryRun {
				ui.Info(fmt.Sprintf("模块 %s 的文件预览完成", moduleName))
				continue
			}

			ui.Success(fmt.Sprintf("模块 %s 添加成功!", moduleName))
			ui.Info(fmt.Sprintf("已生成模块目录: internal/%s/", moduleName))

			// 显示生成的文件
			if len(layers) == 0 {
				layers = []string{"handler", "service", "model", "repository", "schema", "routes"}
			}
			for _, layer := range layers {
				ui.Info(fmt.Sprintf("  - %s.go", layer))
			}

			// 自动注册路由
			if register {
				ui.Step(fmt.Sprintf("正在为 %s 注册路由...", moduleName))
				if err := generator.RegisterModuleRoute(moduleName); err != nil {
					ui.Warning(fmt.Sprintf("模块 %s 路由注册失败: %v", moduleName, err))
					ui.Info(fmt.Sprintf("请手动在 internal/router/router.go 中调用 %s.RegisterRoutes(api)", moduleName))
				} else {
					ui.Success(fmt.Sprintf("模块 %s 路由已自动注册到 internal/router/router.go", moduleName))
				}
			} else {
				ui.Info(fmt.Sprintf("请在 internal/router/router.go 中调用 %s.RegisterRoutes(api)", moduleName))
			}

			successCount++
		}

		if dryRun {
			ui.Info("提示: 移除 --dry-run 标志以实际生成这些文件")
			return
		}

		// 批量操作总结
		if totalModules > 1 {
			ui.Header("批量操作完成")
			ui.Success(fmt.Sprintf("成功添加 %d/%d 个模块", successCount, totalModules))
			if len(failedModules) > 0 {
				ui.Warning(fmt.Sprintf("失败的模块: %v", failedModules))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringSliceP("layers", "l", []string{}, "指定要生成的层 (handler,service,model,repository,schema,routes)")
	addCmd.Flags().Bool("register", true, "自动在 router.go 中注册路由")
	addCmd.Flags().BoolP("verbose", "v", false, "显示详细日志")
	addCmd.Flags().Bool("example", false, "生成包含完整实现的示例代码")
	addCmd.Flags().Bool("dry-run", false, "预览将要生成的文件，不实际创建")
}
