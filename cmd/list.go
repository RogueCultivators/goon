package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "列出项目中的所有模块",
	Long:  `显示 internal 目录下的所有业务模块`,
	Run: func(cmd *cobra.Command, args []string) {
		// 检查是否在项目根目录
		if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
			fmt.Println("错误: 当前目录不是一个 Go 项目")
			return
		}

		// 检查 internal 目录是否存在
		if _, err := os.Stat("internal"); os.IsNotExist(err) {
			fmt.Println("未找到 internal 目录")
			return
		}

		modules := []string{}

		// 遍历 internal 目录
		err := filepath.WalkDir("internal", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			// 跳过非目录和系统目录
			if !d.IsDir() {
				return nil
			}

			// 跳过 internal 本身和系统目录
			if path == "internal" || strings.Contains(path, "config") ||
				strings.Contains(path, "middleware") || strings.Contains(path, "router") ||
				strings.Contains(path, "sqlc") {
				return nil
			}

			// 检查是否是模块目录（包含 handler.go 或 service.go）
			hasHandler := fileExists(filepath.Join(path, "handler.go"))
			hasService := fileExists(filepath.Join(path, "service.go"))

			if hasHandler || hasService {
				moduleName := strings.TrimPrefix(path, "internal/")
				modules = append(modules, moduleName)
				return fs.SkipDir // 不再深入子目录
			}

			return nil
		})

		if err != nil {
			fmt.Printf("读取模块列表失败: %v\n", err)
			return
		}

		if len(modules) == 0 {
			fmt.Println("未找到任何模块")
			fmt.Println("\n使用 'goon add <模块名>' 添加新模块")
			return
		}

		fmt.Println("项目模块列表:")
		fmt.Println()
		for _, module := range modules {
			fmt.Printf("  - %s\n", module)
			showModuleFiles(filepath.Join("internal", module))
		}
		fmt.Println()
		fmt.Printf("共 %d 个模块\n", len(modules))
	},
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func showModuleFiles(modulePath string) {
	files := []string{"handler.go", "service.go", "model.go", "repository.go", "schema.go"}
	existingFiles := []string{}

	for _, file := range files {
		if fileExists(filepath.Join(modulePath, file)) {
			existingFiles = append(existingFiles, file)
		}
	}

	if len(existingFiles) > 0 {
		fmt.Printf("    (%s)\n", strings.Join(existingFiles, ", "))
	}
}

func init() {
	rootCmd.AddCommand(listCmd)
}
