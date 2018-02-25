package classAnaly

import (
	"bytes"
	"errors"

	"../class"
)

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
	minor_version = uint16(context[0])<<8 | uint16(context[1])
	major_version = uint16(context[2])<<8 | uint16(context[3])
	return context[4:], minor_version, major_version
}

/******************************************************************
    功能:读取常量池
	入参:文件内容
    返回值:1、读取后的context
	      2、最小版本号
		  3、主版本号
******************************************************************/
/*
func readConstantPool(context []byte) {
	size := (*uint16)(comFunc.BytesToUnsafePointer(context[0:2]))
	fmt.Println(size)
}
*/
