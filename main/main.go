// main project main.go
package main

import (
	"fmt"

	"../accessOp"
	"../class"
	"../classAnaly"
	"../comFunc"
	"../memoryControl"
	"../methodStack"
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
	//加载主类
	mainClass, err := classAnaly.LoadClass(startup.CmdPara.MainClass)
	if err != nil {
		fmt.Println(mainClass, err)
		panic("main()3")
	}
	//查找main方法,to do 后续补一下可访问性的判断
	methodName, err := memCtrl.PutSymbol([]byte("main"))
	if err != nil {
		fmt.Println(err)
		panic("main()4")
	}
	methodDescriptor, err := memCtrl.PutSymbol([]byte("([Ljava/lang/String;)V"))
	if err != nil {
		fmt.Println(err)
		panic("main()5")
	}
	methodInfo, codeAdr := mainClass.FindMethod(methodName, methodDescriptor)
	if methodInfo == nil || codeAdr == memCtrl.INVALID_MEM {
		fmt.Println(methodInfo, codeAdr)
		panic("main()6")
	}
	codeAttri := (*classAnaly.CODE_ATTRI)(memCtrl.GetPointer(codeAdr, classAnaly.CODE_ATTRI_SIZE))
	//创建方法栈
	methodStack := method.NewMethodStack()
	frame := methodStack.PushFrame(codeAttri.MaxLocal, codeAttri.MaxStack, mainClass.LocalAdr, 0)
	methodStack.PC = codeAdr + classAnaly.CODE_ATTRI_SIZE

	//创建main函数参数String,加入到变量区
	stringAdrs := make([]uint32, len(startup.CmdPara.Args))
	for i, arg := range startup.CmdPara.Args {
		stringAdrs[i], err = access.PutString(access.BytesToUint16([]byte(arg)))
		if err != nil {
			panic("main()7")
		}
	}
	_, arrAdr, err := access.NewArray([]byte("java/lang/String"), 4, uint32(len(stringAdrs)))
	if err != nil {
		panic("main()8")
	}
	_, arrData := access.GetArrayInfo(arrAdr)
	arr := *(*[]uint32)(comFunc.BytesToArray(arrData, 4))
	for i, v := range stringAdrs {
		arr[i] = v
	}
	frame.SetVar(0, arrAdr)
	//开始执行
	methodStack.Excute()
}
