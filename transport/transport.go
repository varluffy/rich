/**
 * @Time: 2021/2/24 11:02 上午
 * @Author: varluffy
 * @Description: transport
 */

package transport

type Server interface {
	Start() error
	Stop() error
}
