package generator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/RogueCultivators/goon/internal/template"
)

func InitProject(projectName, moduleName string, minimal bool, example bool) error {
	// 创建项目目录（如果不存在）
	if err := os.MkdirAll(projectName, 0o755); err != nil {
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
		"scripts",
		"docs",
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
		if err := os.MkdirAll(filepath.Join(projectName, dir), 0o755); err != nil {
			return err
		}
	}

	data := template.ProjectData{
		ProjectName: projectName,
		ModuleName:  moduleName,
	}

	// 核心文件（minimal 模式）
	coreFiles := map[string]string{
		"main.go":                       "main.go.tmpl",
		"cmd/server/server.go":          "server.go.tmpl",
		"internal/config/config.go":     "config.go.tmpl",
		"internal/middleware/cors.go":   "cors.go.tmpl",
		"internal/middleware/logger.go": "logger_middleware.go.tmpl",
		"internal/router/router.go":     "router.go.tmpl",
		"pkg/response/response.go":      "response.go.tmpl",
		"pkg/logger/logger.go":          "logger.go.tmpl",
		"pkg/errors/errors.go":          "errors.go.tmpl",
		"go.mod":                        "go.mod.tmpl",
		"config.yaml":                   "config.yaml.tmpl",
		".gitignore":                    ".gitignore.tmpl",
		"README.md":                     "README.md.tmpl",
		"Makefile":                      "Makefile.tmpl",
		".env.example":                  ".env.example.tmpl",
		"docker-compose.dev.yml":        "docker-compose.dev.yml.tmpl",
		"scripts/setup.sh":              "setup.sh.tmpl",
		"scripts/seed.sh":               "seed.sh.tmpl",
		"docs/api.md":                   "api.md.tmpl",
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

		// 为脚本文件设置执行权限
		perm := os.FileMode(0o600)
		if filepath.Ext(path) == ".sh" {
			perm = 0o755
		}

		if err := os.WriteFile(fullPath, []byte(content), perm); err != nil {
			return err
		}
	}

	// 如果启用示例模式，生成 user 示例模块
	if example {
		// 创建 user 模块目录
		userDir := filepath.Join(projectName, "internal", "user")
		if err := os.MkdirAll(userDir, 0o755); err != nil {
			return fmt.Errorf("创建 user 模块目录失败: %w", err)
		}

		// 生成示例模块文件
		moduleData := template.ModuleData{
			ModuleName:      "user",
			CapitalizedName: "User",
			ProjectModule:   moduleName,
		}

		exampleFiles := map[string]string{
			"handler.go":    "handler_example.go.tmpl",
			"service.go":    "service_example.go.tmpl",
			"model.go":      "model_example.go.tmpl",
			"repository.go": "repository_example.go.tmpl",
			"schema.go":     "schema_example.go.tmpl",
			"routes.go":     "routes.go.tmpl",
		}

		for fileName, tmplName := range exampleFiles {
			fullPath := filepath.Join(userDir, fileName)

			content, err := renderer.Render(tmplName, moduleData)
			if err != nil {
				return fmt.Errorf("渲染示例模板 %s 失败: %w", tmplName, err)
			}

			if err := os.WriteFile(fullPath, []byte(content), 0o600); err != nil {
				return fmt.Errorf("写入示例文件 %s 失败: %w", fileName, err)
			}
		}

		// 创建迁移目录和示例迁移文件
		migrationsDir := filepath.Join(projectName, "migrations")
		if err := os.MkdirAll(migrationsDir, 0o755); err != nil {
			return fmt.Errorf("创建迁移目录失败: %w", err)
		}

		// 创建 users 表迁移
		upMigration := `CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);
`

		downMigration := `DROP TABLE IF EXISTS users;
`

		upPath := filepath.Join(migrationsDir, "000001_create_users_table.up.sql")
		downPath := filepath.Join(migrationsDir, "000001_create_users_table.down.sql")

		if err := os.WriteFile(upPath, []byte(upMigration), 0o600); err != nil {
			return fmt.Errorf("写入迁移文件失败: %w", err)
		}

		if err := os.WriteFile(downPath, []byte(downMigration), 0o600); err != nil {
			return fmt.Errorf("写入迁移文件失败: %w", err)
		}
	}

	return nil
}
