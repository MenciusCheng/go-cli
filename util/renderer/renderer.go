package renderer

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

// Renderer 是模板渲染的封装类
type Renderer struct {
	// 自定义模板函数映射
	funcMap template.FuncMap
}

// New 创建一个新的模板渲染器
func New() *Renderer {
	tr := &Renderer{
		funcMap: template.FuncMap{},
	}
	// 注册默认的函数
	tr.registerDefaultFuncs()
	return tr
}

// 注册默认的模板函数
func (tr *Renderer) registerDefaultFuncs() {
	// 时间相关函数
	tr.funcMap["now"] = time.Now
	tr.funcMap["formatTime"] = func(format string) string {
		return time.Now().Format(format)
	}
	tr.funcMap["formatDate"] = func() string {
		return time.Now().Format("2006-01-02")
	}
	tr.funcMap["formatDateTime"] = func() string {
		return time.Now().Format("2006-01-02 15:04:05")
	}

	// 字符串处理函数
	tr.funcMap["upper"] = strings.ToUpper
	tr.funcMap["lower"] = strings.ToLower
	tr.funcMap["title"] = strings.Title
	tr.funcMap["trim"] = strings.TrimSpace
	tr.funcMap["replace"] = strings.Replace
	tr.funcMap["contains"] = strings.Contains
	tr.funcMap["hasPrefix"] = strings.HasPrefix
	tr.funcMap["hasSuffix"] = strings.HasSuffix

	// 格式化函数
	tr.funcMap["sprintf"] = fmt.Sprintf

	// 命名风格转换
	tr.funcMap["camel"] = ToCamelCase       // 小驼峰: thisIsExample
	tr.funcMap["snake"] = ToSnakeCase       // 蛇形: this_is_example
	tr.funcMap["pascal"] = ToPascalCase     // 大驼峰: ThisIsExample
	tr.funcMap["kebab"] = ToKebabCase       // 短横线: this-is-example
	tr.funcMap["constant"] = ToConstantCase // 常量: THIS_IS_EXAMPLE

	// 首字母处理
	tr.funcMap["firstUpper"] = FirstLetterUpper // 首字母大写
	tr.funcMap["firstLower"] = FirstLetterLower // 首字母小写

	// 数组逗号，非最后元素则返回逗号
	tr.funcMap["joinComma"] = JoinComma
}

// AddFunc 添加自定义模板函数
func (tr *Renderer) AddFunc(name string, fn interface{}) *Renderer {
	tr.funcMap[name] = fn
	return tr
}

// RenderString 渲染模板字符串
func (tr *Renderer) RenderString(tmplContent string, data interface{}) (string, error) {
	t, err := template.New("template").Funcs(tr.funcMap).Parse(tmplContent)
	if err != nil {
		return "", fmt.Errorf("解析模板失败: %w", err)
	}

	var buf strings.Builder
	if err := t.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("渲染模板失败: %w", err)
	}

	return buf.String(), nil
}

// RenderToFile 渲染模板并输出到文件
func (tr *Renderer) RenderToFile(tmplContent string, data interface{}, filePath string) error {
	// 确保目录存在
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	// 渲染模板内容
	content, err := tr.RenderString(tmplContent, data)
	if err != nil {
		return err
	}

	// 写入文件
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	return nil
}

// RenderToFileAskIfExist 渲染模板并输出到文件，如果文件存在则先询问是否覆盖
func (tr *Renderer) RenderToFileAskIfExist(tmplContent string, data interface{}, filePath string) error {
	// 检查目标文件是否存在
	if _, err := os.Stat(filePath); err == nil {
		// 文件已存在，询问是否覆盖
		fmt.Printf("文件已存在: %s\n", filePath)
		fmt.Print("是否覆盖? (y/N): ")

		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("读取输入失败: %w", err)
		}

		input = strings.TrimSpace(strings.ToLower(input))
		if input != "y" && input != "yes" {
			fmt.Println("操作已取消")
			return nil
		}
	}

	return tr.RenderToFile(tmplContent, data, filePath)
}

// RenderFromFile 从模板文件渲染内容
func (tr *Renderer) RenderFromFile(tmplFilePath string, data interface{}) (string, error) {
	// 读取模板文件
	tmplContent, err := os.ReadFile(tmplFilePath)
	if err != nil {
		return "", fmt.Errorf("读取模板文件失败: %w", err)
	}

	// 渲染模板
	return tr.RenderString(string(tmplContent), data)
}

// RenderFileToFile 从模板文件渲染并输出到目标文件
func (tr *Renderer) RenderFileToFile(tmplFilePath string, data interface{}, outputFilePath string) error {
	// 读取模板文件
	tmplContent, err := os.ReadFile(tmplFilePath)
	if err != nil {
		return fmt.Errorf("读取模板文件失败: %w", err)
	}

	// 渲染到文件
	return tr.RenderToFile(string(tmplContent), data, outputFilePath)
}

// MustRenderString 安全版本的渲染字符串，出错时返回空字符串
func (tr *Renderer) MustRenderString(tmplContent string, data interface{}) string {
	result, err := tr.RenderString(tmplContent, data)
	if err != nil {
		return ""
	}
	return result
}
