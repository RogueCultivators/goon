package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/RogueCultivators/goon/internal/config"
	"github.com/RogueCultivators/goon/internal/template"
	"github.com/RogueCultivators/goon/internal/utils"
)

// InitOptions 项目初始化选项
type InitOptions struct {
	ProjectName   string
	ModuleName    string
	Minimal       bool
	Example       bool
	Database      string // PostgreSQL, MySQL, SQLite, 无数据库
	UseAuth       bool
	AuthMethod    string // JWT, Session
	UseDocker     bool
	ExampleModule string // 示例模块名称，默认 "user"
}

func InitProject(opts *InitOptions) error {
	opts.ProjectName = utils.SanitizeInput(opts.ProjectName)
	if opts.ProjectName == "" {
		return fmt.Errorf("项目名称不能为空")
	}

	if err := utils.ValidatePath(".", opts.ProjectName); err != nil {
		return fmt.Errorf("不安全的项目路径: %w", err)
	}

	if err := os.MkdirAll(opts.ProjectName, 0o755); err != nil {
		return err
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("加载配置失败: %w", err)
	}

	renderer, err := template.NewRenderer(cfg.Templates.CustomPath)
	if err != nil {
		return fmt.Errorf("初始化模板渲染器失败: %w", err)
	}

	if err := createProjectDirectories(opts); err != nil {
		return err
	}

	data := template.ProjectData{
		ProjectName: opts.ProjectName,
		ModuleName:  opts.ModuleName,
	}

	if err := renderProjectFiles(opts, renderer, data); err != nil {
		return err
	}

	if opts.Example {
		exampleModule := opts.ExampleModule
		if exampleModule == "" {
			exampleModule = "user"
		}
		if err := generateExampleModule(opts.ProjectName, opts.ModuleName, exampleModule, renderer); err != nil {
			return err
		}
	}

	return nil
}

func createProjectDirectories(opts *InitOptions) error {
	coreDirs := []string{
		"cmd/server", "internal/config", "internal/middleware", "internal/router",
		"pkg/response", "pkg/logger", "pkg/errors", "scripts", "docs",
	}

	dirs := coreDirs
	if !opts.Minimal {
		// 根据选项添加对应的目录
		dirs = append(dirs, "pkg/validator", "pkg/utils", "pkg/pagination", "pkg/testutil")

		if opts.Database != "无数据库" && opts.Database != "" {
			dirs = append(dirs, "internal/sqlc/queries", "internal/sqlc/schema", "pkg/database")
		}
		if opts.UseAuth {
			if opts.AuthMethod == "JWT" || opts.AuthMethod == "" {
				dirs = append(dirs, "pkg/jwt")
			}
		}
		dirs = append(dirs, "pkg/cache", "pkg/email", "pkg/upload")
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(opts.ProjectName, dir), 0o755); err != nil {
			return err
		}
	}
	return nil
}

func renderProjectFiles(opts *InitOptions, renderer *template.Renderer, data template.ProjectData) error {
	files := getCoreFiles(opts)
	if !opts.Minimal {
		for k, v := range getFullFiles(opts) {
			files[k] = v
		}
	}

	for path, tmplName := range files {
		fullPath := filepath.Join(opts.ProjectName, path)
		if _, err := os.Stat(fullPath); err == nil {
			continue
		}

		content, err := renderer.Render(tmplName, data)
		if err != nil {
			return fmt.Errorf("渲染模板 %s 失败: %w", tmplName, err)
		}

		perm := os.FileMode(0o600)
		if filepath.Ext(path) == ".sh" {
			perm = 0o755
		}

		if err := os.WriteFile(fullPath, []byte(content), perm); err != nil {
			return err
		}
	}
	return nil
}

func getCoreFiles(opts *InitOptions) map[string]string {
	files := map[string]string{
		"main.go": "main.go.tmpl", "cmd/server/server.go": "server.go.tmpl",
		"internal/config/config.go": "config.go.tmpl", "internal/middleware/cors.go": "cors.go.tmpl",
		"internal/middleware/logger.go": "logger_middleware.go.tmpl", "internal/router/router.go": "router.go.tmpl",
		"pkg/response/response.go": "response.go.tmpl", "pkg/logger/logger.go": "logger.go.tmpl",
		"pkg/errors/errors.go": "errors.go.tmpl", "go.mod": "go.mod.tmpl",
		"config.yaml": "config.yaml.tmpl", ".gitignore": ".gitignore.tmpl",
		"README.md": "README.md.tmpl", "Makefile": "Makefile.tmpl",
		".env.example":     ".env.example.tmpl",
		"scripts/setup.sh": "setup.sh.tmpl", "scripts/seed.sh": "seed.sh.tmpl",
		"docs/api.md": "api.md.tmpl",
	}

	if opts.UseDocker {
		files["docker-compose.dev.yml"] = "docker-compose.dev.yml.tmpl"
	}

	return files
}

func getFullFiles(opts *InitOptions) map[string]string {
	files := map[string]string{
		"internal/middleware/requestid.go": "requestid.go.tmpl", "internal/middleware/ratelimit.go": "ratelimit.go.tmpl",
		"internal/middleware/gzip.go": "gzip.go.tmpl", "pkg/validator/validator.go": "validator.go.tmpl",
		"pkg/utils/utils.go": "utils.go.tmpl", "pkg/cache/cache.go": "cache.go.tmpl",
		"pkg/email/email.go": "email.go.tmpl", "pkg/upload/upload.go": "upload.go.tmpl",
		"pkg/pagination/pagination.go": "pagination.go.tmpl", "pkg/testutil/testutil.go": "testutil.go.tmpl",
	}

	// 数据库相关文件
	if opts.Database != "无数据库" && opts.Database != "" {
		files["pkg/database/database.go"] = "database.go.tmpl"
		files["sqlc.yaml"] = "sqlc.yaml.tmpl"
	}

	// 认证相关文件
	if opts.UseAuth {
		files["internal/middleware/auth.go"] = "auth.go.tmpl"
		files["internal/middleware/permission.go"] = "permission.go.tmpl"
		if opts.AuthMethod == "JWT" || opts.AuthMethod == "" {
			files["pkg/jwt/jwt.go"] = "jwt.go.tmpl"
		}
	}

	// Docker 相关文件
	if opts.UseDocker {
		files["Dockerfile"] = "Dockerfile.tmpl"
		files["docker-compose.yaml"] = "docker-compose.yaml.tmpl"
	}

	return files
}

func generateExampleModule(projectName, moduleName, exampleModule string, renderer *template.Renderer) error {
	moduleDir := filepath.Join(projectName, "internal", exampleModule)
	if err := os.MkdirAll(moduleDir, 0o755); err != nil {
		return fmt.Errorf("创建 %s 模块目录失败: %w", exampleModule, err)
	}

	capitalizedName := strings.ToUpper(exampleModule[:1]) + exampleModule[1:]
	moduleData := template.ModuleData{
		ModuleName:      exampleModule,
		CapitalizedName: capitalizedName,
		ProjectModule:   moduleName,
	}

	exampleFiles := map[string]string{
		"handler.go": "handler_example.go.tmpl", "service.go": "service_example.go.tmpl",
		"model.go": "model_example.go.tmpl", "repository.go": "repository_example.go.tmpl",
		"schema.go": "schema_example.go.tmpl", "routes.go": "routes.go.tmpl",
	}

	for fileName, tmplName := range exampleFiles {
		fullPath := filepath.Join(moduleDir, fileName)
		content, err := renderer.Render(tmplName, moduleData)
		if err != nil {
			return fmt.Errorf("渲染示例模板 %s 失败: %w", tmplName, err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0o600); err != nil {
			return fmt.Errorf("写入示例文件 %s 失败: %w", fileName, err)
		}
	}

	return createMigrationFiles(projectName)
}

func createMigrationFiles(projectName string) error {
	migrationsDir := filepath.Join(projectName, "migrations")
	if err := os.MkdirAll(migrationsDir, 0o755); err != nil {
		return fmt.Errorf("创建迁移目录失败: %w", err)
	}

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

	return nil
}
