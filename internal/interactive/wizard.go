package interactive

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

// ProjectConfig 项目配置
type ProjectConfig struct {
	ProjectName   string
	ModuleName    string
	Database      string
	UseAuth       bool
	AuthMethod    string
	UseDocker     bool
	UseExample    bool
	ExampleModule string
}

// RunInitWizard 运行项目初始化向导
func RunInitWizard() (*ProjectConfig, error) {
	config := &ProjectConfig{}

	// 1. 项目名称
	projectNamePrompt := &survey.Input{
		Message: "项目名称:",
		Help:    "将作为项目目录名称",
	}
	if err := survey.AskOne(projectNamePrompt, &config.ProjectName, survey.WithValidator(survey.Required)); err != nil {
		return nil, err
	}

	// 2. Go module 名称
	moduleNamePrompt := &survey.Input{
		Message: "Go module 名称:",
		Default: config.ProjectName,
		Help:    "例如: github.com/username/project",
	}
	if err := survey.AskOne(moduleNamePrompt, &config.ModuleName); err != nil {
		return nil, err
	}

	// 3. 选择数据库
	databasePrompt := &survey.Select{
		Message: "选择数据库:",
		Options: []string{"PostgreSQL", "MySQL", "SQLite", "无数据库"},
		Default: "PostgreSQL",
	}
	if err := survey.AskOne(databasePrompt, &config.Database); err != nil {
		return nil, err
	}

	// 4. 是否需要认证功能
	authPrompt := &survey.Confirm{
		Message: "是否需要认证功能?",
		Default: true,
	}
	if err := survey.AskOne(authPrompt, &config.UseAuth); err != nil {
		return nil, err
	}

	// 5. 如果需要认证，选择认证方式
	if config.UseAuth {
		authMethodPrompt := &survey.Select{
			Message: "选择认证方式:",
			Options: []string{"JWT", "Session"},
			Default: "JWT",
		}
		if err := survey.AskOne(authMethodPrompt, &config.AuthMethod); err != nil {
			return nil, err
		}
	}

	// 6. 是否需要 Docker 支持
	dockerPrompt := &survey.Confirm{
		Message: "是否需要 Docker 支持?",
		Default: true,
	}
	if err := survey.AskOne(dockerPrompt, &config.UseDocker); err != nil {
		return nil, err
	}

	// 7. 是否生成示例模块
	examplePrompt := &survey.Confirm{
		Message: "是否生成示例模块?",
		Default: true,
		Help:    "生成包含完整实现的示例代码（推荐新手使用）",
	}
	if err := survey.AskOne(examplePrompt, &config.UseExample); err != nil {
		return nil, err
	}

	// 8. 如果生成示例，询问模块名称
	if config.UseExample {
		exampleModulePrompt := &survey.Input{
			Message: "示例模块名称:",
			Default: "user",
			Help:    "将生成完整的 CRUD API",
		}
		if err := survey.AskOne(exampleModulePrompt, &config.ExampleModule); err != nil {
			return nil, err
		}
	}

	return config, nil
}

// PrintSummary 打印配置摘要
func PrintSummary(config *ProjectConfig) {
	fmt.Println("\n📋 项目配置摘要:")
	fmt.Printf("  项目名称: %s\n", config.ProjectName)
	fmt.Printf("  Module 名称: %s\n", config.ModuleName)
	fmt.Printf("  数据库: %s\n", config.Database)
	if config.UseAuth {
		fmt.Printf("  认证方式: %s\n", config.AuthMethod)
	} else {
		fmt.Println("  认证: 不使用")
	}
	fmt.Printf("  Docker 支持: %v\n", config.UseDocker)
	if config.UseExample {
		fmt.Printf("  示例模块: %s\n", config.ExampleModule)
	}
	fmt.Println()
}
