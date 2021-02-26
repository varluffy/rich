/**
 * @Time: 2021/2/24 3:25 下午
 * @Author: varluffy
 * @Description: //TODO
 */

package log

import "testing"

func TestNewJSONLogger(t *testing.T) {
	logger := NewLogger(WithDebugLevel(), WithFileRotation("/Users/leng/code/go-api/1111.log"), WithConsoleEncoder(), WithDisableConsole())
	defer logger.Sync()
	logger.Debug("debug logger")
	logger.Info("info logger")
	logger.Error("error logger")
}
