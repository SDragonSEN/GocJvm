package access

import (
	"unicode/utf16"

	"../comFunc"
	"../memoryControl"
)

type STRING struct {
	ArrAdr    uint32
	HashValue uint32
}

const STRING_SIZE = 8

type HASH_ITEM struct {
	StrAdr   uint32
	HashCode uint32
}

const HASH_ITEM_SIZE = 8

const INIT_STRING_CLASS_ADR = 0xFEFEFEDD

var StringClassAdr uint32 = INIT_STRING_CLASS_ADR

var StringPoolAdr uint32 = memCtrl.INVALID_MEM

var StringPoolCap uint32 = 0

var StringPoolSize uint32 = 0

/******************************************************************
    功能:计算string的hash值
	入参:[]byte(即utf8的编码)
    返回值:hash值
******************************************************************/
func HashCode(s []uint16) uint32 {
	h := uint32(0)
	for _, v := range s {
		h = h*31 + uint32(v)
	}
	return h
}

/******************************************************************
    功能:判断utf16是否和string相等
	入参:[]uint16(即utf16的编码)
    返回值:bool
******************************************************************/
func IsEqualStr(strAdr uint32, s []uint16) bool {

	//获取string数据
	strData := GetData(strAdr)
	str := (*STRING)(comFunc.BytesToUnsafePointer(strData))
	//获取char数组数据
	_, arrData := GetArrayInfo(str.ArrAdr)
	chars := *(*[]uint16)(comFunc.BytesToArray(arrData, 2))
	//比较长度
	if len(chars) != len(s) {
		return false
	}
	//比较内容
	for i := 0; i < len(chars); i++ {
		if chars[i] != s[i] {
			return false
		}
	}
	return true
}

/******************************************************************
    功能:判断utf16是否和string相等
	入参:[]uint16(即utf16的编码)
    返回值:1、access地址
	      2、error
******************************************************************/
func PutString(s []uint16) (uint32, error) {
	var err error
	//初始化常量池
	if StringPoolAdr == memCtrl.INVALID_MEM {
		StringPoolCap = 10
		StringPoolAdr, err = memCtrl.Malloc(StringPoolCap*HASH_ITEM_SIZE, memCtrl.CONSTANT_POOL_NODE)
		if err != nil {
			return memCtrl.INVALID_MEM, err
		}
	}

	constantPool := *(*[]HASH_ITEM)(comFunc.BytesToArray(memCtrl.Memory[StringPoolAdr:StringPoolAdr+StringPoolCap*HASH_ITEM_SIZE], HASH_ITEM_SIZE))

	//计算字符串Hash值
	hash := HashCode(s)

	index := hash % StringPoolCap
	//找到对应的hash位置
	for {
		if constantPool[index].StrAdr == 0 {
			break
		}
		if IsEqualStr(constantPool[index].StrAdr, s) {
			return constantPool[index].StrAdr, nil

		}
		index++
		index %= StringPoolCap
	}
	//新建类引用
	acc, adr, err := NewAccessInfo()
	if err != nil {
		return memCtrl.INVALID_MEM, err
	}
	acc.TypeAddr = StringClassAdr
	//分配实例数据
	acc.DataAddr, err = memCtrl.Malloc(STRING_SIZE, memCtrl.INSTANCE_NODE)
	if err != nil {
		return memCtrl.INVALID_MEM, err
	}
	str := (*STRING)(comFunc.BytesToUnsafePointer(GetData(adr)))
	//新建数组实例
	_, str.ArrAdr, err = NewArray([]byte("[C"), 2, uint32(len(s)))
	if err != nil {
		return memCtrl.INVALID_MEM, err
	}
	//拷贝字符数组
	_, data := GetArrayInfo(str.ArrAdr)
	chars := *(*[]uint16)(comFunc.BytesToArray(data, 2))
	copy(chars, s)
	//常量池修改
	constantPool[index].StrAdr = adr
	constantPool[index].HashCode = hash
	StringPoolSize++
	//拓展Hash空间
	err = ExtendConstantPool()
	if err != nil {
		return memCtrl.INVALID_MEM, err
	}
	return adr, nil
}

/******************************************************************
    功能:拓展常量池的Hash空间
	入参:无
    返回值:error
******************************************************************/
func ExtendConstantPool() error {
	var err error
	//如果size/cap>4/5,将HashMap空间拓展一倍
	if StringPoolSize*5 <= 4*StringPoolCap {
		return nil
	}
	oldStringPoolAdr := StringPoolAdr
	oldConstantPool := *(*[]HASH_ITEM)(comFunc.BytesToArray(memCtrl.Memory[StringPoolAdr:StringPoolAdr+StringPoolCap*HASH_ITEM_SIZE], HASH_ITEM_SIZE))

	StringPoolCap *= 2
	StringPoolAdr, err = memCtrl.Malloc(StringPoolCap*HASH_ITEM_SIZE, memCtrl.CONSTANT_POOL_NODE)
	if err != nil {
		return err
	}
	ConstantPool := *(*[]HASH_ITEM)(comFunc.BytesToArray(memCtrl.Memory[StringPoolAdr:StringPoolAdr+StringPoolCap*HASH_ITEM_SIZE], HASH_ITEM_SIZE))
	for _, v := range oldConstantPool {
		index := v.HashCode % StringPoolCap
		for {
			if ConstantPool[index].StrAdr == 0 {
				ConstantPool[index].HashCode = v.HashCode
				ConstantPool[index].StrAdr = v.StrAdr
				break
			}
			index = (index + 1) % StringPoolCap
		}
	}
	memCtrl.MemFree(oldStringPoolAdr)
	return nil
}

/******************************************************************
    功能:[]byte转成[]uint16
	入参:[]byte(即utf8的编码)
    返回值:1、uint16（即utf16的编码）
******************************************************************/
func BytesToUint16(s []byte) []uint16 {
	//[]byte转成string,string转成[]rune,[]rune转成[]uint16
	return utf16.Encode([]rune(string(s)))
}
