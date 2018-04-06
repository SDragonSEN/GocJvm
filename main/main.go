// main project main.go
package main

import (
	"fmt"

	"../accessOp"
	"../class"
	"../classAnaly"
	"../memoryControl"
	"../startup"
)

func main() {
	startup.ParseCmd()
	class.InitClassPath(startup.CmdPara.Jar)

	memCtrl.Init(1024 * 1024)

	_, err := classAnaly.LoadClass("java.lang.Object")
	if err != nil {
		panic("main()1")
	}
	//加载数组
	arrayTypeAdr := classAnaly.LoadArrayClass()
	access.ModifyTypeAddr(access.INIT_ARRAY_CLASS_ADR, arrayTypeAdr)
	access.ArrayClassAdr = arrayTypeAdr
	//加载String
	stringClass, err := classAnaly.LoadClass("java.lang.String")
	if err != nil {
		fmt.Println("error:main():", err)
		panic("main()2")
	}
	access.ModifyTypeAddr(access.INIT_STRING_CLASS_ADR, stringClass.LocalAdr)
	access.StringClassAdr = stringClass.LocalAdr
}
