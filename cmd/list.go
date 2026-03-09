package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/RogueCultivators/goon/internal/ui"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "列出项目中的所有模块",
	Long:  `显示 internal 目录下的所有业务模块`,
	Run: func(cmd *cobra.Command, args []string) {
		// 检查是否在项目根目录
		if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
			ui.Error("当前目录不是一个 Go 项目")
			return
		}

		// 检查 internal 目录是否存在
		if _, err := os.Stat("internal"); os.IsNotExist(err) {
			ui.Warning("未找到 internal 目录")
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

			// 检查是否是模块目录（包含 handler.go、service.go 或 routes.go）
			hasHandler := fileExists(filepath.Join(path, "handler.go"))
			hasService := fileExists(filepath.Join(path, "service.go"))
			hasRoutes := fileExists(filepath.Join(path, "routes.go"))

			if hasHandler || hasService || hasRoutes {
				moduleName := strings.TrimPrefix(path, "internal/")
				modules = append(modules, moduleName)
				return fs.SkipDir // 不再深入子目录
			}

			return nil
		})

		if err != nil {
			ui.Error(fmt.Sprintf("读取模块列表失败: %v", err))
			return
		}

		if len(modules) == 0 {
			ui.Warning("未找到任何模块")
			ui.Info("使用 'goon add <模块名>' 添加新模块")
			return
		}

		ui.Header("项目模块列表")
		for _, module := range modules {
			ui.Success(module)
			showModuleFiles(filepath.Join("internal", module))
		}
		ui.Info(fmt.Sprintf("共 %d 个模块", len(modules)))
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
		ui.Info(fmt.Sprintf("    (%s)", strings.Join(existingFiles, ", ")))
	}
}

func init() {
	rootCmd.AddCommand(listCmd)
}
