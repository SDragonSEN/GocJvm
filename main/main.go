// main project main.go
package main

import (
	"../class"
	"../memoryControl"
	"../startup"
)

func main() {
	startup.ParseCmd()
	class.InitClassPath(startup.CmdPara.Jar)

	memCtrl.Init(1024)

	//fmt.Println(byte(0x20103004 >> 24))
}
