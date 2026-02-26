package template

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"
)

//go:embed templates/project/*.tmpl templates/module/*.tmpl
var templateFS embed.FS

type Renderer struct {
	templates *template.Template
}

func NewRenderer() (*Renderer, error) {
	tmpl, err := template.ParseFS(templateFS, "templates/**/*.tmpl")
	if err != nil {
		return nil, fmt.Errorf("解析模板失败: %w", err)
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

type ProjectData struct {
	ProjectName string
	ModuleName  string
}

type ModuleData struct {
	ModuleName      string
	CapitalizedName string
	ProjectModule   string
}
