package access

import (
	"unicode/utf16"

	_ "../comFunc"
	"../memoryControl"
)

var StringHeaderAdr uint32 = memCtrl.INVALID_MEM

const INIT_STRING_CLASS_ADR = 0xFEFEFEDD

var StringClassAdr uint32 = INIT_STRING_CLASS_ADR

/*
func PutString(s []uint16) (*ACCESS_INFO, uint32, error) {
	if StringHeaderAdr == memCtrl.INVALID_MEM {
		//var accArray *ACCESS_INFO
		acc, adr, err := NewAccessInfo()
		acc.TypeAddr = StringClassAdr
		if err != nil {

		}
		length := uint32(len(s))
		_, acc.DataAddr, err = NewArray([]byte("[C"), 2, length)
		if err != nil {

		}
		//accArray := (*ACCESS_INFO)(comFunc.BytesToUnsafePointer(memCtrl.Memroy[acc.DataAddr:]))
		data := memCtrl.Memory[acc.DataAddr+ARRAY_INFO_SIZE : acc.DataAddr+ARRAY_INFO_SIZE+2*length]
		u16 := *(*[]uint16)(comFunc.BytesToArray(data, 2))
		copy(u16, s)
		return acc, adr, nil
	}
}
*/
/******************************************************************
    功能:[]byte转成[]uint16
	入参:[]byte(即utf8的编码)
    返回值:1、uint16（即utf16的编码）
******************************************************************/
func BytesToUint16(s []byte) []uint16 {
	//[]byte转成string,string转成[]rune,[]rune转成[]uint16
	return utf16.Encode([]rune(string(s)))
}
