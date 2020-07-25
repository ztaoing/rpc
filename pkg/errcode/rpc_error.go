/**
* @Author:zhoutao
* @Date:2020/7/25 下午9:30
 */

package errcode

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TogRPCError(err *Error) error {
	s := status.New(ToRPCCode(err.code), err.msg)
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
