package errorspb_test

import (
	"testing"

	"github.com/go-kratos/kratos/v3/errors"
	"github.com/stretchr/testify/require"
	errorspb2 "github.com/yylego/kratos-errors/internal/errorspb"
	"github.com/yylego/neatjson/neatjsons"
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestErrorUnknown(t *testing.T) {
	erk := errorspb2.ErrorUnknown("internal system failure: %s", "database connection lost")
	t.Log(erk)
	t.Log(neatjsons.S(erk))

	require.NotNil(t, erk)
	require.Equal(t, int32(500), erk.Code)
	require.Equal(t, errorspb2.ErrorReason_UNKNOWN.String(), erk.Reason)
	require.Equal(t, "internal system failure: database connection lost", erk.Message)
}

func TestErrorServerDbError(t *testing.T) {
	erk := errorspb2.ErrorServerDbError("database query failed: %s", "connection timeout")
	t.Log(erk)
	t.Log(neatjsons.S(erk))

	require.NotNil(t, erk)
	require.Equal(t, int32(500), erk.Code)
	require.Equal(t, errorspb2.ErrorReason_SERVER_DB_ERROR.String(), erk.Reason)
	require.Equal(t, "database query failed: connection timeout", erk.Message)
}

func TestErrorServerDbTransactionError(t *testing.T) {
	erk := errorspb2.ErrorServerDbTransactionError("transaction failed: %s", "deadlock detected")
	t.Log(erk)
	t.Log(neatjsons.S(erk))

	require.NotNil(t, erk)
	require.Equal(t, int32(500), erk.Code)
	require.Equal(t, errorspb2.ErrorReason_SERVER_DB_TRANSACTION_ERROR.String(), erk.Reason)
	require.Equal(t, "transaction failed: deadlock detected", erk.Message)
}

func TestIsUnknown(t *testing.T) {
	t.Run("match unknown", func(t *testing.T) {
		erk := errorspb2.ErrorUnknown("test unknown error")
		t.Log(erk)
		t.Log(neatjsons.S(erk))
		require.True(t, errorspb2.IsUnknown(erk))
	})

	t.Run("not match db", func(t *testing.T) {
		erk := errorspb2.ErrorServerDbError("test db error")
		t.Log(erk)
		t.Log(neatjsons.S(erk))
		require.False(t, errorspb2.IsUnknown(erk))
	})

	t.Run("match standard", func(t *testing.T) {
		erk := errors.New(500, "UNKNOWN", "standard unknown error")
		t.Log(erk)
		t.Log(neatjsons.S(erk))
		require.True(t, errorspb2.IsUnknown(erk))
	})

	t.Run("nil check", func(t *testing.T) {
		require.False(t, errorspb2.IsUnknown(nil))
	})
}

func TestIsServerDbError(t *testing.T) {
	t.Run("match db", func(t *testing.T) {
		erk := errorspb2.ErrorServerDbError("test db error")
		t.Log(erk)
		t.Log(neatjsons.S(erk))
		require.True(t, errorspb2.IsServerDbError(erk))
	})

	t.Run("not match unknown", func(t *testing.T) {
		erk := errorspb2.ErrorUnknown("test unknown")
		t.Log(erk)
		t.Log(neatjsons.S(erk))
		require.False(t, errorspb2.IsServerDbError(erk))
	})

	t.Run("match standard", func(t *testing.T) {
		erk := errors.New(500, "SERVER_DB_ERROR", "standard db error")
		t.Log(erk)
		t.Log(neatjsons.S(erk))
		require.True(t, errorspb2.IsServerDbError(erk))
	})
}

func TestIsServerDbTransactionError(t *testing.T) {
	t.Run("match transaction", func(t *testing.T) {
		erk := errorspb2.ErrorServerDbTransactionError("test transaction error")
		t.Log(erk)
		t.Log(neatjsons.S(erk))
		require.True(t, errorspb2.IsServerDbTransactionError(erk))
	})

	t.Run("not match db", func(t *testing.T) {
		erk := errorspb2.ErrorServerDbError("test db error")
		t.Log(erk)
		t.Log(neatjsons.S(erk))
		require.False(t, errorspb2.IsServerDbTransactionError(erk))
	})

	t.Run("match standard", func(t *testing.T) {
		erk := errors.New(500, "SERVER_DB_TRANSACTION_ERROR", "standard transaction error")
		t.Log(erk)
		t.Log(neatjsons.S(erk))
		require.True(t, errorspb2.IsServerDbTransactionError(erk))
	})
}

func TestErrorReasonStringConversion(t *testing.T) {
	// Test string representation
	require.Equal(t, "UNKNOWN", errorspb2.ErrorReason_UNKNOWN.String())
	require.Equal(t, "SERVER_DB_ERROR", errorspb2.ErrorReason_SERVER_DB_ERROR.String())
	require.Equal(t, "SERVER_DB_TRANSACTION_ERROR", errorspb2.ErrorReason_SERVER_DB_TRANSACTION_ERROR.String())
}

func TestAllErrorTypesHaveCorrectHttpCode(t *testing.T) {
	t.Run("unknown code", func(t *testing.T) {
		erk := errorspb2.ErrorUnknown("test")
		t.Log(erk)
		t.Log(neatjsons.S(erk))
		require.Equal(t, int32(500), erk.Code)
	})

	t.Run("db code", func(t *testing.T) {
		erk := errorspb2.ErrorServerDbError("test")
		t.Log(erk)
		t.Log(neatjsons.S(erk))
		require.Equal(t, int32(500), erk.Code)
	})

	t.Run("transaction code", func(t *testing.T) {
		erk := errorspb2.ErrorServerDbTransactionError("test")
		t.Log(erk)
		t.Log(neatjsons.S(erk))
		require.Equal(t, int32(500), erk.Code)
	})
}

func TestCrossErrorTypeChecking(t *testing.T) {
	t.Run("db check", func(t *testing.T) {
		erk := errorspb2.ErrorServerDbError("database connection failed")
		t.Log(erk)
		t.Log(neatjsons.S(erk))

		require.True(t, errorspb2.IsServerDbError(erk))
		require.False(t, errorspb2.IsUnknown(erk))
		require.False(t, errorspb2.IsServerDbTransactionError(erk))
	})

	t.Run("transaction check", func(t *testing.T) {
		erk := errorspb2.ErrorServerDbTransactionError("transaction rollback failed")
		t.Log(erk)
		t.Log(neatjsons.S(erk))

		require.True(t, errorspb2.IsServerDbTransactionError(erk))
		require.False(t, errorspb2.IsUnknown(erk))
		require.False(t, errorspb2.IsServerDbError(erk))
	})
}
