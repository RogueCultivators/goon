package generator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/RogueCultivators/goon/internal/template"
	"github.com/RogueCultivators/goon/internal/utils"
)

func AddModule(moduleName string, layers []string) error {
	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		return fmt.Errorf("当前目录不是一个 Go 项目，请先运行 'goon init' 初始化项目")
	}

	projectModule, err := getModuleNameFromGoMod()
	if err != nil {
		return err
	}

	renderer, err := template.NewRenderer()
	if err != nil {
		return fmt.Errorf("初始化模板渲染器失败: %w", err)
	}

	// 标准化模块名称
	moduleName = utils.ToSnakeCase(moduleName)
	capitalizedName := utils.ToPascalCase(moduleName)

	moduleDir := filepath.Join("internal", moduleName)

	// 创建模块目录（如果不存在）
	if err := os.MkdirAll(moduleDir, 0755); err != nil {
		return err
	}

	data := template.ModuleData{
		ModuleName:      moduleName,
		CapitalizedName: capitalizedName,
		ProjectModule:   projectModule,
	}

	// 所有可用的层
	allFiles := map[string]string{
		"handler":    "handler.go.tmpl",
		"service":    "service.go.tmpl",
		"model":      "model.go.tmpl",
		"repository": "repository.go.tmpl",
		"schema":     "schema.go.tmpl",
		"routes":     "routes.go.tmpl",
	}

	// 如果没有指定 layers，生成所有文件
	if len(layers) == 0 {
		layers = []string{"handler", "service", "model", "repository", "schema", "routes"}
	}

	// 验证 layers 是否有效
	validLayers := make(map[string]bool)
	for layer := range allFiles {
		validLayers[layer] = true
	}

	for _, layer := range layers {
		if !validLayers[layer] {
			return fmt.Errorf("无效的层: %s，可用的层: handler, service, model, repository, schema", layer)
		}
	}

	// 根据指定的 layers 构建文件列表
	files := make(map[string]string)
	for _, layer := range layers {
		if tmpl, ok := allFiles[layer]; ok {
			files[filepath.Join(moduleDir, layer+".go")] = tmpl
		}
	}

	for path, tmplName := range files {
		// 跳过已存在的文件（幂等性）
		if _, err := os.Stat(path); err == nil {
			continue
		}

		content, err := renderer.Render(tmplName, data)
		if err != nil {
			return fmt.Errorf("渲染模板 %s 失败: %w", tmplName, err)
		}

		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return err
		}
	}

	return nil
}

// getModuleNameFromGoMod 已移至 modfile.go
