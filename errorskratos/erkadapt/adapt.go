// Package erkadapt converts a Kratos error to the standard error interface
// with correct nil semantics, sidestepping Go's (*T)(nil) != nil trap.
//
// erkadapt 把 Kratos 错误转成标准 error 接口，正确处理 nil，规避 Go 的 (*T)(nil) != nil 陷阱。
package erkadapt

import (
	"github.com/go-kratos/kratos/v3/errors"
)

// Adapt returns erk as an error interface, giving an untyped nil when erk is nil.
//
// A Go interface is a (type, value) tuple; (*errors.Error)(nil) yields a
// non-nil interface holding a nil value, so returning it as-is would
// pollute downstream nil checks. Adapt hands back an untyped nil instead.
// Reference: https://go.dev/doc/faq#nil_error
//
// Adapt 返回 erk 的 error 接口；erk 为 nil 时返回无类型的 nil。
// Go 接口是 (类型, 值) 对，(*errors.Error)(nil) 会得到装着 nil 值的非 nil 接口，
// 直接返回会污染下游的 nil 判断，所以这里返回无类型的 nil。
// 参考: https://go.dev/doc/faq#nil_error
func Adapt(erk *errors.Error) error {
	if erk != nil {
		return erk
	}
	return nil
}
