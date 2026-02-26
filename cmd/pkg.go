package cmd

import (
	"fmt"

	"github.com/RogueCultivators/goon/internal/generator"

	"github.com/spf13/cobra"
)

var addPkgCmd = &cobra.Command{
	Use:   "pkg [包名称]",
	Short: "添加功能包到项目",
	Long:  `添加额外的功能包，如 cache、jwt、email 等`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pkgName := args[0]

		fmt.Printf("正在添加功能包: %s\n", pkgName)

		if err := generator.AddPackage(pkgName); err != nil {
			fmt.Printf("添加功能包失败: %v\n", err)
			return
		}

		fmt.Println("\n✓ 功能包添加成功!")
	},
}

var listPkgCmd = &cobra.Command{
	Use:   "list-pkg",
	Short: "列出所有可用的功能包",
	Long:  `显示所有可以通过 'goon add pkg' 添加的功能包`,
	Run: func(cmd *cobra.Command, args []string) {
		pkgs := generator.ListAvailablePackages()

		fmt.Println("可用的功能包:")
		fmt.Println()
		for _, pkg := range pkgs {
			fmt.Printf("  - %s\n", pkg)
		}
		fmt.Println()
		fmt.Println("使用方法: goon add pkg <包名称>")
		fmt.Println("示例: goon add pkg cache")
	},
}

func init() {
	addCmd.AddCommand(addPkgCmd)
	rootCmd.AddCommand(listPkgCmd)
}
