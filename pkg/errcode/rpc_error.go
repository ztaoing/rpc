/**
* @Author:zhoutao
* @Date:2020/7/25 下午9:30
 */

package errcode

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"rpc/proto"
)

/**
grpc 的状态消息中一共包含三个属性，分别是错误码code、错误信息message、错误详情details
*/
func TogRPCError(err *Error) error {
	s, _ := status.New(ToRPCCode(
		err.Code()),
		err.Msg(),
	).WithDetails(&proto.Error{Code: int32(err.Code()), Message: err.Msg()})
	return s.Err()
}

//把原始业务错误码放入details中
func ToRPCStatus(code int, msg string) *Status {
	s, _ := status.New(ToRPCCode(code), msg).WithDetails(&proto.Error{Code: int32(code), Message: msg})
	return &Status{s}
}

func ToRPCCode(code int) codes.Code {
	var statusCode codes.Code

	switch code {
	case Fail.code:
		statusCode = codes.Internal
	case InvalidParams.code:
		statusCode = codes.InvalidArgument
	case Unauthorized.code:
		statusCode = codes.Unauthenticated
	case AccessDenied.code:
		statusCode = codes.PermissionDenied
	case DeadlineExceeded.code:
		statusCode = codes.DeadlineExceeded
	case NotFound.code:
		statusCode = codes.NotFound
	case LimitExceded.code:
		statusCode = codes.ResourceExhausted
	case MethodNotAllowed.code:
		statusCode = codes.Unimplemented
	default:
		statusCode = codes.Unknown
	}
	return statusCode
}

type Status struct {
	*status.Status
}

//获取错误的类型
func FromError(err error) *Status {
	s, _ := status.FromError(err)
	return &Status{s}
}
