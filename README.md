# go-cli

大模型命令行助手

## 安装

### 方式1：源码安装命令工具

```
cd go-cli
go mod tidy
go install
```

### 方式2：[下载已发布版本](https://github.com/MenciusCheng/go-cli/releases)

### 测试安装结果

```bash
go-cli version
```

## ask 命令

咨询大模型代码补全问题。

### 命令行使用

```
go-cli ask "你是谁" --deepseekApiKey "sk-xx"
```

注意：
- 需要替换成实际的 api key `--deepseekApiKey "sk-xx"`
- 或者设置环境变量 `export DEEPSEEK_API_KEY="sk-xx"`， 就不需要每次带上key了。

### idea 配置

1. 打开设置，选择 **Tools | External Tools**
2. 添加新的外部工具
3. 配置工具参数：

基本设置：
- Name: `go-cli 代码咨询`

程序配置：
- Program: `go-cli`
- Arguments: `ask "$Prompt$" --fileDir "$FileDir$" --filePath "$FilePath$" --selectionStartLine "$SelectionStartLine$" --selectionEndLine "$SelectionEndLine$"  --selectionEndLine "$SelectionEndLine$" --selectionStartColumn "$SelectionStartColumn$" --selectionEndColumn "$SelectionEndColumn$" --deepseekApiKey "sk-xx"`
- Working directory: `$FileDir$`

注意：
- 需要替换成实际的 api key `--deepseekApiKey "sk-xx"`
- 选中代码时，默认会当做选中完整行的代码，而不只是行中局部代码。

### idea 使用

1. 选择代码，右键菜单
2. 选择 **External Tools | go-cli 代码咨询**
3. 输入问题


## code 命令

自定义大模型代码补全策略。

### idea 配置

1. 打开设置，选择 **Tools | External Tools**
2. 添加新的外部工具
3. 配置工具参数：

基本设置：
- Name: `go-cli 代码补全`

程序配置：
- Program: `go-cli`
- Arguments: `code "$Prompt$" --fileDir "$FileDir$" --filePath "$FilePath$" --selectionStartLine "$SelectionStartLine$" --selectionEndLine "$SelectionEndLine$"  --selectionEndLine "$SelectionEndLine$" --selectionStartColumn "$SelectionStartColumn$" --selectionEndColumn "$SelectionEndColumn$" --deepseekApiKey "sk-xx"`
- Working directory: `$FileDir$`

注意：
- 需要替换成实际的 api key `--deepseekApiKey "sk-xx"`
- 选中代码时，默认会当做选中完整行的代码，而不只是行中局部代码。

### idea 使用

1. 选择代码，右键菜单
2. 选择 **External Tools | go-cli 代码补全**
3. 输入补全要求

## add 命令

添加新命令到项目中

语法：`go-cli add [name]`

## 构建步骤

比如构建版本 `v1.0.0`

1. 确保代码已提交
2. 创建标签 `git tag -a v1.0.0 -m "发布版本 v1.0.0"`
3. 推送标签 `git push origin v1.0.0`
4. 使用 Makefile 构建 `make all`
