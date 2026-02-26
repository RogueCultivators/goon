package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "goon",
	Short: "Goon - 快速生成 Gin Web 项目的 CLI 工具",
	Long:  `Goon 是一个基于 Cobra 和 Gin 的命令行工具，用于快速初始化和管理模块化的 Web 项目。`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
