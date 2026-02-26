package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBackupManager(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	originalContent := []byte("original content")

	// 创建测试文件
	if err := os.WriteFile(testFile, originalContent, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	bm := NewBackupManager()

	// 备份文件
	if err := bm.BackupFile(testFile); err != nil {
		t.Fatalf("BackupFile() failed: %v", err)
	}

	// 修改文件
	newContent := []byte("modified content")
	if err := os.WriteFile(testFile, newContent, 0644); err != nil {
		t.Fatalf("Failed to modify file: %v", err)
	}

	// 验证文件已修改
	content, _ := os.ReadFile(testFile)
	if string(content) != string(newContent) {
		t.Errorf("File was not modified correctly")
	}

	// 回滚
	if err := bm.Rollback(); err != nil {
		t.Fatalf("Rollback() failed: %v", err)
	}

	// 验证文件已恢复
	content, _ = os.ReadFile(testFile)
	if string(content) != string(originalContent) {
		t.Errorf("File was not restored correctly, got %s, want %s", string(content), string(originalContent))
	}
}

func TestBackupManagerNewFile(t *testing.T) {
	tmpDir := t.TempDir()
	newFile := filepath.Join(tmpDir, "new.txt")

	bm := NewBackupManager()

	// 备份不存在的文件
	if err := bm.BackupFile(newFile); err != nil {
		t.Fatalf("BackupFile() failed for new file: %v", err)
	}

	// 创建新文件
	if err := os.WriteFile(newFile, []byte("new content"), 0644); err != nil {
		t.Fatalf("Failed to create new file: %v", err)
	}

	// 回滚（应该删除新文件）
	if err := bm.Rollback(); err != nil {
		t.Fatalf("Rollback() failed: %v", err)
	}

	// 验证文件已删除
	if _, err := os.Stat(newFile); !os.IsNotExist(err) {
		t.Errorf("New file was not deleted during rollback")
	}
}

func TestValidatePath(t *testing.T) {
	tests := []struct {
		name       string
		basePath   string
		targetPath string
		wantErr    bool
	}{
		{
			name:       "valid path within base",
			basePath:   "/home/user/project",
			targetPath: "/home/user/project/src/file.go",
			wantErr:    false,
		},
		{
			name:       "path traversal attack",
			basePath:   "/home/user/project",
			targetPath: "/home/user/project/../../../etc/passwd",
			wantErr:    true,
		},
		{
			name:       "relative path within base",
			basePath:   "/home/user/project",
			targetPath: "/home/user/project/./src/file.go",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePath(tt.basePath, tt.targetPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSanitizeInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "clean input",
			input: "user_name",
			want:  "user_name",
		},
		{
			name:  "input with semicolon",
			input: "user; rm -rf /",
			want:  "user rm -rf /",
		},
		{
			name:  "input with pipe",
			input: "user | cat /etc/passwd",
			want:  "user  cat /etc/passwd",
		},
		{
			name:  "input with backticks",
			input: "user`whoami`",
			want:  "userwhoami",
		},
		{
			name:  "multiple dangerous chars",
			input: "user;$(whoami)&|",
			want:  "userwhoami",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizeInput(tt.input)
			if got != tt.want {
				t.Errorf("SanitizeInput() = %v, want %v", got, tt.want)
			}
		})
	}
}
