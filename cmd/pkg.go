package cmd

import (
	"fmt"

	"github.com/RogueCultivators/goon/internal/generator"
	"github.com/RogueCultivators/goon/internal/ui"

	"github.com/spf13/cobra"
)

var addPkgCmd = &cobra.Command{
	Use:   "pkg [包名称]",
	Short: "添加功能包到项目",
	Long:  `添加额外的功能包，如 cache、jwt、email 等`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pkgName := args[0]

		ui.Step(fmt.Sprintf("正在添加功能包: %s", pkgName))

		if err := generator.AddPackage(pkgName); err != nil {
			ui.Error(fmt.Sprintf("添加功能包失败: %v", err))
			return
		}

		ui.Success("功能包添加成功!")
	},
}

var listPkgCmd = &cobra.Command{
	Use:   "list-pkg",
	Short: "列出所有可用的功能包",
	Long:  `显示所有可以通过 'goon add pkg' 添加的功能包`,
	Run: func(cmd *cobra.Command, args []string) {
		pkgs := generator.ListAvailablePackages()

		ui.Header("可用的功能包")
		for _, pkg := range pkgs {
			ui.Info(fmt.Sprintf("  - %s", pkg))
		}
		ui.Info("使用方法: goon add pkg <包名称>")
		ui.Info("示例: goon add pkg cache")
	},
}

func init() {
	addCmd.AddCommand(addPkgCmd)
	rootCmd.AddCommand(listPkgCmd)
}
