/**
 * @Time: 2021/3/8 4:16 下午
 * @Author: varluffy
 */

package biz

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewAuthUsecase)
