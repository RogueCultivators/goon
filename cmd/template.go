package cmd

import (
	"fmt"
	"strings"

	"github.com/RogueCultivators/goon/internal/template"
	"github.com/RogueCultivators/goon/internal/ui"

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
			ui.Error(fmt.Sprintf("获取参数失败: %v", err))
			return
		}

		projectTmpls, moduleTmpls := template.ListTemplates()

		ui.Header("可用模板列表")

		if templateType == "" || templateType == "project" {
			ui.Step(fmt.Sprintf("项目模板 (用于 goon init) - 共 %d 个", len(projectTmpls)))
			for _, t := range projectTmpls {
				name := strings.TrimSuffix(t, ".tmpl")
				ui.Info(fmt.Sprintf("  - %s", name))
			}
		}

		if templateType == "" || templateType == "module" {
			ui.Step(fmt.Sprintf("模块模板 (用于 goon add) - 共 %d 个", len(moduleTmpls)))
			for _, t := range moduleTmpls {
				name := strings.TrimSuffix(t, ".tmpl")
				if strings.Contains(name, "example") {
					ui.Info(fmt.Sprintf("  - %s (--example)", name))
				} else {
					ui.Info(fmt.Sprintf("  - %s", name))
				}
			}
		}

		ui.Info("")
		ui.Info("使用提示:")
		ui.Info("  - 基础模板：生成代码骨架，包含 TODO 注释")
		ui.Info("  - 示例模板：生成完整可运行的代码")
		ui.Info("")
		ui.Info("示例:")
		ui.Info("  goon add user                    # 使用基础模板")
		ui.Info("  goon add user --example          # 使用示例模板")
		ui.Info("  goon add user --dry-run          # 预览将生成的文件")
	},
}

func init() {
	rootCmd.AddCommand(templateCmd)
	templateCmd.AddCommand(templateListCmd)
	templateListCmd.Flags().String("type", "", "模板类型 (project/module)")
}
