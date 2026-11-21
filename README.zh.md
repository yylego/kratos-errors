<!-- TEMPLATE (ZH) BEGIN: BADGES -->

[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/yylego/kratos-errors/release.yml?branch=main&label=BUILD)](https://github.com/yylego/kratos-errors/actions/workflows/release.yml?query=branch%3Amain)
[![GoDoc](https://pkg.go.dev/badge/github.com/yylego/kratos-errors)](https://pkg.go.dev/github.com/yylego/kratos-errors)
[![Coverage Status](https://img.shields.io/coveralls/github/yylego/kratos-errors/main.svg)](https://coveralls.io/github/yylego/kratos-errors?branch=main)
[![Supported Go Versions](https://img.shields.io/badge/Go-1.26%2B-lightgrey.svg)](https://go.dev/)
[![GitHub Release](https://img.shields.io/github/release/yylego/kratos-errors.svg)](https://github.com/yylego/kratos-errors/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/yylego/kratos-errors)](https://goreportcard.com/report/github.com/yylego/kratos-errors)

<!-- TEMPLATE (ZH) CLOSE: BADGES -->

# kratos-errors

高级 Kratos 错误处理包，提供类型安全操作和 nil 接口陷阱防护。

---

<!-- TEMPLATE (ZH) BEGIN: LANGUAGE NAVIGATION -->

## 英文文档

[ENGLISH README](README.md)

<!-- TEMPLATE (ZH) CLOSE: LANGUAGE NAVIGATION -->

## 核心特性

- 🏗️ **按 reason 枚举造错误**：`kratoserrors.NewError` 直接用生成的 reason 枚举造 Kratos 错误，`IsError` 用同一个枚举判断
- ⚡ **不怕 `(*T)(nil)`**：nil 的 `*errors.Error` 当作 `error` 传进来，`As`、`IsError` 和测试助手都当作没有错误
- 🔄 **测试助手**：`erkassert` / `erkrequire` 封装 testify，`erkmust` 遇到 Kratos 错误就 panic
- 🎯 **不撞包名**：用 `errorskratos.Erk` 别名，就不用同时导入两个都叫 `errors` 的包

## 安装

```bash
go get github.com/yylego/kratos-errors
```

## 使用方法

### 基础错误处理

```go
import "github.com/yylego/kratos-errors/errorskratos"

// 类型安全的错误转换
err := someFunction()
if erk, ok := errorskratos.As(err); ok {
    fmt.Printf("Kratos 错误: %s (代码: %d)\n", erk.Reason, erk.Code)
}

// 错误比较
erk1 := errors.BadRequest("INVALID_INPUT", "缺少字段")
erk2 := errors.BadRequest("INVALID_INPUT", "格式错误")
if errorskratos.Is(erk1, erk2) {
    // 相同的错误类型（原因和代码匹配）
}

// 将任何错误转换为 Kratos 错误
erk := errorskratos.From(err)
```

### 简洁错误创建 (kratoserrors)

```go
import "github.com/yylego/kratos-errors/kratoserrors"

// 使用枚举创建类型安全的错误
erk := kratoserrors.NewError(404, ErrorReason_USER_NOT_FOUND, "用户 %d 未找到", userID)

// 用 %w 把原始错误带上，跟 fmt.Errorf 一样，errors.Is 和 errors.As 能找到它
erk = kratoserrors.NewError(500, ErrorReason_SERVER_DB_ERROR, "查询用户 %d: %w", userID, err)

// 检查错误类型
if kratoserrors.IsError(err, ErrorReason_USER_NOT_FOUND, 404) {
    // 处理用户未找到错误
}
```

### reason 与前端枚举的对应

`NewError` 把枚举值的名字放进错误的 reason（如 `USER_NOT_FOUND`），不额外加任何字段。客户端要用枚举的数字值时，用它自己生成的枚举从 reason 反查即可：protobuf-ts 里 `ErrorReason[reason]` 就得到 `ErrorReason.USER_NOT_FOUND`，所以错误里不需要再把数字值放进 metadata。类型约束要求枚举同时有 protoc-gen-go 生成的两个方法（`String` 和取数字值的方法），这样普通的字符串类型在编译时就会被拒绝。

### 使用 Assert 测试

```go
import "github.com/yylego/kratos-errors/errorskratos/must/erkassert"

func TestSomething(t *testing.T) {
    var erk *errors.Error

    // 断言没有错误（正确处理 nil 接口）
    erkassert.NoError(t, erk)

    // 断言错误存在
    erk = errors.InternalServer("SERVER_ERROR", "数据库失败")
    erkassert.Error(t, erk)

    // 断言错误相等
    expected := errors.InternalServer("SERVER_ERROR", "另一条消息")
    erkassert.Is(t, expected, erk) // reason 和 code 一致即可，不比较消息
}
```

### 使用 Require 测试

```go
import "github.com/yylego/kratos-errors/errorskratos/must/erkrequire"

func TestCritical(t *testing.T) {
    var erk *errors.Error

    // 要求没有错误（如果存在错误立即失败）
    erkrequire.NoError(t, erk)

    // 只有在没有错误时继续...
}
```

### 生产环境错误强制执行

```go
import "github.com/yylego/kratos-errors/errorskratos/must/erkmust"

func criticalOperation() {
    erk := doSomethingImportant()

    // 如果存在错误则 panic（带结构化日志）
    erkmust.Done(erk)

    // 或使用 Must（相同行为，不同名称）
    erkmust.Must(erk)
}
```

## 包结构

```
kratos-errors/
├── errorskratos/       # 错误检查 (As, Is, From)
│   ├── errors.go
│   ├── erkadapt/       # Nil 接口适配
│   └── must/           # 测试和强制执行工具
│       ├── erkassert/  # testify/assert 包装器
│       ├── erkrequire/ # testify/require 包装器
│       └── erkmust/    # 生产环境 panic 工具
├── kratoserrors/       # 类型安全的错误创建与判断 (NewError, IsError)
├── Makefile            # CI：测试/覆盖率/基准
└── internal/           # 用 buf 生成 PB（仅内部）
    ├── buf.yaml
    ├── buf.gen.yaml
    ├── Makefile        # 用 buf 生成 PB
    ├── errorspb/       # proto + 生成的 pb.go
    └── proto3ps/       # vendor 的 proto 依赖 (errors + google)
```

## 核心能力

### Nil 接口问题

Go 有一个众所周知的问题，当类型化的 nil 指针转换为接口时不等于 nil：

```go
var erk *errors.Error = nil
var err error = erk
fmt.Println(erk == nil)  // true
fmt.Println(err == nil)  // false (!!)
```

`erkadapt.Adapt` 把 nil 的 `*errors.Error` 转成无类型的 `nil`；`As`、`IsError`、`erkassert`、`erkrequire`、`erkmust` 都把它当作没有错误。

### 清晰的命名

`Erk` 类型别名避免了标准 `errors` 包和 Kratos `errors` 之间的导入冲突：

```go
// 不用这种混乱的方式：
import (
    stderrors "errors"
    "github.com/go-kratos/kratos/v3/errors"
)

// 只需使用：
import "github.com/yylego/kratos-errors/errorskratos"
// 然后使用 errorskratos.Erk
```

## 相关项目

- [kratos-ebz](https://github.com/yylego/kratos-ebz) - Kratos 的业务错误：biz 方法返回 `*Ebz`，它不实现 `error` 接口

<!-- TEMPLATE (ZH) BEGIN: STANDARD PROJECT FOOTER -->
<!-- VERSION 2025-11-25 03:52:28.131064 +0000 UTC -->

## 📄 许可证类型

MIT 许可证 - 详见 [LICENSE](LICENSE)。

---

## 💬 联系与反馈

非常欢迎贡献代码！报告 BUG、建议功能、贡献代码：

- 🐛 **问题报告？** 在 GitHub 上提交问题并附上重现步骤
- 💡 **新颖思路？** 创建 issue 讨论
- 📖 **文档疑惑？** 报告问题，帮助我们完善文档
- 🚀 **需要功能？** 分享使用场景，帮助理解需求
- ⚡ **性能瓶颈？** 报告慢操作，协助解决性能问题
- 🔧 **配置困扰？** 询问复杂设置的相关问题
- 📢 **关注进展？** 关注仓库以获取新版本和功能
- 🌟 **成功案例？** 分享这个包如何改善工作流程
- 💬 **反馈意见？** 欢迎提出建议和意见

---

## 🔧 代码贡献

新代码贡献，请遵循此流程：

1. **Fork**：在 GitHub 上 Fork 仓库（使用网页界面）
2. **克隆**：克隆 Fork 的项目（`git clone https://github.com/yourname/repo-name.git`）
3. **导航**：进入克隆的项目（`cd repo-name`）
4. **分支**：创建功能分支（`git checkout -b feature/xxx`）
5. **编码**：实现您的更改并编写全面的测试
6. **测试**：（Golang 项目）确保测试通过（`go test ./...`）并遵循 Go 代码风格约定
7. **文档**：面向用户的更改需要更新文档
8. **暂存**：暂存更改（`git add .`）
9. **提交**：提交更改（`git commit -m "Add feature xxx"`）确保向后兼容的代码
10. **推送**：推送到分支（`git push origin feature/xxx`）
11. **PR**：在 GitHub 上打开 Merge Request（在 GitHub 网页上）并提供详细描述

请确保测试通过并包含相关的文档更新。

---

## 🌟 项目支持

非常欢迎通过提交 Merge Request 和报告问题来贡献此项目。

**项目支持：**

- ⭐ **给予星标**如果项目对您有帮助
- 🤝 **分享项目**给团队成员和（golang）编程朋友
- 📝 **撰写博客**关于开发工具和工作流程 - 我们提供写作支持
- 🌟 **加入生态** - 致力于支持开源和（golang）开发场景

**祝你用这个包编程愉快！** 🎉🎉🎉

<!-- TEMPLATE (ZH) CLOSE: STANDARD PROJECT FOOTER -->
