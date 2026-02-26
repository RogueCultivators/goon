package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/RogueCultivators/goon/internal/ui"
	"github.com/RogueCultivators/goon/internal/utils"
)

// RenameModule 重命名模块
func RenameModule(oldName, newName string, dryRun bool) error {
	// 标准化模块名称
	oldName = utils.ToSnakeCase(oldName)
	newName = utils.ToSnakeCase(newName)

	oldDir := filepath.Join("internal", oldName)
	newDir := filepath.Join("internal", newName)

	// 检查旧模块是否存在
	if _, err := os.Stat(oldDir); os.IsNotExist(err) {
		return fmt.Errorf("模块 %s 不存在", oldName)
	}

	// 检查新模块名是否已存在
	if _, err := os.Stat(newDir); err == nil {
		return fmt.Errorf("模块 %s 已存在", newName)
	}

	// 获取项目模块名
	projectModule, err := getModuleNameFromGoMod()
	if err != nil {
		return err
	}

	oldCapitalized := utils.ToPascalCase(oldName)
	newCapitalized := utils.ToPascalCase(newName)

	if dryRun {
		ui.Step("将要执行的操作:")
		ui.Info(fmt.Sprintf("1. 重命名目录: %s → %s", oldDir, newDir))
		ui.Info("2. 更新模块内所有文件的包名")
		ui.Info("3. 更新 router.go 中的导入和路由注册")
		ui.Info("4. 更新所有引用此模块的文件")
		return nil
	}

	// 创建备份管理器
	bm := utils.NewBackupManager()

	// 1. 重命名目录
	ui.Step(fmt.Sprintf("重命名目录: %s → %s", oldDir, newDir))
	if err := os.Rename(oldDir, newDir); err != nil {
		return fmt.Errorf("重命名目录失败: %w", err)
	}

	// 2. 更新模块内所有文件
	ui.Step("更新模块内文件...")
	files, err := filepath.Glob(filepath.Join(newDir, "*.go"))
	if err != nil {
		return fmt.Errorf("查找文件失败: %w", err)
	}

	for _, file := range files {
		if err := bm.BackupFile(file); err != nil {
			bm.Rollback()
			return err
		}

		content, err := os.ReadFile(file)
		if err != nil {
			bm.Rollback()
			return fmt.Errorf("读取文件失败 %s: %w", file, err)
		}

		newContent := string(content)
		// 更新包名
		newContent = strings.ReplaceAll(newContent, fmt.Sprintf("package %s", oldName), fmt.Sprintf("package %s", newName))
		// 更新类型名称
		newContent = strings.ReplaceAll(newContent, oldCapitalized, newCapitalized)
		// 更新变量名
		newContent = strings.ReplaceAll(newContent, oldName, newName)

		if err := os.WriteFile(file, []byte(newContent), 0644); err != nil {
			bm.Rollback()
			return fmt.Errorf("写入文件失败 %s: %w", file, err)
		}
	}

	// 3. 更新 router.go
	ui.Step("更新 router.go...")
	routerPath := filepath.Join("internal", "router", "router.go")
	if _, err := os.Stat(routerPath); err == nil {
		if err := bm.BackupFile(routerPath); err != nil {
			bm.Rollback()
			return err
		}

		content, err := os.ReadFile(routerPath)
		if err != nil {
			bm.Rollback()
			return fmt.Errorf("读取 router.go 失败: %w", err)
		}

		newContent := string(content)
		// 更新导入
		oldImport := fmt.Sprintf("\"%s/internal/%s\"", projectModule, oldName)
		newImport := fmt.Sprintf("\"%s/internal/%s\"", projectModule, newName)
		newContent = strings.ReplaceAll(newContent, oldImport, newImport)
		// 更新路由注册
		newContent = strings.ReplaceAll(newContent, fmt.Sprintf("%s.RegisterRoutes", oldName), fmt.Sprintf("%s.RegisterRoutes", newName))

		if err := os.WriteFile(routerPath, []byte(newContent), 0644); err != nil {
			bm.Rollback()
			return fmt.Errorf("写入 router.go 失败: %w", err)
		}
	}

	// 4. 查找并更新所有引用此模块的文件
	ui.Step("搜索并更新其他引用...")
	if err := updateReferences(projectModule, oldName, newName, bm); err != nil {
		bm.Rollback()
		return err
	}

	ui.Success("所有文件已更新")
	return nil
}

func updateReferences(projectModule, oldName, newName string, bm *utils.BackupManager) error {
	// 搜索所有 Go 文件
	var goFiles []string
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			// 跳过 vendor 和隐藏目录
			if !strings.Contains(path, "vendor") && !strings.HasPrefix(path, ".") {
				goFiles = append(goFiles, path)
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("遍历文件失败: %w", err)
	}

	oldImport := fmt.Sprintf("%s/internal/%s", projectModule, oldName)
	newImport := fmt.Sprintf("%s/internal/%s", projectModule, newName)

	for _, file := range goFiles {
		content, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		fileContent := string(content)
		if strings.Contains(fileContent, oldImport) {
			if err := bm.BackupFile(file); err != nil {
				return err
			}

			newContent := strings.ReplaceAll(fileContent, oldImport, newImport)
			if err := os.WriteFile(file, []byte(newContent), 0644); err != nil {
				return fmt.Errorf("更新文件失败 %s: %w", file, err)
			}
		}
	}

	return nil
}
