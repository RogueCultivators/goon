package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "模板管理命令",
	Long:  `管理和查看项目模板`,
}

var templateListCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有可用的模板",
	Long: `显示所有可用的模板文件，包括项目模板和模块模板。

示例:
  goon template list              # 列出所有模板
  goon template list --type=module  # 只列出模块模板
  goon template list --type=project # 只列出项目模板`,
	Run: func(cmd *cobra.Command, args []string) {
		templateType, err := cmd.Flags().GetString("type")
		if err != nil {
			fmt.Printf("获取参数失败: %v\n", err)
			return
		}

		fmt.Println("📋 可用模板列表")

		if templateType == "" || templateType == "project" {
			fmt.Println("🏗️  项目模板 (用于 goon init):")
			fmt.Println("  核心文件:")
			fmt.Println("    - main.go.tmpl                    # 主入口文件")
			fmt.Println("    - go.mod.tmpl                     # Go 模块定义")
			fmt.Println("    - config.yaml.tmpl                # 配置文件")
			fmt.Println("    - .env.example.tmpl               # 环境变量示例")
			fmt.Println("    - README.md.tmpl                  # 项目文档")
			fmt.Println()
			fmt.Println("  服务器:")
			fmt.Println("    - server.go.tmpl                  # 服务器启动逻辑")
			fmt.Println("    - config.go.tmpl                  # 配置管理")
			fmt.Println("    - router.go.tmpl                  # 路由配置")
			fmt.Println()
			fmt.Println("  中间件:")
			fmt.Println("    - cors.go.tmpl                    # CORS 中间件")
			fmt.Println("    - logger_middleware.go.tmpl       # 日志中间件")
			fmt.Println("    - requestid.go.tmpl               # 请求 ID 中间件")
			fmt.Println("    - ratelimit.go.tmpl               # 限流中间件")
			fmt.Println("    - auth.go.tmpl                    # 认证中间件")
			fmt.Println("    - permission.go.tmpl              # 权限中间件")
			fmt.Println("    - gzip.go.tmpl                    # Gzip 压缩中间件")
			fmt.Println()
			fmt.Println("  工具包:")
			fmt.Println("    - response.go.tmpl                # 统一响应格式")
			fmt.Println("    - logger.go.tmpl                  # 日志工具")
			fmt.Println("    - errors.go.tmpl                  # 错误处理")
			fmt.Println("    - database.go.tmpl                # 数据库连接")
			fmt.Println("    - jwt.go.tmpl                     # JWT 工具")
			fmt.Println("    - utils.go.tmpl                   # 通用工具")
			fmt.Println()
			fmt.Println("  开发工具:")
			fmt.Println("    - Makefile.tmpl                   # 开发工具链")
			fmt.Println("    - docker-compose.yaml.tmpl        # Docker 配置")
			fmt.Println("    - docker-compose.dev.yml.tmpl     # 开发环境配置")
			fmt.Println("    - Dockerfile.tmpl                 # Docker 镜像")
			fmt.Println("    - setup.sh.tmpl                   # 初始化脚本")
			fmt.Println("    - seed.sh.tmpl                    # 数据填充脚本")
			fmt.Println()
			fmt.Println("  文档:")
			fmt.Println("    - api.md.tmpl                     # API 文档模板")
			fmt.Println()
		}

		if templateType == "" || templateType == "module" {
			fmt.Println("📦 模块模板 (用于 goon add):")
			fmt.Println()
			fmt.Println("  基础模板 (默认):")
			fmt.Println("    - handler.go.tmpl                 # HTTP 处理器骨架")
			fmt.Println("    - service.go.tmpl                 # 业务逻辑层骨架")
			fmt.Println("    - model.go.tmpl                   # 数据模型骨架")
			fmt.Println("    - repository.go.tmpl              # 数据访问层骨架")
			fmt.Println("    - schema.go.tmpl                  # 请求/响应结构骨架")
			fmt.Println("    - routes.go.tmpl                  # 路由注册")
			fmt.Println()
			fmt.Println("  示例模板 (--example):")
			fmt.Println("    - handler_example.go.tmpl         # 完整 HTTP 处理器实现")
			fmt.Println("    - service_example.go.tmpl         # 完整业务逻辑实现")
			fmt.Println("    - model_example.go.tmpl           # 完整数据模型实现")
			fmt.Println("    - repository_example.go.tmpl      # 完整数据访问实现")
			fmt.Println("    - schema_example.go.tmpl          # 完整请求/响应结构")
			fmt.Println()
		}

		fmt.Println("💡 使用提示:")
		fmt.Println("  - 基础模板：生成代码骨架，包含 TODO 注释")
		fmt.Println("  - 示例模板：生成完整可运行的代码，包含实际字段和业务逻辑")
		fmt.Println()
		fmt.Println("📖 示例:")
		fmt.Println("  goon add user                    # 使用基础模板")
		fmt.Println("  goon add user --example          # 使用示例模板")
		fmt.Println("  goon add user --dry-run          # 预览将生成的文件")
	},
}

func init() {
	rootCmd.AddCommand(templateCmd)
	templateCmd.AddCommand(templateListCmd)
	templateListCmd.Flags().String("type", "", "模板类型 (project/module)")
}
