package method

import (
	. "basic/symbol"
	"fmt"
)

type LocalMethod func(frame *METHOD_FRAME)

var MethodMap map[string]LocalMethod

/******************************************************************
    功能:初始化方法
	入参:无
    返回值:无
******************************************************************/
func init() {
	MethodMap = make(map[string]LocalMethod)
	MethodMap["java/lang/Object~registerNatives~()V"] = java_lang_Object_registerNatives
	MethodMap["java/lang/Class~registerNatives~()V"] = java_lang_Class_registerNatives
	MethodMap["java/lang/Class~getPrimitiveClass~(Ljava/lang/String;)Ljava/lang/Class;"] = java_lang_Class_getPrimitiveClass
	MethodMap["java/lang/Float~floatToRawIntBits~(F)I"] = java_lang_Float_floatToRawIntBits
}

/******************************************************************
    功能:执行本地方法(适配版)
	入参:无
    返回值:无
******************************************************************/
func ExcuteLocalMethodAdp(className, methodName, methodDesp uint32, frame *METHOD_FRAME) {
	ExcuteLocalMethod(string(GetSymbol(className)), string(GetSymbol(methodName)), string(GetSymbol(methodDesp)), frame)
}

/******************************************************************
    功能:执行本地方法
	入参:无
    返回值:无
******************************************************************/
func ExcuteLocalMethod(className, methodName, methodDesp string, frame *METHOD_FRAME) {
	localMethod := MethodMap[className+"~"+methodName+"~"+methodDesp]
	if localMethod != nil {
		localMethod(frame)
		return
	}
	fmt.Println("未知的native方法:" + className + "~" + methodName + "~" + methodDesp)
}

/******************************************************************
    功能:java/lang/Object的registerNatives方法
	入参:frame
    返回值:无
	注:函数栈大小默认为100
******************************************************************/
func java_lang_Object_registerNatives(frame *METHOD_FRAME) {
	fmt.Println("java/lang/Object~registerNatives~()V被执行")
}

/******************************************************************
    功能:java/lang/Class的registerNatives方法
	入参:frame
    返回值:无
	注:函数栈大小默认为100
******************************************************************/
func java_lang_Class_registerNatives(frame *METHOD_FRAME) {
	fmt.Println("java/lang/Class~registerNatives~()V被执行")
}

/******************************************************************
    功能:java/lang/Class的registerNatives方法
	入参:frame
    返回值:无
	注:函数栈大小默认为100
******************************************************************/
func java_lang_Class_getPrimitiveClass(frame *METHOD_FRAME) {
	fmt.Println("java/lang/Class~getPrimitiveClass~(Ljava/lang/String;)Ljava/lang/Class;被执行")
	frame.Pop()
	frame.Push(0)
}

/******************************************************************
    功能:java/lang/Class的registerNatives方法
	入参:frame
    返回值:无
	注:函数栈大小默认为100
******************************************************************/
func java_lang_Float_floatToRawIntBits(frame *METHOD_FRAME) {
	fmt.Println("java/lang/Float~floatToRawIntBits~(F)I被执行")
	//因为java和golang都是用的IEE 754标准的浮点数，所以这里可以不做处理
}
