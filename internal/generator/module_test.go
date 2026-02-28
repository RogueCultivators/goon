package generator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/RogueCultivators/goon/internal/utils"
)

func TestAddModule(t *testing.T) {
	tests := []struct {
		name        string
		moduleName  string
		layers      []string
		wantErr     bool
		errContains string
	}{
		{
			name:       "valid module with all layers",
			moduleName: "user",
			layers:     []string{"handler", "service", "model", "repository", "schema"},
			wantErr:    false,
		},
		{
			name:       "valid module with some layers",
			moduleName: "product",
			layers:     []string{"handler", "service"},
			wantErr:    false,
		},
		{
			name:       "valid module with no layers (should generate all)",
			moduleName: "order",
			layers:     []string{},
			wantErr:    false,
		},
		{
			name:        "invalid layer",
			moduleName:  "test",
			layers:      []string{"invalid_layer"},
			wantErr:     true,
			errContains: "无效的层",
		},
		{
			name:       "camelCase module name",
			moduleName: "userProfile",
			layers:     []string{"handler"},
			wantErr:    false,
		},
		{
			name:       "PascalCase module name",
			moduleName: "UserProfile",
			layers:     []string{"handler"},
			wantErr:    false,
		},
		{
			name:       "kebab-case module name",
			moduleName: "user-profile",
			layers:     []string{"handler"},
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary test directory
			tmpDir := t.TempDir()
			oldWd, _ := os.Getwd()
			defer os.Chdir(oldWd)
			os.Chdir(tmpDir)

			// Create a minimal go.mod file
			goModContent := "module testproject\n\ngo 1.24\n"
			if err := os.WriteFile("go.mod", []byte(goModContent), 0o644); err != nil {
				t.Fatalf("Failed to create go.mod: %v", err)
			}

			err := AddModule(tt.moduleName, tt.layers, false, false)

			if tt.wantErr {
				if err == nil {
					t.Errorf("AddModule() expected error but got none")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("AddModule() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("AddModule() unexpected error = %v", err)
				return
			}

			// Verify module directory was created
			moduleDir := filepath.Join("internal", normalizeModuleName(tt.moduleName))
			if _, err := os.Stat(moduleDir); os.IsNotExist(err) {
				t.Errorf("Module directory %s was not created", moduleDir)
			}

			// Verify expected files were created
			expectedLayers := tt.layers
			if len(expectedLayers) == 0 {
				expectedLayers = []string{"handler", "service", "model", "repository", "schema"}
			}

			for _, layer := range expectedLayers {
				filePath := filepath.Join(moduleDir, layer+".go")
				if _, err := os.Stat(filePath); os.IsNotExist(err) {
					t.Errorf("Expected file %s was not created", filePath)
				}
			}
		})
	}
}

func TestGetModuleNameFromGoMod(t *testing.T) {
	tests := []struct {
		name         string
		goModContent string
		want         string
		wantErr      bool
	}{
		{
			name:         "valid go.mod",
			goModContent: "module github.com/user/project\n\ngo 1.24\n",
			want:         "github.com/user/project",
			wantErr:      false,
		},
		{
			name:         "go.mod with extra spaces",
			goModContent: "module   github.com/user/project  \n\ngo 1.24\n",
			want:         "github.com/user/project",
			wantErr:      false,
		},
		{
			name:         "simple module name",
			goModContent: "module myproject\n\ngo 1.24\n",
			want:         "myproject",
			wantErr:      false,
		},
		{
			name:         "no module line",
			goModContent: "go 1.24\n",
			want:         "",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			oldWd, _ := os.Getwd()
			defer os.Chdir(oldWd)
			os.Chdir(tmpDir)

			if err := os.WriteFile("go.mod", []byte(tt.goModContent), 0644); err != nil {
				t.Fatalf("Failed to create go.mod: %v", err)
			}

			got, err := getModuleNameFromGoMod()

			if tt.wantErr {
				if err == nil {
					t.Errorf("getModuleNameFromGoMod() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("getModuleNameFromGoMod() unexpected error = %v", err)
				return
			}

			if got != tt.want {
				t.Errorf("getModuleNameFromGoMod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddModuleIdempotency(t *testing.T) {
	tmpDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	goModContent := "module testproject\n\ngo 1.24\n"
	if err := os.WriteFile("go.mod", []byte(goModContent), 0644); err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	// First call
	err := AddModule("user", []string{"handler"}, false, false)
	if err != nil {
		t.Fatalf("First AddModule() call failed: %v", err)
	}

	// Get file info
	filePath := filepath.Join("internal", "user", "handler.go")
	info1, err := os.Stat(filePath)
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}

	// Second call (should be idempotent)
	err = AddModule("user", []string{"handler"}, false, false)
	if err != nil {
		t.Fatalf("Second AddModule() call failed: %v", err)
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

// Helper functions
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func normalizeModuleName(name string) string {
	// Use the same normalization as the actual implementation
	return utils.ToSnakeCase(name)
}
