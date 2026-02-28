package generator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/RogueCultivators/goon/internal/template"
	"github.com/RogueCultivators/goon/internal/utils"
)

func AddModule(moduleName string, layers []string, example bool, dryRun bool) error {
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

	moduleName = utils.ToSnakeCase(moduleName)
	capitalizedName := utils.ToPascalCase(moduleName)
	moduleDir := filepath.Join("internal", moduleName)

	if err := os.MkdirAll(moduleDir, 0o755); err != nil {
		return err
	}

	data := template.ModuleData{
		ModuleName:      moduleName,
		CapitalizedName: capitalizedName,
		ProjectModule:   projectModule,
	}

	allFiles := getModuleTemplates(example)
	if len(layers) == 0 {
		layers = []string{"handler", "service", "model", "repository", "schema", "routes"}
	}

	if err := validateLayers(layers, allFiles); err != nil {
		return err
	}

	files := buildFileList(moduleDir, layers, allFiles)
	return renderModuleFiles(files, renderer, data, dryRun)
}

func getModuleTemplates(example bool) map[string]string {
	if example {
		return map[string]string{
			"handler": "handler_example.go.tmpl", "service": "service_example.go.tmpl",
			"model": "model_example.go.tmpl", "repository": "repository_example.go.tmpl",
			"schema": "schema_example.go.tmpl", "routes": "routes.go.tmpl",
		}
	}
	return map[string]string{
		"handler": "handler.go.tmpl", "service": "service.go.tmpl",
		"model": "model.go.tmpl", "repository": "repository.go.tmpl",
		"schema": "schema.go.tmpl", "routes": "routes.go.tmpl",
	}
}

func validateLayers(layers []string, allFiles map[string]string) error {
	validLayers := make(map[string]bool)
	for layer := range allFiles {
		validLayers[layer] = true
	}

	for _, layer := range layers {
		if !validLayers[layer] {
			return fmt.Errorf("无效的层: %s，可用的层: handler, service, model, repository, schema", layer)
		}
	}
	return nil
}

func buildFileList(moduleDir string, layers []string, allFiles map[string]string) map[string]string {
	files := make(map[string]string)
	for _, layer := range layers {
		if tmpl, ok := allFiles[layer]; ok {
			files[filepath.Join(moduleDir, layer+".go")] = tmpl
		}
	}
	return files
}

func renderModuleFiles(files map[string]string, renderer *template.Renderer, data template.ModuleData, dryRun bool) error {
	for path, tmplName := range files {
		if _, err := os.Stat(path); err == nil {
			if dryRun {
				fmt.Printf("  ⏭  %s (已存在，将跳过)\n", path)
			}
			continue
		}

		if dryRun {
			fmt.Printf("  ✓ %s\n", path)
			fmt.Printf("     模板: %s\n", tmplName)
			continue
		}

		content, err := renderer.Render(tmplName, data)
		if err != nil {
			return fmt.Errorf("渲染模板 %s 失败: %w", tmplName, err)
		}

		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			return err
		}
	}
	return nil
}

// getModuleNameFromGoMod 已移至 modfile.go
