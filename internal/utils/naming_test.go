package utils

import (
	"testing"
)

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"空字符串", "", ""},
		{"小写单词", "user", "user"},
		{"大写单词", "USER", "user"},
		{"驼峰命名", "userName", "user_name"},
		{"PascalCase", "UserName", "user_name"},
		{"多个单词", "userProfileData", "user_profile_data"},
		{"连续大写", "HTTPServer", "httpserver"},
		{"短横线", "user-name", "user_name"},
		{"空格", "user name", "user_name"},
		{"混合格式", "user-Profile Name", "user_profile_name"},
		{"已经是snake_case", "user_name", "user_name"},
		{"多个下划线", "user__name", "user_name"},
		{"开头下划线", "_user", "user"},
		{"结尾下划线", "user_", "user"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToSnakeCase(tt.input)
			if result != tt.expected {
				t.Errorf("ToSnakeCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToPascalCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"空字符串", "", ""},
		{"小写单词", "user", "User"},
		{"snake_case", "user_name", "UserName"},
		{"多个单词", "user_profile_data", "UserProfileData"},
		{"驼峰命名", "userName", "UserName"},
		{"短横线", "user-name", "UserName"},
		{"空格", "user name", "UserName"},
		{"混合格式", "user-Profile_name", "UserProfileName"},
		{"已经是PascalCase", "UserName", "UserName"},
		{"多个分隔符", "user__name--data", "UserNameData"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToPascalCase(tt.input)
			if result != tt.expected {
				t.Errorf("ToPascalCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToCamelCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"空字符串", "", ""},
		{"小写单词", "user", "user"},
		{"snake_case", "user_name", "userName"},
		{"多个单词", "user_profile_data", "userProfileData"},
		{"PascalCase", "UserName", "userName"},
		{"短横线", "user-name", "userName"},
		{"空格", "user name", "userName"},
		{"混合格式", "User-Profile_name", "userProfileName"},
		{"已经是camelCase", "userName", "userName"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToCamelCase(tt.input)
			if result != tt.expected {
				t.Errorf("ToCamelCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToKebabCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"空字符串", "", ""},
		{"小写单词", "user", "user"},
		{"snake_case", "user_name", "user-name"},
		{"驼峰命名", "userName", "user-name"},
		{"PascalCase", "UserName", "user-name"},
		{"多个单词", "userProfileData", "user-profile-data"},
		{"空格", "user name", "user-name"},
		{"已经是kebab-case", "user-name", "user-name"},
		{"混合格式", "user_Profile-Name", "user-profile-name"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToKebabCase(tt.input)
			if result != tt.expected {
				t.Errorf("ToKebabCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// 测试实际使用场景
func TestRealWorldScenarios(t *testing.T) {
	scenarios := []struct {
		name     string
		input    string
		snake    string
		pascal   string
		camel    string
		kebab    string
	}{
		{
			name:   "用户资料",
			input:  "userProfile",
			snake:  "user_profile",
			pascal: "UserProfile",
			camel:  "userProfile",
			kebab:  "user-profile",
		},
		{
			name:   "订单历史",
			input:  "order-history",
			snake:  "order_history",
			pascal: "OrderHistory",
			camel:  "orderHistory",
			kebab:  "order-history",
		},
		{
			name:   "产品分类",
			input:  "ProductCategory",
			snake:  "product_category",
			pascal: "ProductCategory",
			camel:  "productCategory",
			kebab:  "product-category",
		},
		{
			name:   "用户管理",
			input:  "user_management",
			snake:  "user_management",
			pascal: "UserManagement",
			camel:  "userManagement",
			kebab:  "user-management",
		},
	}

	for _, sc := range scenarios {
		t.Run(sc.name, func(t *testing.T) {
			if result := ToSnakeCase(sc.input); result != sc.snake {
				t.Errorf("ToSnakeCase(%q) = %q, want %q", sc.input, result, sc.snake)
			}
			if result := ToPascalCase(sc.input); result != sc.pascal {
				t.Errorf("ToPascalCase(%q) = %q, want %q", sc.input, result, sc.pascal)
			}
			if result := ToCamelCase(sc.input); result != sc.camel {
				t.Errorf("ToCamelCase(%q) = %q, want %q", sc.input, result, sc.camel)
			}
			if result := ToKebabCase(sc.input); result != sc.kebab {
				t.Errorf("ToKebabCase(%q) = %q, want %q", sc.input, result, sc.kebab)
			}
		})
	}
}
