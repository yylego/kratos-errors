package kratoserrors_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-kratos/kratos/v3/errors"
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

	// %w keeps the cause on the result, the same as in fmt.Errorf
	// %w 把原始错误作为 cause 留在结果上，跟 fmt.Errorf 一样
	t.Run("with cause", func(t *testing.T) {
		_, err := os.Open(filepath.Join(t.TempDir(), "missing.txt"))
		require.Error(t, err)

		erk := kratoserrors.NewError(500, errorspb.ErrorReason_UNKNOWN, "open the config: %w", err)

		require.Equal(t, "open the config: "+err.Error(), erk.Message)
		require.Equal(t, err, erk.Unwrap())
		require.ErrorIs(t, erk, os.ErrNotExist)
	})

	// %v puts the text of the error into the message and keeps no cause
	// %v 只把错误的文字写进消息，不带 cause
	t.Run("without cause", func(t *testing.T) {
		_, err := os.Open(filepath.Join(t.TempDir(), "missing.txt"))
		require.Error(t, err)

		erk := kratoserrors.NewError(500, errorspb.ErrorReason_UNKNOWN, "open the config: %v", err)

		require.Equal(t, "open the config: "+err.Error(), erk.Message)
		require.NoError(t, erk.Unwrap())
		require.NotErrorIs(t, erk, os.ErrNotExist)
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

	t.Run("nil erk passed as error", func(t *testing.T) {
		var erk *errors.Error
		result := kratoserrors.IsError(erk, errorspb.ErrorReason_UNKNOWN, 500)
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
