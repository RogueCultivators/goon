package utils

import (
	"strings"
	"unicode"
)

// ToSnakeCase 将字符串转换为 snake_case
// 支持: camelCase, PascalCase, kebab-case, 空格分隔
func ToSnakeCase(s string) string {
	if s == "" {
		return ""
	}

	// 替换连字符和空格为下划线
	s = strings.ReplaceAll(s, "-", "_")
	s = strings.ReplaceAll(s, " ", "_")

	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			// 如果不是第一个字符，且前一个字符不是下划线，添加下划线
			if i > 0 && s[i-1] != '_' {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}

	// 清理多余的下划线
	return cleanUnderscores(result.String())
}

// ToPascalCase 将字符串转换为 PascalCase
func ToPascalCase(s string) string {
	if s == "" {
		return ""
	}

	// 先转换为 snake_case，然后转换为 PascalCase
	snake := ToSnakeCase(s)
	parts := strings.Split(snake, "_")

	var result strings.Builder
	for _, part := range parts {
		if part == "" {
			continue
		}
		// 首字母大写
		runes := []rune(part)
		runes[0] = unicode.ToUpper(runes[0])
		result.WriteString(string(runes))
	}

	return result.String()
}

// ToCamelCase 将字符串转换为 camelCase
func ToCamelCase(s string) string {
	pascal := ToPascalCase(s)
	if pascal == "" {
		return ""
	}

	// 首字母小写
	runes := []rune(pascal)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

// ToKebabCase 将字符串转换为 kebab-case
func ToKebabCase(s string) string {
	snake := ToSnakeCase(s)
	return strings.ReplaceAll(snake, "_", "-")
}

// cleanUnderscores 清理多余的下划线
func cleanUnderscores(s string) string {
	// 移除开头和结尾的下划线
	s = strings.Trim(s, "_")

	// 将多个连续下划线替换为单个
	for strings.Contains(s, "__") {
		s = strings.ReplaceAll(s, "__", "_")
	}

	return s
}
