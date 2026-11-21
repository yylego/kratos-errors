<!-- TEMPLATE (EN) BEGIN: BADGES -->

[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/yylego/kratos-errors/release.yml?branch=main&label=BUILD)](https://github.com/yylego/kratos-errors/actions/workflows/release.yml?query=branch%3Amain)
[![GoDoc](https://pkg.go.dev/badge/github.com/yylego/kratos-errors)](https://pkg.go.dev/github.com/yylego/kratos-errors)
[![Coverage Status](https://img.shields.io/coveralls/github/yylego/kratos-errors/main.svg)](https://coveralls.io/github/yylego/kratos-errors?branch=main)
[![Supported Go Versions](https://img.shields.io/badge/Go-1.26%2B-lightgrey.svg)](https://go.dev/)
[![GitHub Release](https://img.shields.io/github/release/yylego/kratos-errors.svg)](https://github.com/yylego/kratos-errors/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/yylego/kratos-errors)](https://goreportcard.com/report/github.com/yylego/kratos-errors)

<!-- TEMPLATE (EN) CLOSE: BADGES -->

# kratos-errors

Advanced Kratos error handling package with type-safe operations and nil interface trap prevention.

---

<!-- TEMPLATE (EN) BEGIN: LANGUAGE NAVIGATION -->

## CHINESE README

[中文说明](README.zh.md)

<!-- TEMPLATE (EN) CLOSE: LANGUAGE NAVIGATION -->

## Main Features

- 🏗️ **Errors from Reason Enums**: `kratoserrors.NewError` builds a Kratos error straight from a generated reason enum, and `IsError` matches it
- ⚡ **Safe with `(*T)(nil)`**: a nil `*errors.Error` passed as `error` counts as no error in `As`, `IsError` and the test helpers
- 🔄 **Test Helpers**: `erkassert` / `erkrequire` wrap testify, `erkmust` panics on a Kratos error
- 🎯 **No Import Clash**: the `errorskratos.Erk` alias saves importing two packages both named `errors`

## Installation

```bash
go get github.com/yylego/kratos-errors
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

### Concise Error Creation (kratoserrors)

```go
import "github.com/yylego/kratos-errors/kratoserrors"

// Create type-safe error with enum
erk := kratoserrors.NewError(404, ErrorReason_USER_NOT_FOUND, "user %d not found", userID)

// Keep the cause with %w, the same as in fmt.Errorf, so errors.Is and errors.As reach it
erk = kratoserrors.NewError(500, ErrorReason_SERVER_DB_ERROR, "find user %d: %w", userID, err)

// Check error type
if kratoserrors.IsError(err, ErrorReason_USER_NOT_FOUND, 404) {
    // Handle user not found error
}
```

### Reason and the Client Enum

`NewError` puts the enum value name into the fault's reason (e.g. `USER_NOT_FOUND`) and adds no extra field. A client that needs the numeric enum value maps the reason back through its own generated enum: with protobuf-ts, `ErrorReason[reason]` gives `ErrorReason.USER_NOT_FOUND`, so the fault does not need to hold the numeric value in metadata. The type constraint asks the enum to have the two methods that protoc-gen-go generates (`String` and the numeric method), so a plain string type gets rejected at compile time.

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
    expected := errors.InternalServer("SERVER_ERROR", "a different text")
    erkassert.Is(t, expected, erk) // Reason and code match, the message is not compared
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
kratos-errors/
├── errorskratos/       # Error inspection (As, Is, From)
│   ├── errors.go
│   ├── erkadapt/       # Nil interface adaptation
│   └── must/           # Testing and enforcement tools
│       ├── erkassert/  # testify/assert helpers
│       ├── erkrequire/ # testify/require helpers
│       └── erkmust/    # Production panic utilities
├── kratoserrors/       # Typed error creation and check (NewError, IsError)
├── Makefile            # CI: test/coverage/bench
└── internal/           # PB generation via buf
    ├── buf.yaml
    ├── buf.gen.yaml
    ├── Makefile        # generate PB via buf
    ├── errorspb/       # proto + generated pb.go
    └── proto3ps/       # vendored proto imports (errors + google)
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

`erkadapt.Adapt` turns a nil `*errors.Error` into an untyped `nil`, and `As`, `IsError`, `erkassert`, `erkrequire` and `erkmust` treat it as no error.

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

- [kratos-ebz](https://github.com/yylego/kratos-ebz) - The biz error of Kratos: biz methods return `*Ebz`, which does not implement `error`

<!-- TEMPLATE (EN) BEGIN: STANDARD PROJECT FOOTER -->
<!-- VERSION 2025-11-25 03:52:28.131064 +0000 UTC -->

## 📄 License

MIT License - see [LICENSE](LICENSE).

---

## 💬 Contact & Feedback

Contributions are welcome! Report bugs, suggest features, and contribute code:

- 🐛 **Mistake reports?** Open an issue on GitHub with reproduction steps
- 💡 **Fresh ideas?** Create an issue to discuss
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
7. **Documentation**: Update documentation to support client-facing changes
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

<!-- TEMPLATE (EN) CLOSE: STANDARD PROJECT FOOTER -->
