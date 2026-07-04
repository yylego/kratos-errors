// Package erkmust panics on a non-nil Kratos error, logging it first.
// Use it on operations that must succeed.
//
// erkmust 在 Kratos 错误非 nil 时 panic 并先记录日志；用于必须成功的操作。
package erkmust

import (
	"github.com/yylego/kratos-errors/errorskratos"
	"github.com/yylego/zaplog"
	"go.uber.org/zap"
)

// Done panics when erk is non-nil, logging the error first.
//
// Done 在 erk 非 nil 时 panic，先记录错误。
func Done(erk *errorskratos.Erk) {
	if erk != nil {
		zaplog.ZAPS.Skip1.LOG.Panic("UNEXPECTED ERROR", zap.Error(erk))
	}
}

// Must is an alias of Done.
//
// Must 是 Done 的别名。
func Must(erk *errorskratos.Erk) {
	if erk != nil {
		zaplog.ZAPS.Skip1.LOG.Panic("UNEXPECTED ERROR", zap.Error(erk))
	}
}
