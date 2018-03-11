package classAnaly

import (
	"errors"

	"../comFunc"
	"../memoryControl"
)

/******************************************************************
    功能:从常量池中读取class name
	入参:1、常量池
	    2、Class索引
    返回值:1、符号表地址
	      2、error
	注:常量池是从1开始计算的，不是从0
******************************************************************/
func GetClassFromConstPool(constPool []byte, index uint32) (uint32, error) {
	if index == 0 {
		return memCtrl.INVALID_MEM, errors.New("GetClassFromConstPool():索引为0!")
	}
	if index > uint32(len(constPool)) {
		return memCtrl.INVALID_MEM, errors.New("GetClassFromConstPool():索引越界!")
	}
	constItem := (*CONSTANT_TYPE_32)(comFunc.BytesToUnsafePointer(constPool[(index-1)*4:]))
	return GetUtf8FromConstPool(constPool, constItem.param)
}

/******************************************************************
    功能:从常量池中读取UTF8
	入参:1、常量池
	    2、Utf8索引
    返回值:1、符号表地址
	      2、error
	注:常量池是从1开始计算的，不是从0
******************************************************************/
func GetUtf8FromConstPool(constPool []byte, index uint32) (uint32, error) {
	if index == 0 {
		return memCtrl.INVALID_MEM, errors.New("GetUtf8FromConstPool():索引为0!")
	}
	if index > uint32(len(constPool)) {
		return memCtrl.INVALID_MEM, errors.New("GetUtf8FromConstPool():索引越界!")
	}
	constItem := (*CONSTANT_TYPE_32)(comFunc.BytesToUnsafePointer(constPool[(index-1)*4:]))
	return constItem.param, nil
}
