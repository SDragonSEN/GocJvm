package access

import (
	"../comFunc"
	"../memoryControl"
)

type ARRAY_INFO struct {
	ArrayType uint32 //Symbol
	Width     uint8  //字宽,4字节还是8字节
	Length    uint32 //数组长度
}

const ARRAY_INFO_SIZE = 9

const INIT_ARRAY_CLASS_ADR = 0xFEFEFECC

var ArrayClassAdr uint32 = INIT_ARRAY_CLASS_ADR

/******************************************************************
    功能:新建引用
	入参:无
    返回值:1、引用指针
	      2、地址
		  3、error
******************************************************************/
func NewArray(symbol []byte, width uint8, length uint32) (*ACCESS_INFO, uint32, error) {
	//新建引用
	access, accAdr, err := NewAccessInfo()
	if err != nil {
		return nil, memCtrl.INVALID_MEM, err
	}
	//分配数组数据的内存
	arrAdr, err := memCtrl.Malloc(uint32(ARRAY_INFO_SIZE)+uint32(width)*length, memCtrl.ARRAY_NODE)
	if err != nil {
		return nil, memCtrl.INVALID_MEM, err
	}
	access.DataAddr = arrAdr
	access.TypeAddr = ArrayClassAdr

	//数组描述符,字宽,长度赋值
	array := (*ARRAY_INFO)(comFunc.BytesToUnsafePointer(memCtrl.Memory[arrAdr:]))
	array.ArrayType, err = memCtrl.PutSymbol(symbol)
	if err != nil {
		return nil, memCtrl.INVALID_MEM, err
	}
	array.Width = width
	array.Length = length
	//数组数据刷成0
	for i := arrAdr + ARRAY_INFO_SIZE; i < arrAdr+ARRAY_INFO_SIZE+uint32(width)*length; i++ {
		memCtrl.Memory[i] = 0
	}
	return access, accAdr, nil
}
