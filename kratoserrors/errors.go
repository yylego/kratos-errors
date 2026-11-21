package kratoserrors

import (
	stderrors "errors"
	"fmt"

	"github.com/go-kratos/kratos/v3/errors"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// reasonEnum constrains protobuf-generated enum types
//
// reasonEnum protobuf 生成的枚举类型约束
type reasonEnum interface {
	String() string
	Number() protoreflect.EnumNumber
}

// IsError checks if the given error matches the specified enum and HTTP code
//
// IsError 检查给定错误是否匹配指定的枚举和 HTTP 代码
func IsError[E reasonEnum](err error, enum E, httpCode int) bool {
	erk := errors.FromError(err)
	if erk == nil { // A nil err, and also a nil *errors.Error passed as error // err 为 nil，或者是当作 error 传进来的 nil *errors.Error
		return false
	}
	return erk.Reason == enum.String() && erk.Code == int32(httpCode)
}

// NewError creates a new error with the specified parameters
// The format follows fmt.Errorf: a %w verb puts that error on the result as its cause, so errors.Is and errors.As reach it,
// while the message reads the same as with %v. With no %w the result has no cause
//
// NewError 创建带指定参数的新错误
// 格式串的规矩跟 fmt.Errorf 一样：用 %w 时，对应的错误会作为 cause 挂在结果上，errors.Is 和 errors.As 能找到它，
// 消息文字跟用 %v 时相同。不写 %w 时结果不带 cause
func NewError[E reasonEnum](httpCode int, enum E, format string, args ...any) *errors.Error {
	caught := fmt.Errorf(format, args...)
	erk := errors.New(httpCode, enum.String(), caught.Error())
	if cause := stderrors.Unwrap(caught); cause != nil {
		return erk.WithCause(cause)
	}
	return erk
}
