package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/RogueCultivators/goon/internal/ui"
	"github.com/RogueCultivators/goon/internal/utils"
)

// RenameModule 重命名模块
func RenameModule(oldName, newName string, dryRun bool) error {
	oldName = utils.SanitizeInput(oldName)
	newName = utils.SanitizeInput(newName)
	if oldName == "" || newName == "" {
		return fmt.Errorf("模块名称不能为空")
	}

	oldName = utils.ToSnakeCase(oldName)
	newName = utils.ToSnakeCase(newName)

	oldDir := filepath.Join("internal", oldName)
	newDir := filepath.Join("internal", newName)

	if err := utils.ValidatePath(".", newDir); err != nil {
		return fmt.Errorf("不安全的模块路径: %w", err)
	}

	if err := validateRenameOperation(oldDir, newDir, oldName, newName); err != nil {
		return err
	}

	projectModule, err := getModuleNameFromGoMod()
	if err != nil {
		return err
	}

	if dryRun {
		return showRenameDryRun(oldDir, newDir)
	}

	return executeRename(oldDir, newDir, oldName, newName, projectModule)
}

func validateRenameOperation(oldDir, newDir, oldName, newName string) error {
	if _, err := os.Stat(oldDir); os.IsNotExist(err) {
		return fmt.Errorf("模块 %s 不存在", oldName)
	}

	if _, err := os.Stat(newDir); err == nil {
		return fmt.Errorf("模块 %s 已存在", newName)
	}

	return nil
}

func showRenameDryRun(oldDir, newDir string) error {
	ui.Step("将要执行的操作:")
	ui.Info(fmt.Sprintf("1. 重命名目录: %s → %s", oldDir, newDir))
	ui.Info("2. 更新模块内所有文件的包名")
	ui.Info("3. 更新 router.go 中的导入和路由注册")
	ui.Info("4. 更新所有引用此模块的文件")
	return nil
}

func executeRename(oldDir, newDir, oldName, newName, projectModule string) error {
	bm := utils.NewBackupManager()

	ui.Step(fmt.Sprintf("重命名目录: %s → %s", oldDir, newDir))
	if renameErr := os.Rename(oldDir, newDir); renameErr != nil {
		return fmt.Errorf("重命名目录失败: %w", renameErr)
	}

	if err := updateModuleFiles(newDir, oldName, newName, bm); err != nil {
		return err
	}

	if err := updateRouterFile(oldName, newName, projectModule, bm); err != nil {
		return err
	}

	if err := updateReferences(projectModule, oldName, newName, bm); err != nil {
		return err
	}

	ui.Success("所有文件已更新")
	return nil
}

func updateModuleFiles(newDir, oldName, newName string, bm *utils.BackupManager) error {
	ui.Step("更新模块内文件...")
	files, err := filepath.Glob(filepath.Join(newDir, "*.go"))
	if err != nil {
		return fmt.Errorf("查找文件失败: %w", err)
	}

	oldCapitalized := utils.ToPascalCase(oldName)
	newCapitalized := utils.ToPascalCase(newName)

	for _, file := range files {
		if err := processModuleFile(file, oldName, newName, oldCapitalized, newCapitalized, bm); err != nil {
			return err
		}
	}
	return nil
}

func processModuleFile(file, oldName, newName, oldCapitalized, newCapitalized string, bm *utils.BackupManager) error {
	if err := bm.BackupFile(file); err != nil {
		if rollbackErr := bm.Rollback(); rollbackErr != nil {
			return fmt.Errorf("备份失败且回滚失败: %w, 回滚错误: %v", err, rollbackErr)
		}
		return err
	}

	content, err := os.ReadFile(file)
	if err != nil {
		if rollbackErr := bm.Rollback(); rollbackErr != nil {
			return fmt.Errorf("读取文件失败 %s: %w, 回滚错误: %v", file, err, rollbackErr)
		}
		return fmt.Errorf("读取文件失败 %s: %w", file, err)
	}

	newContent := string(content)
	// 1. 替换 package 声明（精确匹配）
	newContent = strings.ReplaceAll(newContent, fmt.Sprintf("package %s", oldName), fmt.Sprintf("package %s", newName))
	// 2. 使用单词边界替换 PascalCase 标识符
	reCapitalized := regexp.MustCompile(`\b` + regexp.QuoteMeta(oldCapitalized) + `\b`)
	newContent = reCapitalized.ReplaceAllString(newContent, newCapitalized)
	// 3. 使用单词边界替换 snake_case 标识符（避免替换子串）
	reName := regexp.MustCompile(`\b` + regexp.QuoteMeta(oldName) + `\b`)
	newContent = reName.ReplaceAllString(newContent, newName)

	if err := os.WriteFile(file, []byte(newContent), 0o600); err != nil {
		if rollbackErr := bm.Rollback(); rollbackErr != nil {
			return fmt.Errorf("写入文件失败 %s: %w, 回滚错误: %v", file, err, rollbackErr)
		}
		return fmt.Errorf("写入文件失败 %s: %w", file, err)
	}
	return nil
}

func updateRouterFile(oldName, newName, projectModule string, bm *utils.BackupManager) error {
	ui.Step("更新 router.go...")
	routerPath := filepath.Join("internal", "router", "router.go")

	if _, err := os.Stat(routerPath); os.IsNotExist(err) {
		return nil
	}

	if err := bm.BackupFile(routerPath); err != nil {
		if rollbackErr := bm.Rollback(); rollbackErr != nil {
			return fmt.Errorf("备份 router.go 失败且回滚失败: %w, 回滚错误: %v", err, rollbackErr)
		}
		return err
	}

	content, err := os.ReadFile(routerPath)
	if err != nil {
		if rollbackErr := bm.Rollback(); rollbackErr != nil {
			return fmt.Errorf("读取 router.go 失败: %w, 回滚错误: %v", err, rollbackErr)
		}
		return fmt.Errorf("读取 router.go 失败: %w", err)
	}

	newContent := updateRouterContent(string(content), oldName, newName, projectModule)

	if err := os.WriteFile(routerPath, []byte(newContent), 0o600); err != nil {
		if rollbackErr := bm.Rollback(); rollbackErr != nil {
			return fmt.Errorf("写入 router.go 失败: %w, 回滚错误: %v", err, rollbackErr)
		}
		return fmt.Errorf("写入 router.go 失败: %w", err)
	}
	return nil
}

func updateRouterContent(content, oldName, newName, projectModule string) string {
	oldImport := fmt.Sprintf("\"%s/internal/%s\"", projectModule, oldName)
	newImport := fmt.Sprintf("\"%s/internal/%s\"", projectModule, newName)
	content = strings.ReplaceAll(content, oldImport, newImport)
	content = strings.ReplaceAll(content, fmt.Sprintf("%s.RegisterRoutes", oldName), fmt.Sprintf("%s.RegisterRoutes", newName))
	return content
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
			if err := os.WriteFile(file, []byte(newContent), 0o600); err != nil {
				return fmt.Errorf("更新文件失败 %s: %w", file, err)
			}
		}
	}

	return nil
}
