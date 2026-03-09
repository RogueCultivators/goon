package generator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitProject(t *testing.T) {
	tests := []struct {
		name    string
		opts    InitOptions
		wantErr bool
	}{
		{
			name: "minimal project",
			opts: InitOptions{
				ProjectName: "testapp",
				ModuleName:  "github.com/user/testapp",
				Minimal:     true,
			},
			wantErr: false,
		},
		{
			name: "full project",
			opts: InitOptions{
				ProjectName: "fullapp",
				ModuleName:  "github.com/user/fullapp",
				Database:    "PostgreSQL",
				UseAuth:     true,
				AuthMethod:  "JWT",
				UseDocker:   true,
			},
			wantErr: false,
		},
		{
			name: "simple module name",
			opts: InitOptions{
				ProjectName: "simpleapp",
				ModuleName:  "simpleapp",
				Minimal:     true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			oldWd, _ := os.Getwd()
			defer os.Chdir(oldWd)
			os.Chdir(tmpDir)

			err := InitProject(&tt.opts)

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
			if _, statErr := os.Stat(tt.opts.ProjectName); os.IsNotExist(statErr) {
				t.Errorf("Project directory %s was not created", tt.opts.ProjectName)
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
				dirPath := filepath.Join(tt.opts.ProjectName, dir)
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
				filePath := filepath.Join(tt.opts.ProjectName, file)
				if _, statErr := os.Stat(filePath); os.IsNotExist(statErr) {
					t.Errorf("Core file %s was not created", filePath)
				}
			}

			// Verify go.mod contains correct module name
			goModPath := filepath.Join(tt.opts.ProjectName, "go.mod")
			content, err := os.ReadFile(goModPath)
			if err != nil {
				t.Errorf("Failed to read go.mod: %v", err)
			}
			if !containsSubstring(string(content), tt.opts.ModuleName) {
				t.Errorf("go.mod does not contain module name %s", tt.opts.ModuleName)
			}

			// If not minimal, verify additional directories based on options
			if !tt.opts.Minimal {
				if tt.opts.Database != "\u65e0\u6570\u636e\u5e93" && tt.opts.Database != "" {
					for _, dir := range []string{"internal/sqlc/queries", "internal/sqlc/schema", "pkg/database"} {
						dirPath := filepath.Join(tt.opts.ProjectName, dir)
						if _, statErr := os.Stat(dirPath); os.IsNotExist(statErr) {
							t.Errorf("Database directory %s was not created", dirPath)
						}
					}
				}

				if tt.opts.UseDocker {
					for _, file := range []string{"Dockerfile", "docker-compose.yaml"} {
						filePath := filepath.Join(tt.opts.ProjectName, file)
						if _, statErr := os.Stat(filePath); os.IsNotExist(statErr) {
							t.Errorf("Docker file %s was not created", filePath)
						}
					}
				}

				if tt.opts.UseAuth && (tt.opts.AuthMethod == "JWT" || tt.opts.AuthMethod == "") {
					jwtPath := filepath.Join(tt.opts.ProjectName, "pkg", "jwt", "jwt.go")
					if _, statErr := os.Stat(jwtPath); os.IsNotExist(statErr) {
						t.Errorf("JWT file was not created")
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

	opts := InitOptions{
		ProjectName: projectName,
		ModuleName:  moduleName,
		Minimal:     true,
	}

	// First call
	err := InitProject(&opts)
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
	err = InitProject(&opts)
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
	opts := InitOptions{
		ProjectName: "/invalid/path/that/does/not/exist/project",
		ModuleName:  "test",
		Minimal:     true,
	}
	err := InitProject(&opts)
	if err == nil {
		t.Errorf("InitProject() should fail with invalid path")
	}
}
