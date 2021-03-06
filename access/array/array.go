package array

import (
	. "access/access"
	. "basic/com"
	. "basic/memCtrl"
)

type ARRAY_INFO struct {
	ArrayType uint32 //Symbol地址
	Length    uint32 //数组长度
	Width     uint32 //字宽,4字节还是8字节
}

const ARRAY_INFO_SIZE = 3 * 4

const INIT_ARRAY_CLASS_ADR = 0xFEFEFECC

var ArrayClassAdr uint32 = INIT_ARRAY_CLASS_ADR

/******************************************************************
    功能:新建引用
	入参:无
    返回值:1、引用指针
	      2、地址
		  3、error
******************************************************************/
func NewArray(symbol, width, length uint32) (*ACCESS_INFO, uint32, error) {
	//新建引用
	access, accAdr, err := NewAccessInfo()
	if err != nil {
		return nil, INVALID_MEM, err
	}
	//分配数组数据的内存
	leng := ARRAY_INFO_SIZE + width*length
	arrAdr, err := Malloc(leng, ARRAY_NODE)
	if err != nil {
		return nil, INVALID_MEM, err
	}
	access.DataAddr = arrAdr
	access.TypeAddr = ArrayClassAdr
	//数组描述符,字宽,长度赋值
	array := (*ARRAY_INFO)(GetPointer(arrAdr, ARRAY_INFO_SIZE))
	array.ArrayType = symbol
	array.Width = width
	array.Length = length

	//数组数据刷成0
	for i := arrAdr + ARRAY_INFO_SIZE; i < arrAdr+ARRAY_INFO_SIZE+width*length; i++ {
		Memory[i] = 0
	}
	return access, accAdr, nil
}

/******************************************************************
    功能:获取数组信息
	入参:access地址
    返回值:1、ARRAY_INFO
	      2、数组数据切片
******************************************************************/
func GetArrayInfo(accAdr uint32) (ARRAY_INFO, []byte) {
	data := GetData(accAdr)
	arr := (*ARRAY_INFO)(BytesToUnsafePointer(data[0:ARRAY_INFO_SIZE]))
	return *arr, data[ARRAY_INFO_SIZE : ARRAY_INFO_SIZE+arr.Length*uint32(arr.Width)]
}
