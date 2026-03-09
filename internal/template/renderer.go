package template

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"text/template"
)

//go:embed templates/project/*.tmpl templates/module/*.tmpl
var templateFS embed.FS

type Renderer struct {
	templates *template.Template
}

func NewRenderer(customPaths ...string) (*Renderer, error) {
	tmpl, err := template.ParseFS(templateFS, "templates/**/*.tmpl")
	if err != nil {
		return nil, fmt.Errorf("解析内置模板失败: %w", err)
	}

	// 如果提供了自定义模板路径，加载并覆盖内置模板
	if len(customPaths) > 0 && customPaths[0] != "" {
		customPath := customPaths[0]
		if _, statErr := os.Stat(customPath); statErr == nil {
			pattern := filepath.Join(customPath, "*.tmpl")
			tmpl, err = tmpl.ParseGlob(pattern)
			if err != nil {
				return nil, fmt.Errorf("解析自定义模板 %s 失败: %w", customPath, err)
			}
		}
	}

	return &Renderer{
		templates: tmpl,
	}, nil
}

func (r *Renderer) Render(templateName string, data interface{}) (string, error) {
	var buf bytes.Buffer
	if err := r.templates.ExecuteTemplate(&buf, templateName, data); err != nil {
		return "", fmt.Errorf("渲染模板 %s 失败: %w", templateName, err)
	}
	return buf.String(), nil
}

// ListTemplates 从 embed.FS 动态读取模板文件列表，按类别分组返回
func ListTemplates() (projectTemplates []string, moduleTemplates []string) {
	if entries, err := fs.ReadDir(templateFS, "templates/project"); err == nil {
		for _, e := range entries {
			if !e.IsDir() {
				projectTemplates = append(projectTemplates, e.Name())
			}
		}
	}
	if entries, err := fs.ReadDir(templateFS, "templates/module"); err == nil {
		for _, e := range entries {
			if !e.IsDir() {
				moduleTemplates = append(moduleTemplates, e.Name())
			}
		}
	}
	return
}

type ProjectData struct {
	ProjectName string
	ModuleName  string
}

type ModuleData struct {
	ModuleName      string
	CapitalizedName string
	ProjectModule   string
}
