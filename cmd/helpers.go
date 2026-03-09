package cmd

import (
	"fmt"

	"github.com/RogueCultivators/goon/internal/ui"

	"github.com/spf13/cobra"
)

// mustGetBool 获取 bool 类型的 flag，失败时 panic（cobra 注册过的 flag 不应失败）
func mustGetBool(cmd *cobra.Command, name string) (value, ok bool) {
	val, err := cmd.Flags().GetBool(name)
	if err != nil {
		ui.Error(fmt.Sprintf("获取参数 %s 失败: %v", name, err))
		return false, false
	}
	return val, true
}

// mustGetString 获取 string 类型的 flag
func mustGetString(cmd *cobra.Command, name string) (string, bool) {
	val, err := cmd.Flags().GetString(name)
	if err != nil {
		ui.Error(fmt.Sprintf("获取参数 %s 失败: %v", name, err))
		return "", false
	}
	return val, true
}

// mustGetStringSlice 获取 string slice 类型的 flag
func mustGetStringSlice(cmd *cobra.Command, name string) ([]string, bool) {
	val, err := cmd.Flags().GetStringSlice(name)
	if err != nil {
		ui.Error(fmt.Sprintf("获取参数 %s 失败: %v", name, err))
		return nil, false
	}
	return val, true
}
