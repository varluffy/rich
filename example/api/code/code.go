/**
 * @Time: 2021/2/27 7:26 下午
 * @Author: varluffy
 * @Description: //TODO
 */

package code

import "github.com/varluffy/ginx/errcode"

var (
	ErrUnauthorizedInvalid = errcode.Unauthorized(401, "签名错误")
	ErrUnauthorizedExpired = errcode.Unauthorized(401, "登录过期")
	ErrUnauthorizedError   = errcode.Unauthorized(401, "签名异常")
)
