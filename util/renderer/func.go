package renderer

import (
	"reflect"
	"strings"
	"unicode"
)

// 首字母大写
func FirstLetterUpper(s string) string {
	if s == "" {
		return ""
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

// 首字母小写
func FirstLetterLower(s string) string {
	if s == "" {
		return ""
	}
	r := []rune(s)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}

// 将字符串转换为驼峰式命名（小驼峰）
// 例如: "hello_world" -> "helloWorld"
func ToCamelCase(s string) string {
	// 处理各种分隔符
	s = strings.NewReplacer("-", " ", "_", " ").Replace(s)
	parts := strings.Fields(s)

	// 处理单词
	for i := 1; i < len(parts); i++ {
		if parts[i] != "" {
			parts[i] = FirstLetterUpper(parts[i])
		}
	}

	// 处理第一个单词
	if len(parts) > 0 && parts[0] != "" {
		parts[0] = FirstLetterLower(parts[0])
	}

	return strings.Join(parts, "")
}

// 将字符串转换为帕斯卡命名（大驼峰）
// 例如: "hello_world" -> "HelloWorld"
func ToPascalCase(s string) string {
	// 处理各种分隔符
	s = strings.NewReplacer("-", " ", "_", " ").Replace(s)
	parts := strings.Fields(s)

	// 处理每个单词
	for i := 0; i < len(parts); i++ {
		if parts[i] != "" {
			parts[i] = FirstLetterUpper(parts[i])
		}
	}

	return strings.Join(parts, "")
}

// 将字符串转换为蛇形命名法
// 例如: "HelloWorld" -> "hello_world"
func ToSnakeCase(s string) string {
	var result strings.Builder
	runes := []rune(s)

	for i, r := range runes {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else if unicode.IsDigit(r) && i > 0 {
			// 检查前一个字符是否不是数字，如果不是则添加下划线
			prevRune := runes[i-1]
			if !unicode.IsDigit(prevRune) {
				result.WriteRune('_')
			}
			result.WriteRune(r)
		} else {
			result.WriteRune(r)
		}
	}

	// 将空格和短横线替换为下划线
	return strings.NewReplacer(" ", "_", "-", "_").Replace(result.String())
}

// 将字符串转换为短横线分隔命名法
// 例如: "HelloWorld" -> "hello-world"
func ToKebabCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteRune('-')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}

	// 将空格和下划线替换为短横线
	return strings.NewReplacer(" ", "-", "_", "-").Replace(result.String())
}

// 将字符串转换为常量命名法（全大写加下划线）
// 例如: "HelloWorld" -> "HELLO_WORLD"
func ToConstantCase(s string) string {
	return strings.ToUpper(ToSnakeCase(s))
}

func JoinComma(index int, arr interface{}) string {
	// 使用反射获取数组/切片长度
	v := reflect.ValueOf(arr)
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return ""
	}

	// 如果不是最后一个元素，返回逗号
	if index < v.Len()-1 {
		return ","
	}
	return ""
}
