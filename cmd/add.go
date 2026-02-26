package cmd

import (
	"fmt"

	"github.com/RogueCultivators/goon/internal/generator"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [模块名称]",
	Short: "添加一个新模块到项目",
	Long: `add 命令用于向项目中添加新的模块。

示例:
  goon add user                           # 生成所有层（handler, service, model, repository, schema, routes）
  goon add product -l handler,service     # 只生成指定的层
  goon add order --layers handler         # 只生成 handler 层
  goon add post --no-register             # 不自动注册路由

可用的层:
  - handler: HTTP 处理器
  - service: 业务逻辑层
  - model: 数据模型
  - repository: 数据访问层
  - schema: 请求/响应结构体
  - routes: 路由注册文件`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		moduleName := args[0]
		layers, _ := cmd.Flags().GetStringSlice("layers")
		register, _ := cmd.Flags().GetBool("register")
		verbose, _ := cmd.Flags().GetBool("verbose")

		if verbose {
			fmt.Printf("正在添加模块: %s\n", moduleName)
			if len(layers) > 0 {
				fmt.Printf("指定的层: %v\n", layers)
			}
		} else {
			fmt.Printf("正在添加模块: %s\n", moduleName)
		}

		if err := generator.AddModule(moduleName, layers); err != nil {
			fmt.Printf("添加模块失败: %v\n", err)
			return
		}

		fmt.Println("✓ 模块添加成功!")
		fmt.Printf("\n已生成模块目录: internal/%s/\n", moduleName)

		// 显示生成的文件
		if len(layers) == 0 {
			layers = []string{"handler", "service", "model", "repository", "schema", "routes"}
		}
		for _, layer := range layers {
			fmt.Printf("  - %s.go\n", layer)
		}

		// 自动注册路由
		if register {
			fmt.Println("\n正在注册路由...")
			if err := generator.RegisterModuleRoute(moduleName); err != nil {
				fmt.Printf("⚠ 路由注册失败: %v\n", err)
				fmt.Printf("请手动在 internal/router/router.go 中调用 %s.RegisterRoutes(api)\n", moduleName)
			} else {
				fmt.Println("✓ 路由已自动注册到 internal/router/router.go")
			}
		} else {
			fmt.Printf("\n请在 internal/router/router.go 中调用 %s.RegisterRoutes(api)\n", moduleName)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringSliceP("layers", "l", []string{}, "指定要生成的层 (handler,service,model,repository,schema,routes)")
	addCmd.Flags().Bool("register", true, "自动在 router.go 中注册路由")
	addCmd.Flags().BoolP("verbose", "v", false, "显示详细日志")
}
