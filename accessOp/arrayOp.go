package access

import (
	"../comFunc"
	"../memoryControl"
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
func NewArray(symbol []byte, width uint32, length uint32) (*ACCESS_INFO, uint32, error) {
	//新建引用
	access, accAdr, err := NewAccessInfo()
	if err != nil {
		return nil, memCtrl.INVALID_MEM, err
	}
	//分配数组数据的内存
	leng := uint32(ARRAY_INFO_SIZE) + uint32(width)*length
	arrAdr, err := memCtrl.Malloc(leng, memCtrl.ARRAY_NODE)
	if err != nil {
		return nil, memCtrl.INVALID_MEM, err
	}
	access.DataAddr = arrAdr
	access.TypeAddr = ArrayClassAdr
	//数组描述符,字宽,长度赋值
	array := (*ARRAY_INFO)(comFunc.BytesToUnsafePointer(memCtrl.Memory[arrAdr : arrAdr+leng]))
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

/******************************************************************
    功能:获取数组信息
	入参:access地址
    返回值:1、ARRAY_INFO
	      2、数组数据切片
******************************************************************/
func GetArrayInfo(accAdr uint32) (ARRAY_INFO, []byte) {
	data := GetData(accAdr)
	arr := (*ARRAY_INFO)(comFunc.BytesToUnsafePointer(data[0:ARRAY_INFO_SIZE]))
	return *arr, data[ARRAY_INFO_SIZE : ARRAY_INFO_SIZE+arr.Length*uint32(arr.Width)]
}
