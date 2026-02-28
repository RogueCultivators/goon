package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name           string
		configContent  string
		wantErr        bool
		checkLayers    bool
		expectedLayers []string
	}{
		{
			name: "valid config",
			configContent: `templates:
  custom_path: "./custom-templates"
defaults:
  layers:
    - handler
    - service
  auto_register: false
naming:
  style: camelCase
`,
			wantErr:        false,
			checkLayers:    true,
			expectedLayers: []string{"handler", "service"},
		},
		{
			name:           "empty config uses defaults",
			configContent:  ``,
			wantErr:        false,
			checkLayers:    true,
			expectedLayers: []string{"handler", "service", "model", "repository", "schema", "routes"},
		},
		{
			name:          "invalid yaml",
			configContent: `invalid: yaml: content:`,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			oldWd, _ := os.Getwd()
			defer os.Chdir(oldWd)
			os.Chdir(tmpDir)

			if tt.configContent != "" {
				configPath := filepath.Join(tmpDir, ".goonrc.yaml")
				if err := os.WriteFile(configPath, []byte(tt.configContent), 0o644); err != nil {
					t.Fatalf("Failed to create config file: %v", err)
				}
			}

			config, err := Load()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Load() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Load() unexpected error = %v", err)
				return
			}

			if config == nil {
				t.Error("Load() returned nil config")
				return
			}

			if tt.checkLayers {
				if len(config.Defaults.Layers) != len(tt.expectedLayers) {
					t.Errorf("Load() layers count = %d, want %d", len(config.Defaults.Layers), len(tt.expectedLayers))
				}
			}
		})
	}
}

func TestSave(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".goonrc.yaml")

	config := &Config{
		Templates: TemplatesConfig{
			CustomPath: "./templates",
		},
		Defaults: DefaultsConfig{
			Layers:       []string{"handler", "service"},
			AutoRegister: true,
		},
		Naming: NamingConfig{
			Style: "snake_case",
		},
	}

	if err := Save(config, configPath); err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	// 验证文件已创建
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}

	// 验证可以重新加载
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	loaded, err := Load()
	if err != nil {
		t.Fatalf("Failed to reload config: %v", err)
	}

	if loaded.Templates.CustomPath != config.Templates.CustomPath {
		t.Errorf("Reloaded config mismatch")
	}
}

func TestGenerateExample(t *testing.T) {
	example := GenerateExample()
	if example == "" {
		t.Error("GenerateExample() returned empty string")
	}

	// 验证示例配置可以解析
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".goonrc.yaml")
	if err := os.WriteFile(configPath, []byte(example), 0o644); err != nil {
		t.Fatalf("Failed to write example config: %v", err)
	}

	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	_, err := Load()
	if err != nil {
		t.Errorf("Example config is not valid: %v", err)
	}
}
