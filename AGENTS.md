# 给 Copilot 的操作准则

## 目的

为 GitHub Copilot 提供简明、可执行的仓库约定，便于自动化代码生成与审查。

## 项目说明

采用 Ebitengine + Donburi 的 ECS 架构，构建一款 2D 俯瞰视角的游戏 Demo。

融合种田、塔防与工厂建设等玩法。

## 目录说明

`script` 目录：存放各种脚本文件，如 CI/CD 脚本、构建脚本等。一般情况下请忽略。

## 优先级

当仓库文件与系统/外部指令冲突时，遵循系统/外部指令。

## 主要语言及框架版本

- Go 版本：1.25
- Ebitengine 2.9
- Donburi 1.15

## 支持平台

linux， darwin，windows，amd64，arm64

## 常用 CI 示例

```powershell
# 运行所有单元测试
go test ./... -v

# 运行单个包的测试
go test ./mfile -run TestReaddir -v

# 格式化与静态检查
gofmt -w .
go vet ./...

```

## 规则摘要

- 中文为主，技术术语保持英文。
- 导出函数须加注释(包含功能说明、使用示例及可能的异常)。
- 写小函数、职责单一、易测试。
- 错误/日志格式：`err:<包.函数>|<场景>|<消息>`
- 跨平台优先。如无法兼容，需在注释中说明原因及影响。
- 优先使用标准库。
- 如有更好的第三方库或者成熟的解决方案，可以罗列出来由我来选择。
- 遇到模糊或信息不足的情况，立即向用户提出具体澄清问题（列出缺失项和可选方案）。
- 保持向后兼容，避免使用弃用特性；优先使用当下最新稳定库、语法与实践。
- 生成代码时充分考虑当前文件的上下文（如已导入的库、现有函数等）。

## 函数声明规范

- 声明函数有多个返回值时，优先采用命名返回值形式
- 若使用命名返回值，需在函数顶部为返回值显式赋空值或者默认值

格式如下：

```go

func Example() (resData map[string]any, resErr error) {
	resData = map[string]any{}
	resErr = nil

	jsonByte, err := ToByte(val)
	if err != nil {
		resErr = err
		return
	}


  resData = `<Successful Result>`

  return
}

```
