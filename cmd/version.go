package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Version 版本号（通过 ldflags 注入）
	Version = "dev"
	// Commit Git 提交哈希（通过 ldflags 注入）
	Commit = "none"
	// BuildDate 构建日期（通过 ldflags 注入）
	BuildDate = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示版本信息",
	Long:  `显示 Goon 的版本信息，包括版本号、Git 提交哈希和构建日期。`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Goon %s\n", Version)
		fmt.Printf("Commit: %s\n", Commit)
		fmt.Printf("Built: %s\n", BuildDate)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
