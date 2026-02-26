package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// FileBackup 文件备份结构
type FileBackup struct {
	Path    string
	Content []byte
	IsNew   bool // 是否是新创建的文件
}

// BackupManager 备份管理器
type BackupManager struct {
	backups []FileBackup
}

// NewBackupManager 创建新的备份管理器
func NewBackupManager() *BackupManager {
	return &BackupManager{
		backups: make([]FileBackup, 0),
	}
}

// BackupFile 备份文件
func (bm *BackupManager) BackupFile(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// 文件不存在，标记为新文件
			bm.backups = append(bm.backups, FileBackup{
				Path:  path,
				IsNew: true,
			})
			return nil
		}
		return fmt.Errorf("备份文件失败 %s: %w", path, err)
	}

	bm.backups = append(bm.backups, FileBackup{
		Path:    path,
		Content: content,
		IsNew:   false,
	})

	return nil
}

// Rollback 回滚所有更改
func (bm *BackupManager) Rollback() error {
	var errors []error

	// 反向遍历备份，恢复文件
	for i := len(bm.backups) - 1; i >= 0; i-- {
		backup := bm.backups[i]

		if backup.IsNew {
			// 删除新创建的文件
			if err := os.Remove(backup.Path); err != nil && !os.IsNotExist(err) {
				errors = append(errors, fmt.Errorf("删除文件失败 %s: %w", backup.Path, err))
			}
		} else {
			// 恢复原文件内容
			if err := os.WriteFile(backup.Path, backup.Content, 0644); err != nil {
				errors = append(errors, fmt.Errorf("恢复文件失败 %s: %w", backup.Path, err))
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("回滚过程中发生 %d 个错误: %v", len(errors), errors)
	}

	return nil
}

// Clear 清空备份
func (bm *BackupManager) Clear() {
	bm.backups = make([]FileBackup, 0)
}

// ValidatePath 验证路径安全性，防止路径遍历攻击
func ValidatePath(basePath, targetPath string) error {
	// 获取绝对路径
	absBase, err := filepath.Abs(basePath)
	if err != nil {
		return fmt.Errorf("无法解析基础路径: %w", err)
	}

	absTarget, err := filepath.Abs(targetPath)
	if err != nil {
		return fmt.Errorf("无法解析目标路径: %w", err)
	}

	// 检查目标路径是否在基础路径内
	rel, err := filepath.Rel(absBase, absTarget)
	if err != nil {
		return fmt.Errorf("路径验证失败: %w", err)
	}

	// 如果相对路径以 .. 开头，说明目标路径在基础路径之外
	if len(rel) >= 2 && rel[0] == '.' && rel[1] == '.' {
		return fmt.Errorf("不安全的路径: 目标路径在基础路径之外")
	}

	return nil
}

// SanitizeInput 清理用户输入，防止注入攻击
func SanitizeInput(input string) string {
	// 移除危险字符
	dangerous := []string{";", "&", "|", "`", "$", "(", ")", "<", ">", "\n", "\r"}
	result := input

	for _, char := range dangerous {
		result = replaceAll(result, char, "")
	}

	return result
}

func replaceAll(s, old, new string) string {
	result := ""
	for i := 0; i < len(s); i++ {
		if i <= len(s)-len(old) && s[i:i+len(old)] == old {
			result += new
			i += len(old) - 1
		} else {
			result += string(s[i])
		}
	}
	return result
}
