// Package erkassert wraps testify/assert to test Kratos errors, handling the
// nil interface trap via erkadapt.
//
// erkassert 为 Kratos 错误封装 testify/assert，经 erkadapt 处理 nil 接口陷阱。
package erkassert

import (
	"github.com/go-kratos/kratos/v3/errors"
	"github.com/stretchr/testify/assert"
	"github.com/yylego/kratos-errors/errorskratos/erkadapt"
)

// NoError asserts that erk is nil.
//
// NoError 断言 erk 为 nil。
func NoError(t assert.TestingT, erk *errors.Error, msgAndArgs ...any) bool {
	return assert.NoError(t, erkadapt.Adapt(erk), msgAndArgs...)
}

// Error asserts that erk is non-nil.
//
// Error 断言 erk 非 nil。
func Error(t assert.TestingT, erk *errors.Error, msgAndArgs ...any) bool {
	return assert.Error(t, erkadapt.Adapt(erk), msgAndArgs...)
}

// Is asserts that two Kratos errors match on nil state, Reason and Code.
//
// Is 断言两个 Kratos 错误在 nil 状态、Reason 和 Code 上一致。
func Is(t assert.TestingT, expected *errors.Error, erk *errors.Error, msgAndArgs ...any) bool {
	if ok := assert.Equal(t, expected == nil, erk == nil, msgAndArgs...); !ok {
		return false
	}
	if expected != nil && erk != nil {
		if ok := assert.Equal(t, expected.Reason, erk.Reason, msgAndArgs...); !ok {
			return false
		}
		if ok := assert.Equal(t, expected.Code, erk.Code, msgAndArgs...); !ok {
			return false
		}
	}
	return true
}
