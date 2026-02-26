package generator

import (
	"fmt"

	"os"

	"golang.org/x/mod/modfile"
)

// getModuleNameFromGoMod 使用 golang.org/x/mod/modfile 解析 go.mod
func getModuleNameFromGoMod() (string, error) {
	data, err := os.ReadFile("go.mod")
	if err != nil {
		return "", fmt.Errorf("读取 go.mod 失败: %w", err)
	}

	f, err := modfile.Parse("go.mod", data, nil)
	if err != nil {
		return "", fmt.Errorf("解析 go.mod 失败: %w", err)
	}

	if f.Module == nil {
		return "", fmt.Errorf("go.mod 中没有 module 声明")
	}

	return f.Module.Mod.Path, nil
}
