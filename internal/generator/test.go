package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/RogueCultivators/goon/internal/template"
	"github.com/RogueCultivators/goon/internal/utils"
)

// GenerateModuleTests 为模块生成测试文件
func GenerateModuleTests(moduleName string, layers []string) error {
	moduleName = utils.ToSnakeCase(moduleName)
	moduleDir := filepath.Join("internal", moduleName)

	// 检查模块是否存在
	if _, err := os.Stat(moduleDir); os.IsNotExist(err) {
		return fmt.Errorf("模块 %s 不存在", moduleName)
	}

	// 获取项目模块名
	projectModule, err := getModuleNameFromGoMod()
	if err != nil {
		return err
	}

	renderer, err := template.NewRenderer()
	if err != nil {
		return fmt.Errorf("初始化模板渲染器失败: %w", err)
	}

	capitalizedName := utils.ToPascalCase(moduleName)

	data := template.ModuleData{
		ModuleName:      moduleName,
		CapitalizedName: capitalizedName,
		ProjectModule:   projectModule,
	}

	// 如果没有指定 layers，为所有存在的文件生成测试
	if len(layers) == 0 {
		files, err := filepath.Glob(filepath.Join(moduleDir, "*.go"))
		if err != nil {
			return fmt.Errorf("查找文件失败: %w", err)
		}

		for _, file := range files {
			base := filepath.Base(file)
			if strings.HasSuffix(base, "_test.go") {
				continue
			}
			layer := strings.TrimSuffix(base, ".go")
			layers = append(layers, layer)
		}
	}

	// 测试模板映射
	testTemplates := map[string]string{
		"handler":    "handler_test.go.tmpl",
		"service":    "service_test.go.tmpl",
		"repository": "repository_test.go.tmpl",
	}

	for _, layer := range layers {
		tmplName, ok := testTemplates[layer]
		if !ok {
			continue // 跳过没有测试模板的层
		}

		testFile := filepath.Join(moduleDir, layer+"_test.go")

		// 跳过已存在的测试文件
		if _, err := os.Stat(testFile); err == nil {
			continue
		}

		content, err := renderer.Render(tmplName, data)
		if err != nil {
			return fmt.Errorf("渲染模板 %s 失败: %w", tmplName, err)
		}

		if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
			return fmt.Errorf("写入文件 %s 失败: %w", testFile, err)
		}
	}

	return nil
}

// GenerateAllTests 为所有模块生成测试文件
func GenerateAllTests() error {
	internalDir := "internal"

	entries, err := os.ReadDir(internalDir)
	if err != nil {
		return fmt.Errorf("读取 internal 目录失败: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// 跳过特殊目录
		if entry.Name() == "config" || entry.Name() == "middleware" ||
		   entry.Name() == "router" || entry.Name() == "sqlc" ||
		   entry.Name() == "template" {
			continue
		}

		moduleName := entry.Name()
		if err := GenerateModuleTests(moduleName, []string{}); err != nil {
			return fmt.Errorf("为模块 %s 生成测试失败: %w", moduleName, err)
		}
	}

	return nil
}
