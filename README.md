[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/yylego/kratos-errors/release.yml?branch=main&label=BUILD)](https://github.com/yylego/kratos-errors/actions/workflows/release.yml?query=branch%3Amain)
[![GoDoc](https://pkg.go.dev/badge/github.com/yylego/kratos-errors)](https://pkg.go.dev/github.com/yylego/kratos-errors)
[![Coverage Status](https://img.shields.io/coveralls/github/yylego/kratos-errors/main.svg)](https://coveralls.io/github/yylego/kratos-errors?branch=main)
[![Supported Go Versions](https://img.shields.io/badge/Go-1.25+-lightgrey.svg)](https://go.dev/)
[![GitHub Release](https://img.shields.io/github/release/yylego/kratos-errors.svg)](https://github.com/yylego/kratos-errors/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/yylego/kratos-errors)](https://goreportcard.com/report/github.com/yylego/kratos-errors)

# kratos-errors

Advanced Kratos error handling package with type-safe operations and nil interface trap prevention.

---

<!-- TEMPLATE (EN) BEGIN: LANGUAGE NAVIGATION -->

## CHINESE README

[中文说明](README.zh.md)

<!-- TEMPLATE (EN) END: LANGUAGE NAVIGATION -->

## Main Features

🎯 **Type-Safe Error Handling**: Simplified API to manipulate Kratos errors without naming conflicts
⚡ **Safe Error Handling**: Solves Go's notorious (\*T)(nil) != nil trap through intelligent adaptation
🔄 **Testing Integration**: Complete testify/assert and testify/require helpers to test Kratos errors

## Installation

```bash
go get github.com/yylego/kratos-errors/errorskratos
```

## Usage

### Basic Error Handling

```go
import "github.com/yylego/kratos-errors/errorskratos"

// Type-safe error conversion
err := someFunction()
if erk, ok := errorskratos.As(err); ok {
    fmt.Printf("Kratos error: %s (code: %d)\n", erk.Reason, erk.Code)
}

// Error comparison
erk1 := errors.BadRequest("INVALID_INPUT", "missing field")
erk2 := errors.BadRequest("INVALID_INPUT", "wrong format")
if errorskratos.Is(erk1, erk2) {
    // Same error type (reason and code match)
}

// Convert generic error to Kratos error
erk := errorskratos.From(err)
```

### Concise Error Creation (newerk)

```go
import "github.com/yylego/kratos-errors/errorskratos/newerk"

// Configure reason code field name to store enum numeric value
newerk.SetReasonCodeFieldName("numeric_reason_code_enum")

// Create type-safe error with enum
erk := newerk.NewError(404, ErrorReason_USER_NOT_FOUND, "user %d not found", userID)

// Check error type
if newerk.IsError(err, ErrorReason_USER_NOT_FOUND, 404) {
    // Handle user not found error
}
```

### Testing with Assert

```go
import "github.com/yylego/kratos-errors/errorskratos/must/erkassert"

func TestSomething(t *testing.T) {
    var erk *errors.Error

    // Assert no error (handles nil interface with safe checks)
    erkassert.NoError(t, erk)

    // Assert error exists
    erk = errors.InternalServer("SERVER_ERROR", "database failed")
    erkassert.Error(t, erk)

    // Assert error equivalence
    expected := errors.BadRequest("INVALID_INPUT", "test")
    erkassert.Is(t, expected, erk)
}
```

### Testing with Require

```go
import "github.com/yylego/kratos-errors/errorskratos/must/erkrequire"

func TestCritical(t *testing.T) {
    var erk *errors.Error

    // Require no error (stops test at once if error exists)
    erkrequire.NoError(t, erk)

    // Continue when no error...
}
```

### Production Error Enforcement

```go
import "github.com/yylego/kratos-errors/errorskratos/must/erkmust"

func criticalOperation() {
    erk := doSomethingImportant()

    // Panic if error exists (with structured logging)
    erkmust.Done(erk)

    // Use Must (same function, different name)
    erkmust.Must(erk)
}
```

## Package Structure

```
errorskratos/
├── errors.go           # Core API (As, Is, From)
├── newerk/             # Concise error creation API
├── erkadapt/           # Nil interface adaptation
├── must/               # Testing and enforcement tools
│   ├── erkassert/      # testify/assert helpers
│   ├── erkrequire/     # testify/require helpers
│   └── erkmust/        # Production panic utilities
└── internal/
    └── errorspb/       # Example error definitions
```

## Core Capabilities

### The Nil Interface Issue

Go has a known issue where a typed nil value doesn't match nil when converted to interface:

```go
var erk *errors.Error = nil
var err error = erk
fmt.Println(erk == nil)  // true
fmt.Println(err == nil)  // false (!!)
```

This causes issues in error handling. kratos-errors solves this through intelligent adaptation in each function.

### Clean Naming

The `Erk` type alias avoids import conflicts between standard `errors` package and Kratos `errors`:

```go
// Instead of this confusion:
import (
    stderrors "errors"
    "github.com/go-kratos/kratos/v3/errors"
)

// Just use:
import "github.com/yylego/kratos-errors/errorskratos"
// And work with errorskratos.Erk
```

## More Projects

- [kratos-ebz](https://github.com/yylego/kratos-ebz) - Error type that doesn't implement error interface

<!-- TEMPLATE (EN) BEGIN: STANDARD PROJECT FOOTER -->
<!-- VERSION 2025-09-26 07:39:27.188023 +0000 UTC -->

## 📄 License

MIT License. See [LICENSE](LICENSE).

---

## 🤝 Contributing

Contributions are welcome! Report bugs, suggest features, and contribute code:

- 🐛 **Found a mistake?** Open an issue on GitHub with reproduction steps
- 💡 **Have a feature idea?** Create an issue to discuss the suggestion
- 📖 **Documentation confusing?** Report it so we can improve
- 🚀 **Need new features?** Share the use cases to help us understand requirements
- ⚡ **Performance issue?** Help us optimize through reporting slow operations
- 🔧 **Configuration problem?** Ask questions about complex setups
- 📢 **Follow project progress?** Watch the repo to get new releases and features
- 🌟 **Success stories?** Share how this package improved the workflow
- 💬 **Feedback?** We welcome suggestions and comments

---

## 🔧 Development

New code contributions, follow this process:

1. **Fork**: Fork the repo on GitHub (using the webpage UI).
2. **Clone**: Clone the forked project (`git clone https://github.com/yourname/repo-name.git`).
3. **Navigate**: Navigate to the cloned project (`cd repo-name`)
4. **Branch**: Create a feature branch (`git checkout -b feature/xxx`).
5. **Code**: Implement the changes with comprehensive tests
6. **Testing**: (Golang project) Ensure tests pass (`go test ./...`) and follow Go code style conventions
7. **Documentation**: Update documentation to support client-facing changes and use significant commit messages
8. **Stage**: Stage changes (`git add .`)
9. **Commit**: Commit changes (`git commit -m "Add feature xxx"`) ensuring backward compatible code
10. **Push**: Push to the branch (`git push origin feature/xxx`).
11. **PR**: Open a merge request on GitHub (on the GitHub webpage) with detailed description.

Please ensure tests pass and include relevant documentation updates.

---

## 🌟 Support

Welcome to contribute to this project via submitting merge requests and reporting issues.

**Project Support:**

- ⭐ **Give GitHub stars** if this project helps you
- 🤝 **Share with teammates** and (golang) programming friends
- 📝 **Write tech blogs** about development tools and workflows - we provide content writing support
- 🌟 **Join the ecosystem** - committed to supporting open source and the (golang) development scene

**Have Fun Coding with this package!** 🎉🎉🎉

<!-- TEMPLATE (EN) END: STANDARD PROJECT FOOTER -->

---

## GitHub Stars

[![Stargazers](https://starchart.cc/yylego/kratos-errors.svg?variant=adaptive)](https://starchart.cc/yylego/kratos-errors)
