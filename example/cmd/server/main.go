/**
 * @Time: 2021/3/11 3:32 下午
 * @Author: varluffy
 */

package main

import (
	"fmt"
	"os"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
