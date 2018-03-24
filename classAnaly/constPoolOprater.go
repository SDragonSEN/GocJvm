package classAnaly

import (
	"errors"

	"../comFunc"
	"../memoryControl"
)

const ACONST_NULL = 0x01 //null值入栈
const ICONST_M1 = 0x02   //-1(int)入栈
const ICONST_0 = 0x03    //0(int)入栈，下同
const ICONST_1 = 0x04
const ICONST_2 = 0x05
const ICONST_3 = 0x06
const ICONST_4 = 0x07
const ICONST_5 = 0x08
const LCONST_0 = 0x09 //0(long)入栈
const LCONST_1 = 0x0a
const FCONST_0 = 0x0b //0(float)入栈
const FCONST_1 = 0x0c
const FCONST_2 = 0x0d
const DCONST_0 = 0x0e //0(double)入栈
const DCONST_1 = 0x0f
const BIPUSH = 0x10 //操作数byte,拓展成int型入栈
const SIPUSH = 0x11 //操作数int16,拓展成int型入栈
const LDC = 0x12    //操作数byte,将常量(int,float,string)入栈
const LDC_W = 0x13  //操作数int16,将常量入栈
const LDC2_W = 0x14 //操作数int16,将常量(long,double)入栈

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
