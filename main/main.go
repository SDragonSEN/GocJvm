// main project main.go
package main

import (
	"fmt"

	. "access/access"
	"methodStack"
	. "startup"

	. "access/array"
	. "access/string"
	. "basic/com"
	. "basic/memCtrl"
	. "basic/symbol"
	. "class/classFind"
	. "class/classInterface"
	. "class/classParse"
)

func main() {

	ParseCmd()
	InitClassPath(CmdPara.Jar)
	object, err := LoadClass("java.lang.Object")
	if err != nil {
		panic("main()1")
	}
	method.CInit(object.LocalAdr)

	//加载数组
	arrayTypeAdr := LoadArrayClass()
	ModifyTypeAddr(INIT_ARRAY_CLASS_ADR, arrayTypeAdr)
	ArrayClassAdr = arrayTypeAdr
	//加载String
	stringClass, err := LoadClass("java.lang.String")
	if err != nil {
		fmt.Println("error:main():", err)
		panic("main()2")
	}
	ModifyTypeAddr(INIT_STRING_CLASS_ADR, stringClass.LocalAdr)
	StringClassAdr = stringClass.LocalAdr
	//加载主类
	mainClass, err := LoadClass(CmdPara.MainClass)
	if err != nil {
		fmt.Println(mainClass, err)
		panic("main()3")
	}
	method.CInit(mainClass.LocalAdr)
	//查找main方法,to do 后续补一下可访问性的判断
	methodName, err := PutSymbol([]byte("main"))
	if err != nil {
		fmt.Println(err)
		panic("main()4")
	}
	methodDescriptor, err := PutSymbol([]byte("([Ljava/lang/String;)V"))
	if err != nil {
		fmt.Println(err)
		panic("main()5")
	}
	methodInfo, codeAdr := mainClass.FindMethod(methodName, methodDescriptor)
	if methodInfo == nil || codeAdr == INVALID_MEM {
		fmt.Println(methodInfo, codeAdr)
		panic("main()6")
	}
	codeAttri := (*CODE_ATTRI)(GetPointer(codeAdr, CODE_ATTRI_SIZE))
	//创建方法栈
	methodStack := method.NewMethodStack()
	frame := methodStack.PushFrame(codeAttri.MaxLocal, codeAttri.MaxStack, mainClass.LocalAdr, 0)
	methodStack.PC = codeAdr + CODE_ATTRI_SIZE

	//创建main函数参数String,加入到变量区
	stringAdrs := make([]uint32, len(CmdPara.Args))
	for i, arg := range CmdPara.Args {
		stringAdrs[i], err = PutString(BytesToUtf16([]byte(arg)))
		if err != nil {
			panic("main()7")
		}
	}
	_, arrAdr, err := NewArray(SYM_Kjava_lang_String, 4, uint32(len(stringAdrs)))
	if err != nil {
		panic("main()8")
	}
	_, arrData := GetArrayInfo(arrAdr)
	arr := *(*[]uint32)(BytesToArray(arrData, 4))
	for i, v := range stringAdrs {
		arr[i] = v
	}
	frame.SetVar(0, arrAdr)
	//开始执行
	methodStack.Excute()
}
