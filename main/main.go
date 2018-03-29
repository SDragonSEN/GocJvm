// main project main.go
package main

import (
	"../class"
	"../classAnaly"
	"../memoryControl"
	"../startup"
)

func main() {
	startup.ParseCmd()
	class.InitClassPath(startup.CmdPara.Jar)

	memCtrl.Init(1024)

}
