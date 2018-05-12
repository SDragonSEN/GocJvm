package method

import (
	. "access/access"
	. "access/array"
	. "basic/com"
	. "basic/memCtrl"
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
	MethodMap["java/lang/Class~getClassLoader0~()Ljava/lang/ClassLoader;"] = java_lang_Class_getClassLoader0
	MethodMap["java/lang/Class~desiredAssertionStatus0~(Ljava/lang/Class;)Z"] = java_lang_Class_desiredAssertStatus0
	MethodMap["java/lang/System~arraycopy~(Ljava/lang/Object;ILjava/lang/Object;II)V"] = java_lang_System_arraycopy
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
		//fmt.Println(className + "~" + methodName + "~" + methodDesp + "被执行")
		return
	}
	fmt.Println("未知的native方法:" + className + "~" + methodName + "~" + methodDesp)
}

/******************************************************************
    功能:java/lang/Object的registerNatives方法
	入参:frame
    返回值:无
******************************************************************/
func java_lang_Object_registerNatives(frame *METHOD_FRAME) {
}

/******************************************************************
    功能:java/lang/Class的registerNatives方法
	入参:frame
    返回值:无
******************************************************************/
func java_lang_Class_registerNatives(frame *METHOD_FRAME) {
}

/******************************************************************
    功能:java/lang/Class的registerNatives方法
	入参:frame
    返回值:无
******************************************************************/
func java_lang_Class_getPrimitiveClass(frame *METHOD_FRAME) {
	frame.Pop()
	frame.Push(0)
}

/******************************************************************
    功能:java/lang/Float的floatToRawIntBits方法
	入参:frame
    返回值:无
******************************************************************/
func java_lang_Float_floatToRawIntBits(frame *METHOD_FRAME) {
	//因为java和golang都是用的IEE 754标准的浮点数，所以这里可以不做处理
}

/******************************************************************
    功能:java/lang/Class的getClassLoader0方法
	入参:frame
    返回值:无
******************************************************************/
func java_lang_Class_getClassLoader0(frame *METHOD_FRAME) {
	frame.Pop()
	frame.Push(0)
}

/******************************************************************
    功能:java/lang/Class的desiredAssertStatus0方法
	入参:frame
    返回值:无
******************************************************************/
func java_lang_Class_desiredAssertStatus0(frame *METHOD_FRAME) {
	frame.Pop()
	frame.Push(GOJVM_FALSE)
}

/******************************************************************
    功能:java/lang/Class的arraycopy方法
	入参:frame
    返回值:无
******************************************************************/
func java_lang_System_arraycopy(frame *METHOD_FRAME) {
	length := int32(frame.Pop())
	destPos := int32(frame.Pop())
	dest := frame.Pop()
	srcPos := int32(frame.Pop())
	src := frame.Pop()
	if !IsArrayAccess(src) || !IsArrayAccess(dest) {
		panic("java_lang_System_arraycopy():src或dest不是Array类型！")
	}
	if length == 0 {
		return
	}
	if length < 0 || srcPos < 0 || destPos < 0 {
		fmt.Println(length, srcPos, destPos)
		panic("java_lang_System_arraycopy():长度有负数！")
	}
	srcArray, srcData := GetArrayInfo(src)
	destArray, destData := GetArrayInfo(dest)
	if uint32(srcPos+length) > srcArray.Length ||
		uint32(destPos+length) > destArray.Length {
		panic("java_lang_System_arraycopy():长度越界！")
	}
	r1 := []rune(string(GetSymbol(srcArray.ArrayType)))
	r2 := []rune(string(GetSymbol(destArray.ArrayType)))
	if r1[1] != r2[1] {
		panic("java_lang_System_arraycopy():类型不一致！")
	}
	if r1[1] == 'L' {
		//待完善
	} else if r1[1] == '[' {
		//待完善
	} else {
		copyBasicArray(srcData, srcPos, destData, destPos, length, int32(srcArray.Width))
	}
}

/******************************************************************
    功能:拷贝基本类型数组
	入参:
    返回值:无
******************************************************************/
func copyBasicArray(srcData []byte, srcPos int32, destData []byte, destPos, length, width int32) {
	srcPos *= width
	destPos *= width
	for i := int32(0); i < length*width; i++ {
		destData[i+destPos] = srcData[i+srcPos]
	}
}

/******************************************************************
    功能:判断是否为数组类型
	入参:引用地址
    返回值:无
******************************************************************/
func IsArrayAccess(adr uint32) bool {
	access := (*ACCESS_INFO)(GetPointer(adr, ACCESS_INFO_SIZE))
	if access.TypeAddr == ArrayClassAdr {
		return true
	}
	return false
}
