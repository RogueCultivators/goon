package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/RogueCultivators/goon/internal/utils"
)

// RegisterModuleRoute 自动在 router.go 中注册模块路由
func RegisterModuleRoute(moduleName string) error {
	routerPath := filepath.Join("internal", "router", "router.go")

	// 检查 router.go 是否存在
	if _, err := os.Stat(routerPath); os.IsNotExist(err) {
		return fmt.Errorf("router.go 不存在，请确保在项目根目录运行此命令")
	}

	// 读取文件内容
	content, err := os.ReadFile(routerPath)
	if err != nil {
		return fmt.Errorf("读取 router.go 失败: %w", err)
	}

	fileContent := string(content)

	// 标准化模块名称
	moduleName = utils.ToSnakeCase(moduleName)

	// 获取项目模块名
	projectModule, err := getModuleNameFromGoMod()
	if err != nil {
		return err
	}

	// 构建 import 语句
	importLine := fmt.Sprintf("\t\"%s/internal/%s\"", projectModule, moduleName)

	// 检查是否已经导入
	if strings.Contains(fileContent, importLine) {
		return fmt.Errorf("模块 %s 的路由已经注册", moduleName)
	}

	// 在 import 块中添加新的导入
	fileContent = addImport(fileContent, importLine)

	// 构建路由注册代码（使用新的 RegisterRoutes 函数）
	routeCode := fmt.Sprintf("\t\t%s.RegisterRoutes(api)\n", moduleName)

	// 在注释标记后添加路由代码
	marker := "// 在这里调用各个模块的 RegisterRoutes 函数"
	if strings.Contains(fileContent, marker) {
		fileContent = strings.Replace(fileContent, marker, marker+"\n"+routeCode, 1)
	} else {
		// 如果没有标记，在 api.Group 块的末尾添加
		fileContent = addRouteToAPIGroup(fileContent, routeCode)
	}

	// 写回文件
	if err := os.WriteFile(routerPath, []byte(fileContent), 0644); err != nil {
		return fmt.Errorf("写入 router.go 失败: %w", err)
	}

	return nil
}

// addImport 在 import 块中添加新的导入
func addImport(content, importLine string) string {
	// 查找 import 块
	importStart := strings.Index(content, "import (")
	if importStart == -1 {
		return content
	}

	// 查找 import 块的结束位置
	importEnd := strings.Index(content[importStart:], ")")
	if importEnd == -1 {
		return content
	}

	importEnd += importStart

	// 在 import 块末尾添加新的导入
	before := content[:importEnd]
	after := content[importEnd:]

	return before + "\n" + importLine + "\n" + after
}

// addRouteToAPIGroup 在 API 路由组中添加路由代码
func addRouteToAPIGroup(content, routeCode string) string {
	// 查找 api := r.Group("/api/v1") 后的代码块
	apiGroupStart := strings.Index(content, `api := r.Group("/api/v1")`)
	if apiGroupStart == -1 {
		return content
	}

	// 查找该代码块的结束大括号
	braceCount := 0
	inBlock := false
	insertPos := -1

	for i := apiGroupStart; i < len(content); i++ {
		if content[i] == '{' {
			braceCount++
			inBlock = true
		} else if content[i] == '}' {
			braceCount--
			if inBlock && braceCount == 0 {
				// 找到了匹配的结束大括号
				insertPos = i
				break
			}
		}
	}

	if insertPos == -1 {
		return content
	}

	// 在结束大括号前插入路由代码
	before := content[:insertPos]
	after := content[insertPos:]

	return before + routeCode + "\t" + after
}
