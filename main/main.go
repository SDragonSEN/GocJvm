// main project main.go
package main

import (
	"fmt"
	_ "fmt"

	"../class"
	"../memoryControl"
	"../startup"
)

func main() {
	startup.ParseCmd()
	class.InitClassPath(startup.CmdPara.Jar)
	fmt.Println(memCtrl.BytesToUint32([]byte{0x01, 0x02, 0x03, 0x04}))
	memCtrl.InitEx(100, 20)
	memCtrl.Malloc(10, memCtrl.HEADER_NODE)
	memCtrl.Malloc(8, memCtrl.HEADER_NODE)
	memCtrl.LogMem()

	//fmt.Println(byte(0x20103004 >> 24))
}
