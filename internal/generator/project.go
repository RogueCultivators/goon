package generator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/RogueCultivators/goon/internal/template"
)

func InitProject(projectName, moduleName string, minimal bool) error {
	// 创建项目目录（如果不存在）
	if err := os.MkdirAll(projectName, 0755); err != nil {
		return err
	}

	renderer, err := template.NewRenderer()
	if err != nil {
		return fmt.Errorf("初始化模板渲染器失败: %w", err)
	}

	// 核心目录（minimal 模式）
	coreDirs := []string{
		"cmd/server",
		"internal/config",
		"internal/middleware",
		"internal/router",
		"pkg/response",
		"pkg/logger",
		"pkg/errors",
	}

	// 完整模式额外目录
	fullDirs := []string{
		"internal/sqlc/queries",
		"internal/sqlc/schema",
		"pkg/validator",
		"pkg/database",
		"pkg/jwt",
		"pkg/utils",
		"pkg/cache",
		"pkg/email",
		"pkg/upload",
		"pkg/pagination",
		"pkg/testutil",
	}

	dirs := coreDirs
	if !minimal {
		dirs = append(dirs, fullDirs...)
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(projectName, dir), 0755); err != nil {
			return err
		}
	}

	data := template.ProjectData{
		ProjectName: projectName,
		ModuleName:  moduleName,
	}

	// 核心文件（minimal 模式）
	coreFiles := map[string]string{
		"main.go":                           "main.go.tmpl",
		"cmd/server/server.go":              "server.go.tmpl",
		"internal/config/config.go":         "config.go.tmpl",
		"internal/middleware/cors.go":       "cors.go.tmpl",
		"internal/middleware/logger.go":     "logger_middleware.go.tmpl",
		"internal/router/router.go":         "router.go.tmpl",
		"pkg/response/response.go":          "response.go.tmpl",
		"pkg/logger/logger.go":              "logger.go.tmpl",
		"pkg/errors/errors.go":              "errors.go.tmpl",
		"go.mod":                            "go.mod.tmpl",
		"config.yaml":                       "config.yaml.tmpl",
		".gitignore":                        ".gitignore.tmpl",
		"README.md":                         "README.md.tmpl",
	}

	// 完整模式额外文件
	fullFiles := map[string]string{
		"internal/middleware/requestid.go":  "requestid.go.tmpl",
		"internal/middleware/ratelimit.go":  "ratelimit.go.tmpl",
		"internal/middleware/auth.go":       "auth.go.tmpl",
		"internal/middleware/permission.go": "permission.go.tmpl",
		"internal/middleware/gzip.go":       "gzip.go.tmpl",
		"pkg/validator/validator.go":        "validator.go.tmpl",
		"pkg/database/database.go":          "database.go.tmpl",
		"pkg/jwt/jwt.go":                    "jwt.go.tmpl",
		"pkg/utils/utils.go":                "utils.go.tmpl",
		"pkg/cache/cache.go":                "cache.go.tmpl",
		"pkg/email/email.go":                "email.go.tmpl",
		"pkg/upload/upload.go":              "upload.go.tmpl",
		"pkg/pagination/pagination.go":      "pagination.go.tmpl",
		"pkg/testutil/testutil.go":          "testutil.go.tmpl",
		"sqlc.yaml":                         "sqlc.yaml.tmpl",
		"Dockerfile":                        "Dockerfile.tmpl",
		"docker-compose.yaml":               "docker-compose.yaml.tmpl",
	}

	files := coreFiles
	if !minimal {
		for k, v := range fullFiles {
			files[k] = v
		}
	}

	for path, tmplName := range files {
		fullPath := filepath.Join(projectName, path)

		// 跳过已存在的文件（幂等性）
		if _, err := os.Stat(fullPath); err == nil {
			continue
		}

		content, err := renderer.Render(tmplName, data)
		if err != nil {
			return fmt.Errorf("渲染模板 %s 失败: %w", tmplName, err)
		}

		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return err
		}
	}

	return nil
}
