package generator

import (
	"fmt"

	"github.com/RogueCultivators/goon/internal/config"
	"github.com/RogueCultivators/goon/internal/template"
)

// Generator 代码生成器，避免重复初始化
type Generator struct {
	renderer      *template.Renderer
	projectModule string
	config        *config.Config
}

// NewGenerator 创建新的代码生成器实例
func NewGenerator() (*Generator, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("加载配置失败: %w", err)
	}

	renderer, err := template.NewRenderer(cfg.Templates.CustomPath)
	if err != nil {
		return nil, fmt.Errorf("初始化模板渲染器失败: %w", err)
	}

	projectModule, err := getModuleNameFromGoMod()
	if err != nil {
		return nil, fmt.Errorf("获取项目模块名失败: %w", err)
	}

	return &Generator{
		renderer:      renderer,
		projectModule: projectModule,
		config:        cfg,
	}, nil
}

// GetRenderer 获取模板渲染器
func (g *Generator) GetRenderer() *template.Renderer {
	return g.renderer
}

// GetProjectModule 获取项目模块名
func (g *Generator) GetProjectModule() string {
	return g.projectModule
}

// GetConfig 获取配置
func (g *Generator) GetConfig() *config.Config {
	return g.config
}
