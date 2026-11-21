package erkmust_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yylego/kratos-errors/errorskratos"
	"github.com/yylego/kratos-errors/errorskratos/must/erkmust"
	"github.com/yylego/kratos-errors/internal/errorspb"
	"github.com/yylego/must"
)

// TestDone validates Done function panic actions on error conditions
// Tests successful path (no panic) and error path (immediate panic)
// Demonstrates incompatibility with generic must.Done due to nil interface handling
//
// TestDone 验证 Done 函数在错误条件下的 panic 行为
// 测试成功路径（无 panic）和错误路径（立即 panic）
// 演示由于 nil 接口处理导致与通用 must.Done 的不兼容性
func TestDone(t *testing.T) {
	t.Run("no panic", func(t *testing.T) {
		var erk *errorskratos.Erk
		erkmust.Done(erk)
	})

	t.Run("panic on error", func(t *testing.T) {
		require.Panics(t, func() {
			erk := errorspb.ErrorServerDbError("wrong db")
			erkmust.Done(erk)
		})
	})

	t.Run("must done incompatible", func(t *testing.T) {
		require.Panics(t, func() {
			var erk *errorskratos.Erk
			must.Done(erk) // A nil *errors.Error passed as error is not nil, so must.Done panics, use erkmust.Done // nil 的 *errors.Error 当作 error 传进去就不再是 nil，所以 must.Done 会 panic，要用 erkmust.Done
		})
	})
}

// TestMust validates Must function panic actions with structured logging
// Verifies aggressive error enforcement through immediate application crash
// Shows that generic must.Must is incompatible with Kratos error nil handling
//
// TestMust 验证 Must 函数的 panic 行为和结构化日志记录
// 通过立即应用崩溃验证激进的错误强制执行
// 显示通用 must.Must 与 Kratos 错误 nil 处理不兼容
func TestMust(t *testing.T) {
	t.Run("no panic", func(t *testing.T) {
		var erk *errorskratos.Erk
		erkmust.Must(erk)
	})

	t.Run("panic on error", func(t *testing.T) {
		require.Panics(t, func() {
			erk := errorspb.ErrorServerDbTransactionError("wrong tx")
			erkmust.Must(erk)
		})
	})

	t.Run("must must incompatible", func(t *testing.T) {
		require.Panics(t, func() {
			var erk *errorskratos.Erk
			must.Must(erk) // Same as above, use erkmust.Must // 同上，要用 erkmust.Must
		})
	})
}
