// Package errorskratos provides type-safe helpers to handle Kratos errors,
// avoiding naming conflicts with the standard errors package.
//
// errorskratos 提供 Kratos 错误的类型安全辅助函数，避免与标准 errors 包命名冲突。
package errorskratos

import (
	"github.com/go-kratos/kratos/v3/errors"
)

// Erk aliases errors.Error to sidestep conflicts with the standard errors package.
//
// Erk 是 errors.Error 的别名，避免与标准 errors 包命名冲突。
type Erk = errors.Error

// As converts a generic error to a Kratos error and reports success.
//
// As 把通用 error 转成 Kratos 错误，并返回是否成功。
func As(err error) (erk *errors.Error, ok bool) {
	if ok = errors.As(err, &erk); ok {
		return erk, true
	}
	return nil, false
}

// Is reports if two Kratos errors match, handling nil on both sides.
//
// Is 判断两个 Kratos 错误是否一致，两侧 nil 都安全处理。
func Is(erk1 *errors.Error, erk2 *errors.Error) bool {
	if erk1 == nil || erk2 == nil {
		return erk1 == nil && erk2 == nil
	}
	return erk1.Is(erk2)
}

// FromError converts a plain error to the Kratos error type.
//
// FromError 把普通 error 转成 Kratos 错误类型。
func FromError(err error) *errors.Error {
	return errors.FromError(err)
}

// From is a short alias of FromError.
//
// From 是 FromError 的简短别名。
func From(err error) *errors.Error {
	return FromError(err)
}
