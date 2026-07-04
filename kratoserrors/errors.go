package kratoserrors

import (
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
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == enum.String() && e.Code == int32(httpCode)
}

// NewError creates a new error with the specified parameters
//
// NewError 创建带指定参数的新错误
func NewError[E reasonEnum](httpCode int, enum E, format string, args ...interface{}) *errors.Error {
	return errors.New(httpCode, enum.String(), fmt.Sprintf(format, args...))
}
