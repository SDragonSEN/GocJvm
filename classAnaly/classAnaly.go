package classAnaly

import (
	"bytes"
	"errors"

	"../class"
	"../comFunc"
	"../memoryControl"
)

type CONSTANT_TYPE_16 struct {
	param1 uint16
	param2 uint16
}
type CONSTANT_TYPE_32 struct {
	param uint32
}

/*
type CONSTANT_TYPE_INT struct {
	param int32
}*/

var magicNum = []byte{0xCA, 0xFE, 0xBA, 0xBE}

func LoadClass(className string) {
	context, err := class.ReadClass(className)
	if err != nil {
		panic("classAnaly.LoadClass():class not found！")
	}
	//读取魔数
	context, err = readMagicNum(context)
	if err != nil {
		panic("classAnaly.LoadClass():魔数不正确！")
	}
	//读取版本号
	context, _, _ = readVersion(context)

	//读取常量池
}

/******************************************************************
    功能:读取魔数
	入参:文件内容
    返回值:1、读取后的context
	      2、error
******************************************************************/
func readMagicNum(context []byte) ([]byte, error) {
	if bytes.Compare(context[0:4], magicNum) != 0 {
		return nil, errors.New("classAnaly.readMagicNum():魔数不正确")
	}
	return context[4:], nil
}

/******************************************************************
    功能:读取版本号
	入参:文件内容
    返回值:1、读取后的context
	      2、最小版本号
		  3、主版本号
******************************************************************/
func readVersion(context []byte) ([]byte, uint16, uint16) {
	var minor_version, major_version uint16
	minor_version = comFunc.BytesToUint16(context[0:2])
	major_version = comFunc.BytesToUint16(context[2:4])
	return context[4:], minor_version, major_version
}

/******************************************************************
    功能:读取常量池
	入参:文件内容
    返回值:1、
******************************************************************/
func readConstantPool(context []byte) {
	//符号数量从1到size-1
	size := comFunc.BytesToUint16(context[0:2])
	//消耗的码流数量
	var count uint32
	count = 2
	//结果
	result := []byte{}
	var i uint16
	var constantBytes []byte
	var consume uint32
	for i = 1; i < size; i++ {
		tag := context[count]
		count++
		switch tag {
		//Utf8_info
		case 0x01:
			constantBytes, consume = readConstantUtf8Info(context[count:])
		//Integer_info
		case 0x03:
			constantBytes, consume = readConstantIntegerInfo(context[count:])
		//Float_info
		case 0x04:
			constantBytes, consume = readConstantFloatInfo(context[count:])
		//LongInfo
		case 0x05:
			//一个long型要占两个4字节和两个slot位(slot字宽固定32位，即使在64位的机子上也一样)
			i++
			constantBytes, consume = readConstantLongInfo(context[count:])
		//LongInfo
		case 0x06:
			//一个Double型要占两个4字节和两个slot位(slot字宽固定32位，即使在64位的机子上也一样)
			i++
			constantBytes, consume = readConstantDoubleInfo(context[count:])

		}
		count += consume
		result = append(result, constantBytes...)
	}
}

/******************************************************************
    功能:读取UTF8_INFO
	入参:文件内容
    返回值:1、转化后的码流，即符号表中的地址
	      2、消耗的码流数量
******************************************************************/
func readConstantUtf8Info(context []byte) ([]byte, uint32) {
	//获取utf8长度
	length := comFunc.BytesToUint16(context[0:2])
	var count uint32
	count = 2
	result := [4]byte{}
	constant32 := (*CONSTANT_TYPE_32)(comFunc.BytesToUnsafePointer(result[:]))
	//将utf8码流加到符号表中
	constant32.param = memCtrl.PutSymbol(context[count : count+uint32(length)])
	count += uint32(length)
	return result[:], count
}

/******************************************************************
    功能:读取INTEGER_INFO
	入参:文件内容
    返回值:1、转化后的码流，即Integer值
	      2、消耗的码流数量
******************************************************************/
func readConstantIntegerInfo(context []byte) ([]byte, uint32) {
	//获取integer值
	integer := comFunc.BytesToUint32(context[0:4])
	result := [4]byte{}
	constantInt := (*CONSTANT_TYPE_32)(comFunc.BytesToUnsafePointer(result[:]))
	constantInt.param = integer
	return result[:], 4
}

/******************************************************************
    功能:读取FLOAT_INFO
	入参:文件内容
    返回值:1、转化后的码流，即Float值
	      2、消耗的码流数量
******************************************************************/
func readConstantFloatInfo(context []byte) ([]byte, uint32) {
	//实现同Integer
	return readConstantIntegerInfo(context)
}

/******************************************************************
    功能:读取LONG_INFO
	入参:文件内容
    返回值:1、转化后的码流，高位在低地址，低位在高地址
	      2、消耗的码流数量
******************************************************************/
func readConstantLongInfo(context []byte) ([]byte, uint32) {
	result := [8]byte{}
	//获取
	long := comFunc.BytesToUint32(context[0:4])

	constant64 := (*CONSTANT_TYPE_32)(comFunc.BytesToUnsafePointer(result[0:4]))
	constant64.param = long

	long = comFunc.BytesToUint32(context[4:8])
	constant64 = (*CONSTANT_TYPE_32)(comFunc.BytesToUnsafePointer(result[4:8]))
	constant64.param = long
	return result[:], 8
}

/******************************************************************
    功能:读取DOUBLE_INFO
	入参:文件内容
    返回值:1、转化后的码流，高位在低地址，低位在高地址
	      2、消耗的码流数量
******************************************************************/
func readConstantDoubleInfo(context []byte) ([]byte, uint32) {
	//实现同Long
	return readConstantLongInfo(context)
}
