package generator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/RogueCultivators/goon/internal/template"
)

// 可用的功能包及其对应的模板文件
var availablePackages = map[string]map[string]string{
	"validator": {
		"pkg/validator/validator.go": "validator.go.tmpl",
	},
	"database": {
		"pkg/database/database.go": "database.go.tmpl",
	},
	"jwt": {
		"pkg/jwt/jwt.go": "jwt.go.tmpl",
	},
	"utils": {
		"pkg/utils/utils.go": "utils.go.tmpl",
	},
	"cache": {
		"pkg/cache/cache.go": "cache.go.tmpl",
	},
	"email": {
		"pkg/email/email.go": "email.go.tmpl",
	},
	"upload": {
		"pkg/upload/upload.go": "upload.go.tmpl",
	},
	"pagination": {
		"pkg/pagination/pagination.go": "pagination.go.tmpl",
	},
	"testutil": {
		"pkg/testutil/testutil.go": "testutil.go.tmpl",
	},
}

// AddPackage 添加功能包到项目
func AddPackage(pkgName string) error {
	// 检查是否在项目根目录
	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		return fmt.Errorf("当前目录不是一个 Go 项目，请先运行 'goon init' 初始化项目")
	}

	// 检查包是否存在
	files, ok := availablePackages[pkgName]
	if !ok {
		return fmt.Errorf("未知的功能包: %s\n可用的包: validator, database, jwt, utils, cache, email, upload, pagination, testutil", pkgName)
	}

	gen, err := NewGenerator()
	if err != nil {
		return err
	}

	// 创建模板数据
	data := template.ProjectData{
		ProjectName: "pkg", // 功能包使用固定名称
		ModuleName:  gen.GetProjectModule(),
	}

	// 生成文件
	for path, tmplName := range files {
		// 跳过已存在的文件（幂等性）
		if _, err := os.Stat(path); err == nil {
			fmt.Printf("⚠ 文件已存在，跳过: %s\n", path)
			continue
		}

		// 确保目录存在
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("创建目录失败: %w", err)
		}

		content, err := gen.GetRenderer().Render(tmplName, data)
		if err != nil {
			return fmt.Errorf("渲染模板 %s 失败: %w", tmplName, err)
		}

		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			return fmt.Errorf("写入文件失败: %w", err)
		}

		fmt.Printf("✓ 已生成: %s\n", path)
	}

	return nil
}

// ListAvailablePackages 列出所有可用的功能包
func ListAvailablePackages() []string {
	pkgs := make([]string, 0, len(availablePackages))
	for pkg := range availablePackages {
		pkgs = append(pkgs, pkg)
	}
	return pkgs
}
