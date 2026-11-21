// Package erkrequire wraps testify/require to test Kratos errors, handling the
// nil interface trap via erkadapt. Unlike erkassert, a failed check stops the test.
//
// erkrequire 为 Kratos 错误封装 testify/require，经 erkadapt 处理 nil 接口陷阱；
// 与 erkassert 不同，断言失败会直接终止测试。
package erkrequire

import (
	"github.com/go-kratos/kratos/v3/errors"
	"github.com/stretchr/testify/require"
	"github.com/yylego/kratos-errors/errorskratos/erkadapt"
)

// NoError requires that erk is nil.
//
// NoError 要求 erk 为 nil。
func NoError(t require.TestingT, erk *errors.Error, msgAndArgs ...any) {
	require.NoError(t, erkadapt.Adapt(erk), msgAndArgs...)
}

// Error requires that erk is non-nil.
//
// Error 要求 erk 非 nil。
func Error(t require.TestingT, erk *errors.Error, msgAndArgs ...any) {
	require.Error(t, erkadapt.Adapt(erk), msgAndArgs...)
}

// Is requires that two Kratos errors match on nil state, Reason and Code.
//
// Is 要求两个 Kratos 错误在 nil 状态、Reason 和 Code 上一致。
func Is(t require.TestingT, expected *errors.Error, erk *errors.Error, msgAndArgs ...any) {
	require.Equal(t, expected == nil, erk == nil, msgAndArgs...)
	if expected != nil && erk != nil {
		require.Equal(t, expected.Reason, erk.Reason, msgAndArgs...)
		require.Equal(t, expected.Code, erk.Code, msgAndArgs...)
	}
}
