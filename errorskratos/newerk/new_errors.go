package newerk

import (
	"fmt"

	"github.com/go-kratos/kratos/v3/errors"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// ProtoErrorEnum constraint for protobuf generated enum types
//
// ProtoErrorEnum protobuf 生成的枚举类型约束
type ProtoErrorEnum interface {
	String() string
	Number() protoreflect.EnumNumber
}

// IsError checks if the given error matches the specified enum and HTTP code
//
// IsError 检查给定错误是否匹配指定的枚举和 HTTP 代码
func IsError[E ProtoErrorEnum](err error, enum E, httpCode int) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == enum.String() && e.Code == int32(httpCode)
}

// NewError creates a new error with the specified parameters
//
// NewError 创建带指定参数的新错误
func NewError[E ProtoErrorEnum](httpCode int, enum E, format string, args ...interface{}) *errors.Error {
	return errors.New(httpCode, enum.String(), fmt.Sprintf(format, args...))
}
