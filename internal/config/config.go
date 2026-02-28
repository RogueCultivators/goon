package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config 配置结构
type Config struct {
	Templates TemplatesConfig `yaml:"templates"`
	Defaults  DefaultsConfig  `yaml:"defaults"`
	Naming    NamingConfig    `yaml:"naming"`
}

// TemplatesConfig 模板配置
type TemplatesConfig struct {
	CustomPath string `yaml:"custom_path"`
}

// DefaultsConfig 默认配置
type DefaultsConfig struct {
	Layers       []string `yaml:"layers"`
	AutoRegister bool     `yaml:"auto_register"`
}

// NamingConfig 命名风格配置
type NamingConfig struct {
	Style string `yaml:"style"` // snake_case, camelCase, kebab-case
}

var defaultConfig = Config{
	Templates: TemplatesConfig{
		CustomPath: "",
	},
	Defaults: DefaultsConfig{
		Layers:       []string{"handler", "service", "model", "repository", "schema", "routes"},
		AutoRegister: true,
	},
	Naming: NamingConfig{
		Style: "snake_case",
	},
}

// Load 加载配置文件
func Load() (*Config, error) {
	configPath := findConfigFile()
	if configPath == "" {
		// 没有找到配置文件，使用默认配置
		return &defaultConfig, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 合并默认配置
	if len(config.Defaults.Layers) == 0 {
		config.Defaults.Layers = defaultConfig.Defaults.Layers
	}
	if config.Naming.Style == "" {
		config.Naming.Style = defaultConfig.Naming.Style
	}

	return &config, nil
}

// findConfigFile 查找配置文件
func findConfigFile() string {
	// 查找顺序：当前目录 -> 父目录 -> 用户主目录
	searchPaths := []string{
		".goonrc.yaml",
		".goonrc.yml",
		"../.goonrc.yaml",
		"../.goonrc.yml",
	}

	// 添加用户主目录
	if home, err := os.UserHomeDir(); err == nil {
		searchPaths = append(searchPaths,
			filepath.Join(home, ".goonrc.yaml"),
			filepath.Join(home, ".goonrc.yml"),
		)
	}

	for _, path := range searchPaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return ""
}

// Save 保存配置文件
func Save(config *Config, path string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	return nil
}

// GenerateExample 生成示例配置文件
func GenerateExample() string {
	return `# Goon 配置文件
# 在项目根目录创建 .goonrc.yaml 文件来自定义配置

# 模板配置
templates:
  # 自定义模板目录路径（可选）
  custom_path: ""

# 默认配置
defaults:
  # 默认生成的层
  layers:
    - handler
    - service
    - model
    - repository
    - schema
    - routes
  # 是否自动注册路由
  auto_register: true

# 命名风格配置
naming:
  # 命名风格: snake_case, camelCase, kebab-case
  style: snake_case
`
}
