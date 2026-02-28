package generator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitProject(t *testing.T) {
	tests := []struct {
		name        string
		projectName string
		moduleName  string
		minimal     bool
		wantErr     bool
	}{
		{
			name:        "minimal project",
			projectName: "testapp",
			moduleName:  "github.com/user/testapp",
			minimal:     true,
			wantErr:     false,
		},
		{
			name:        "full project",
			projectName: "fullapp",
			moduleName:  "github.com/user/fullapp",
			minimal:     false,
			wantErr:     false,
		},
		{
			name:        "simple module name",
			projectName: "simpleapp",
			moduleName:  "simpleapp",
			minimal:     true,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			oldWd, _ := os.Getwd()
			defer os.Chdir(oldWd)
			os.Chdir(tmpDir)

			err := InitProject(tt.projectName, tt.moduleName, tt.minimal, false)

			if tt.wantErr {
				if err == nil {
					t.Errorf("InitProject() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("InitProject() unexpected error = %v", err)
				return
			}

			// Verify project directory was created
			if _, statErr := os.Stat(tt.projectName); os.IsNotExist(statErr) {
				t.Errorf("Project directory %s was not created", tt.projectName)
			}

			// Verify core directories exist
			coreDirs := []string{
				"cmd/server",
				"internal/config",
				"internal/middleware",
				"internal/router",
				"pkg/response",
				"pkg/logger",
				"pkg/errors",
			}

			for _, dir := range coreDirs {
				dirPath := filepath.Join(tt.projectName, dir)
				if _, statErr := os.Stat(dirPath); os.IsNotExist(statErr) {
					t.Errorf("Core directory %s was not created", dirPath)
				}
			}

			// Verify core files exist
			coreFiles := []string{
				"main.go",
				"go.mod",
				"config.yaml",
				".gitignore",
				"README.md",
			}

			for _, file := range coreFiles {
				filePath := filepath.Join(tt.projectName, file)
				if _, statErr := os.Stat(filePath); os.IsNotExist(statErr) {
					t.Errorf("Core file %s was not created", filePath)
				}
			}

			// Verify go.mod contains correct module name
			goModPath := filepath.Join(tt.projectName, "go.mod")
			content, err := os.ReadFile(goModPath)
			if err != nil {
				t.Errorf("Failed to read go.mod: %v", err)
			}
			if !containsSubstring(string(content), tt.moduleName) {
				t.Errorf("go.mod does not contain module name %s", tt.moduleName)
			}

			// If not minimal, verify additional directories
			if !tt.minimal {
				fullDirs := []string{
					"internal/sqlc/queries",
					"internal/sqlc/schema",
					"pkg/validator",
					"pkg/database",
				}

				for _, dir := range fullDirs {
					dirPath := filepath.Join(tt.projectName, dir)
					if _, statErr := os.Stat(dirPath); os.IsNotExist(statErr) {
						t.Errorf("Full mode directory %s was not created", dirPath)
					}
				}

				// Verify additional files
				fullFiles := []string{
					"Dockerfile",
					"docker-compose.yaml",
					"sqlc.yaml",
				}

				for _, file := range fullFiles {
					filePath := filepath.Join(tt.projectName, file)
					if _, statErr := os.Stat(filePath); os.IsNotExist(statErr) {
						t.Errorf("Full mode file %s was not created", filePath)
					}
				}
			}
		})
	}
}

func TestInitProjectIdempotency(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	projectName := "testapp"
	moduleName := "github.com/user/testapp"

	// First call
	err := InitProject(projectName, moduleName, true, false)
	if err != nil {
		t.Fatalf("First InitProject() call failed: %v", err)
	}

	// Get file info
	filePath := filepath.Join(projectName, "main.go")
	info1, err := os.Stat(filePath)
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}

	// Second call (should be idempotent)
	err = InitProject(projectName, moduleName, true, false)
	if err != nil {
		t.Fatalf("Second InitProject() call failed: %v", err)
	}

	// Verify file wasn't modified
	info2, err := os.Stat(filePath)
	if err != nil {
		t.Fatalf("Failed to stat file after second call: %v", err)
	}

	if info1.ModTime() != info2.ModTime() {
		t.Errorf("File was modified on second call, idempotency violated")
	}
}

func TestInitProjectInvalidPath(t *testing.T) {
	// Try to create project in a path that would fail
	err := InitProject("/invalid/path/that/does/not/exist/project", "test", true, false)
	if err == nil {
		t.Errorf("InitProject() should fail with invalid path")
	}
}
