package cmd

import (
	"fmt"

	"github.com/RogueCultivators/goon/internal/generator"
	"github.com/RogueCultivators/goon/internal/ui"

	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "生成各种代码文件",
	Long:  `generate 命令用于生成各种类型的代码文件，如测试文件、文档等。`,
}

var generateTestCmd = &cobra.Command{
	Use:   "test [模块名称]",
	Short: "为模块生成测试文件",
	Long: `generate test 命令为指定模块生成测试文件。

示例:
  goon generate test user              # 为 user 模块生成所有测试文件
  goon generate test product -l handler # 只为 handler 层生成测试
  goon generate test order --all       # 为所有模块生成测试文件`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		all, _ := cmd.Flags().GetBool("all")
		layers, _ := cmd.Flags().GetStringSlice("layers")
		noColor, _ := cmd.Flags().GetBool("no-color")

		ui.NoColor = noColor

		if all {
			ui.Header("为所有模块生成测试文件")
			if err := generator.GenerateAllTests(); err != nil {
				ui.Error(fmt.Sprintf("生成测试失败: %v", err))
				return
			}
			ui.Success("所有模块的测试文件已生成")
			return
		}

		if len(args) == 0 {
			ui.Error("请提供模块名称或使用 --all 标志")
			return
		}

		moduleName := args[0]
		ui.Step(fmt.Sprintf("为模块 %s 生成测试文件", moduleName))

		if err := generator.GenerateModuleTests(moduleName, layers); err != nil {
			ui.Error(fmt.Sprintf("生成测试失败: %v", err))
			return
		}

		ui.Success(fmt.Sprintf("模块 %s 的测试文件已生成", moduleName))
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.AddCommand(generateTestCmd)

	generateTestCmd.Flags().Bool("all", false, "为所有模块生成测试")
	generateTestCmd.Flags().StringSliceP("layers", "l", []string{}, "指定要生成测试的层")
	generateTestCmd.Flags().Bool("no-color", false, "禁用彩色输出")
}
