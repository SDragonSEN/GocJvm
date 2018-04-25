package classInterface

import (
	. "basic/com"
	. "basic/memCtrl"
	. "basic/symbol"
	. "class/classTable"
)

/******************************************************************
    功能:加载数组的定义,自定义
	入参:无
    返回值:1、类型地址
******************************************************************/
func LoadArrayClass() uint32 {
	var err error
	arrayClassMem := make([]byte, CLASS_INFO_SIZE)
	classInfo := (*CLASS_INFO)(BytesToUnsafePointer(arrayClassMem))
	classInfo.ClassName, err = PutSymbol([]byte("java.lang.array"))
	if err != nil {
		panic("LoadArrayClass()1")
	}
	objectClass, err := PutSymbol([]byte("java.lang.Object"))
	if err != nil {
		panic("LoadArrayClass()2")
	}
	classInfo.SuperClassAddr = GetClassMemAddr(objectClass)
	classInfo.AccessFlag = CLASS_ACC_PUBLIC | CLASS_ACC_FINAL | CLASS_ACC_SUPER
	classInfo.ConstNum = 0
	classInfo.FiledInfoDev = CLASS_INFO_SIZE
	classInfo.UnstaticParaDev = CLASS_INFO_SIZE
	classInfo.UnstaticParaSize = 0
	classInfo.UnstaticParaTotalSize = 0
	classInfo.StaticParaDev = CLASS_INFO_SIZE
	classInfo.StaticParaSize = 0
	classInfo.InterfaceDev = CLASS_INFO_SIZE
	classInfo.InterfaceNum = 0
	classInfo.MethodDev = CLASS_INFO_SIZE
	classInfo.MethodNum = 0
	memAdr, err := PutClass(classInfo.ClassName, arrayClassMem)
	(*CLASS_INFO)(GetPointer(memAdr, CLASS_INFO_SIZE)).LocalAdr = memAdr

	if err != nil {
		panic("LoadArrayClass()3")
	}
	return memAdr
}
