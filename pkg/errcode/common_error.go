/**
* @Author:zhoutao
* @Date:2020/7/25 下午9:23
 */

package errcode

var (
	Success          = NewError(0, "成功")
	Fail             = NewError(10000000, "内部错误")
	InvalidParams    = NewError(10000001, "无效参数")
	Unauthorized     = NewError(10000002, "认证错误")
	NotFound         = NewError(10000003, "没有找到")
	Unkonwn          = NewError(1000004, "未知错误")
	DeadlineExceeded = NewError(10000005, "超过最后截止期限")
	AccessDenied     = NewError(10000006, "访问被拒绝")
	LimitExceded     = NewError(10000007, "访问受限")
	MethodNotAllowed = NewError(10000008, "不支持该方法")
)
