package classInterface

import (
	. "basic/com"
	. "basic/memCtrl"
)

type FiledInfo struct {
	ClassName uint32 //都是符号表中的地址
	FiledName uint32
	FiledType uint32
}
type MethodInfo struct {
	ClassName  uint32 //都是符号表中的地址
	MethodName uint32
	MethodDesp uint32
}

/******************************************************************
    功能:获取常量池中MethodRef类型
	入参:1、常量池切片
	    2、常量池中的索引
    返回值:MethodInfo
******************************************************************/
func GetMethodInfo(constPool []byte, index uint32) MethodInfo {
	var methodInfo MethodInfo
	if index == 0 {
		panic("GetStaticFiledInfo()error0")
	}
	if index > uint32(len(constPool)/4) {
		panic("GetStaticFiledInfo()error1")
	}
	constItem := (*CONSTANT_TYPE_16)(BytesToUnsafePointer(constPool[(index-1)*4:]))
	//class类型在常量池中的位置
	classIndex := GetUint32FromConstPool(constPool, uint32(constItem.Param1))
	//class的名字在符号表里的位置
	methodInfo.ClassName = GetUint32FromConstPool(constPool, classIndex)
	//Name和Type在符号表里的位置
	methodInfo.MethodName, methodInfo.MethodDesp = GetNameAndType(constPool, uint32(constItem.Param2))
	return methodInfo
}

/******************************************************************
    功能:获取常量池中FiledRef类型
	入参:1、常量池切片
	    2、常量池中的索引
    返回值:FiledInfo
******************************************************************/
func GetFiledInfo(constPool []byte, index uint32) FiledInfo {
	var filedInfo FiledInfo
	if index == 0 {
		panic("GetStaticFiledInfo()error0")
	}
	if index > uint32(len(constPool)/4) {
		panic("GetStaticFiledInfo()error1")
	}
	constItem := (*CONSTANT_TYPE_16)(BytesToUnsafePointer(constPool[(index-1)*4:]))
	//class类型在常量池中的位置
	classIndex := GetUint32FromConstPool(constPool, uint32(constItem.Param1))
	//class的名字在符号表里的位置
	filedInfo.ClassName = GetUint32FromConstPool(constPool, classIndex)
	//Name和Type在符号表里的位置
	filedInfo.FiledName, filedInfo.FiledType = GetNameAndType(constPool, uint32(constItem.Param2))
	return filedInfo
}

/******************************************************************
    功能:获取常量池中NameAndType类型
	入参:1、常量池切片
	    2、常量池中的索引
    返回值:1、Name在符号表中的位置
	      2、Type在符号表中的位置
******************************************************************/
func GetNameAndType(constPool []byte, index uint32) (uint32, uint32) {
	if index == 0 {
		panic("GetNameAndType()error0")
	}
	if index > uint32(len(constPool)/4) {
		panic("GetNameAndType()error1")
	}
	constItem := (*CONSTANT_TYPE_16)(BytesToUnsafePointer(constPool[(index-1)*4:]))

	return GetUint32FromConstPool(constPool, uint32(constItem.Param1)), GetUint32FromConstPool(constPool, uint32(constItem.Param2))
}

/******************************************************************
    功能:获取常量池切片
	入参:1、Class索引
    返回值:1、常量池切片
******************************************************************/
func GetConstantPoolSlice(classAdr uint32) []byte {
	classInfo := (*CLASS_INFO)(GetPointer(classAdr, CLASS_INFO_SIZE))
	constPoolAdr := classAdr + CLASS_INFO_SIZE
	return Memory[constPoolAdr : constPoolAdr+classInfo.ConstNum*4]
}

/******************************************************************
    功能:是否是Class常量
	入参:1、Class索引
    返回值:1、常量池切片
******************************************************************/
func IsClassConstant(classAdr, index uint32) bool {
	classInfo := (*CLASS_INFO)(GetPointer(classAdr, CLASS_INFO_SIZE))

	m := Memory[classInfo.LocalAdr+classInfo.ClassConstDev : classInfo.LocalAdr+classInfo.InterfaceDev]
	for _, v := range *(*[]uint32)(BytesToArray(m, 4)) {
		if v == index {
			return true
		}
	}
	return false
}

/******************************************************************
    功能:从常量池中读取class name
	入参:1、常量池
	    2、Class索引
    返回值:1、符号表地址
	注:常量池是从1开始计算的，不是从0
******************************************************************/
func GetClassNameFromConstPool(constPool []byte, index uint32) uint32 {
	if index == 0 {
		panic("GetUtf8FromConstPool() 1")
	}
	if index > uint32(len(constPool)/4) {
		panic("GetUtf8FromConstPool() 2")
	}
	constItem := (*CONSTANT_TYPE_32)(BytesToUnsafePointer(constPool[(index-1)*4:]))
	return GetUtf8FromConstPool(constPool, constItem.Param)
}

/******************************************************************
    功能:从常量池中读取UTF8
	入参:1、常量池
	    2、Utf8索引
    返回值:1、符号表地址
	      2、error
	注:常量池是从1开始计算的，不是从0
******************************************************************/
func GetUtf8FromConstPool(constPool []byte, index uint32) uint32 {
	if index == 0 {
		panic("GetUtf8FromConstPool() 1")
	}
	if index > uint32(len(constPool)/4) {
		panic("GetUtf8FromConstPool() 2")
	}
	constItem := (*CONSTANT_TYPE_32)(BytesToUnsafePointer(constPool[(index-1)*4:]))
	return constItem.Param
}

/******************************************************************
    功能:从常量池中读取Uint32
	入参:1、常量池
	    2、常量池索引
    返回值:1、Uint32值
	      2、error
	注:常量池是从1开始计算的，不是从0
******************************************************************/
func GetUint32FromConstPool(constPool []byte, index uint32) uint32 {
	if index == 0 {
		panic("GetUint32FromConstPool() 1")
	}
	if index > uint32(len(constPool)/4) {
		panic("GetUint32FromConstPool() 2")
	}
	constItem := (*CONSTANT_TYPE_32)(BytesToUnsafePointer(constPool[(index-1)*4:]))
	return constItem.Param
}

/******************************************************************
    功能:从常量池中读取Uint64
	入参:1、常量池
	    2、常量池索引
    返回值:1、Uint32值
	      2、error
	注:常量池是从1开始计算的，不是从0
******************************************************************/
func GetUint64FromConstPool(constPool []byte, index uint32) (uint32, uint32) {
	if index == 0 {
		panic("GetUint32FromConstPool() 1")
	}
	if index > uint32(len(constPool)/4) {
		panic("GetUint32FromConstPool() 2")
	}
	constItem := (*CONSTANT_TYPE_32)(BytesToUnsafePointer(constPool[(index-1)*4:]))
	constItem1 := (*CONSTANT_TYPE_32)(BytesToUnsafePointer(constPool[index*4:]))
	return constItem.Param, constItem1.Param
}

/******************************************************************
    功能:从常量池中读取String
	入参:1、常量池
	    2、String索引
    返回值:1、实例地址
	      2、error
	注:常量池是从1开始计算的，不是从0
******************************************************************/
func GetStringFromConstPool(constPool []byte, index uint32) uint32 {
	//实现同utf8
	return GetUtf8FromConstPool(constPool, index)
}

/******************************************************************
    功能:从常量池中读取Class
	入参:1、常量池
	    2、String索引
    返回值:1、实例地址
	注:常量池是从1开始计算的，不是从0
******************************************************************/
func GetClassFromConstPool(constPool []byte, index uint32) uint32 {
	classIndex := GetUint32FromConstPool(constPool, index)
	//class的名字在符号表里的位置
	return GetUint32FromConstPool(constPool, classIndex)
}
