package kratoserrors_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yylego/kratos-errors/internal/errorspb"
	"github.com/yylego/kratos-errors/kratoserrors"
)

// TestNewError tests creating errors with enum types
//
// TestNewError 测试使用枚举类型创建错误
func TestNewError(t *testing.T) {
	t.Run("basic error", func(t *testing.T) {
		erk := kratoserrors.NewError(404, errorspb.ErrorReason_UNKNOWN, "test error")

		require.NotNil(t, erk)
		require.Equal(t, int32(404), erk.Code)
		require.Equal(t, "UNKNOWN", erk.Reason)
		require.Equal(t, "test error", erk.Message)
	})

	t.Run("database error", func(t *testing.T) {
		erk := kratoserrors.NewError(500, errorspb.ErrorReason_SERVER_DB_ERROR, "database failed")

		require.NotNil(t, erk)
		require.Equal(t, int32(500), erk.Code)
		require.Equal(t, "SERVER_DB_ERROR", erk.Reason)
		require.Equal(t, "database failed", erk.Message)
	})

	t.Run("with format args", func(t *testing.T) {
		erk := kratoserrors.NewError(404, errorspb.ErrorReason_UNKNOWN, "user %d not found in %s", 123, "database")

		require.NotNil(t, erk)
		require.Equal(t, int32(404), erk.Code)
		require.Equal(t, "user 123 not found in database", erk.Message)
	})
}

// TestIsError tests checking if error matches enum and HTTP code
//
// TestIsError 测试检查错误是否匹配枚举和 HTTP 代码
func TestIsError(t *testing.T) {
	t.Run("nil error", func(t *testing.T) {
		result := kratoserrors.IsError(nil, errorspb.ErrorReason_UNKNOWN, 500)
		require.False(t, result)
	})

	t.Run("match success", func(t *testing.T) {
		erk := kratoserrors.NewError(500, errorspb.ErrorReason_SERVER_DB_ERROR, "test")
		result := kratoserrors.IsError(erk, errorspb.ErrorReason_SERVER_DB_ERROR, 500)
		require.True(t, result)
	})

	t.Run("wrong reason", func(t *testing.T) {
		erk := kratoserrors.NewError(500, errorspb.ErrorReason_SERVER_DB_ERROR, "test")
		result := kratoserrors.IsError(erk, errorspb.ErrorReason_UNKNOWN, 500)
		require.False(t, result)
	})

	t.Run("wrong code", func(t *testing.T) {
		erk := kratoserrors.NewError(500, errorspb.ErrorReason_SERVER_DB_ERROR, "test")
		result := kratoserrors.IsError(erk, errorspb.ErrorReason_SERVER_DB_ERROR, 404)
		require.False(t, result)
	})

	t.Run("both wrong", func(t *testing.T) {
		erk := kratoserrors.NewError(500, errorspb.ErrorReason_SERVER_DB_ERROR, "test")
		result := kratoserrors.IsError(erk, errorspb.ErrorReason_UNKNOWN, 404)
		require.False(t, result)
	})
}
