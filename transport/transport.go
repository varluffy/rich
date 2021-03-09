/**
 * @Time: 2021/2/24 11:02 上午
 * @Author: varluffy
 */

package transport

type Server interface {
	Start() error
	Stop() error
}
