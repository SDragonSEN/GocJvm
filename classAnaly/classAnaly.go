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

type CLASS_INFO struct {
	ClassName             uint32 //类名
	SuperClassAddr        uint32 //父类地址,为0代表是Object类
	AccessFlag            uint16 //可访问属性
	ConstNum              uint32 //常量数量
	ParaInfoDev           uint32 //参数信息偏移
	UnstaticParaDev       uint32 //非static参数地址
	UnstaticParaSize      uint32 //非static参数大小
	UnstaticParaTotalSize uint32 //非static参数内存总大小(即，分配实例的大小)
	StaticParaDev         uint32 //static参数地址
	StaticParaSize        uint32 //static参数大小
	StaticParaTotalSize   uint32 //static参数内存总大小(即，类变量的大小)
	InterfaceDev          uint32 //接口定义偏移
	InterfaceNum          uint32 //接口数量
}

const CLASS_INFO_SIZE = 50

var magicNum = []byte{0xCA, 0xFE, 0xBA, 0xBE}

func LoadClass(className string) (uint32, error) {

	//读取字节码文件内容
	context, err := class.ReadClass(className)
	if err != nil {
		return memCtrl.INVALID_MEM, err
	}

	//解析结果定义
	result := make([]byte, CLASS_INFO_SIZE)

	//Class Info信息定义
	classInfo := (*CLASS_INFO)(comFunc.BytesToUnsafePointer(result[0:]))

	//读取魔数
	context, err = readMagicNum(context)
	if err != nil {
		return memCtrl.INVALID_MEM, err
	}

	//读取版本号，暂时不使用
	context, _, _ = readVersion(context)

	//读取常量池
	context, constPool, num, err := readConstantPool(context)
	classInfo.ConstNum = num
	if err != nil {
		return memCtrl.INVALID_MEM, err
	}
	result = append(result, constPool...)

	//读取类信息
	context, classInfo.AccessFlag, classInfo.ClassName, classInfo.SuperClassAddr, err = readClassInfo(context, constPool)
	if err != nil {
		return memCtrl.INVALID_MEM, err
	}

	//读取接口信息
	classInfo.InterfaceDev = uint32(len(result))
	context, num, interfaceInfo, err := readInterfaces(context, constPool)
	if err != nil {
		return memCtrl.INVALID_MEM, err
	}
	classInfo.InterfaceNum = num
	result = append(result, interfaceInfo...)

	//保存到内存中
	memAdr, err := memCtrl.PutClass(classInfo.ClassName, result)
	if err != nil {
		return memCtrl.INVALID_MEM, err
	}

	return memAdr, nil
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
    返回值:1、读取后的context
	      2、解析后的常量池码流
		  3、常量池数量
		  3、error
******************************************************************/
func readConstantPool(context []byte) ([]byte, []byte, uint32, error) {
	//符号数量从1到size-1
	size := comFunc.BytesToUint16(context[0:2])
	//消耗的码流数量
	var count uint32
	count = 2
	//结果
	result := make([]byte, 0)

	var i uint16
	var constantBytes []byte
	var consume uint32
	var err error
	for i = 1; i < size; i++ {
		tag := context[count]
		count++
		switch tag {
		//Utf8_info
		case 0x01:
			constantBytes, consume, err = readConstantUtf8Info(context[count:])
			if err != nil {
				return nil, nil, 0, err
			}
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
		//Long_info
		case 0x06:
			//一个Double型要占两个4字节和两个slot位(slot字宽固定32位，即使在64位的机子上也一样)
			i++
			constantBytes, consume = readConstantDoubleInfo(context[count:])
		//Class_info
		case 0x07:
			constantBytes, consume = readConstantClassInfo(context[count:])
		//String_info
		case 0x08:
			constantBytes, consume = readConstantStringInfo(context[count:])
		//Fieldref_info
		case 0x09:
			constantBytes, consume = readConstantFieldrefInfo(context[count:])
		//Methodref_info
		case 0x0A:
			constantBytes, consume = readConstantMethodrefInfo(context[count:])
		//InterfaceMethodref_info
		case 0x0B:
			constantBytes, consume = readConstantInterfaceMethodrefInfo(context[count:])
		//NameAndType_info
		case 0x0C:
			constantBytes, consume = readConstantNameAndTypeInfo(context[count:])
		default:
			return nil, nil, 0, errors.New("常量池解析错误")
		}
		count += consume
		result = append(result, constantBytes...)
	}
	return context[count:], result, uint32(size), nil
}

/******************************************************************
    功能:读取UTF8_INFO
	入参:文件内容
    返回值:1、转化后的码流，即符号表中的地址
	      2、消耗的码流数量
******************************************************************/
func readConstantUtf8Info(context []byte) ([]byte, uint32, error) {
	//获取utf8长度
	length := comFunc.BytesToUint16(context[0:2])
	var count uint32
	var err error
	count = 2
	result := [4]byte{}
	constant32 := (*CONSTANT_TYPE_32)(comFunc.BytesToUnsafePointer(result[:]))
	//将utf8码流加到符号表中
	constant32.param, err = memCtrl.PutSymbol(context[count : count+uint32(length)])
	if err != nil {
		return nil, 0, err
	}
	count += uint32(length)
	return result[:], count, nil
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
	//获取高位
	long := comFunc.BytesToUint32(context[0:4])

	constant64 := (*CONSTANT_TYPE_32)(comFunc.BytesToUnsafePointer(result[0:4]))
	constant64.param = long
	//获取低位
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

/******************************************************************
    功能:读取CLASS_INFO
	入参:文件内容
    返回值:1、转化后的码流，即常量Index值
	      2、消耗的码流数量
******************************************************************/
func readConstantClassInfo(context []byte) ([]byte, uint32) {
	//获取index值
	index := comFunc.BytesToUint16(context[0:2])
	result := [4]byte{}
	constantInt := (*CONSTANT_TYPE_32)(comFunc.BytesToUnsafePointer(result[:]))
	constantInt.param = uint32(index)
	return result[:], 2
}

/******************************************************************
    功能:读取STRING_INFO
	入参:文件内容
    返回值:1、转化后的码流，即常量Index值
	      2、消耗的码流数量
******************************************************************/
func readConstantStringInfo(context []byte) ([]byte, uint32) {
	//实现同Class
	return readConstantClassInfo(context)
}

/******************************************************************
    功能:读取Fieldref_INFO
	入参:文件内容
    返回值:1、转化后的码流，即常量Index值
	      2、消耗的码流数量
******************************************************************/
func readConstantFieldrefInfo(context []byte) ([]byte, uint32) {

	result := [4]byte{}
	constantInt := (*CONSTANT_TYPE_16)(comFunc.BytesToUnsafePointer(result[:]))
	//获取class index值
	index := comFunc.BytesToUint16(context[0:2])
	constantInt.param1 = index
	//获取name and type index值
	index = comFunc.BytesToUint16(context[2:4])
	constantInt.param2 = index
	return result[:], 4
}

/******************************************************************
    功能:读取Methodref_Info
	入参:文件内容
    返回值:1、转化后的码流，即常量Index值
	      2、消耗的码流数量
******************************************************************/
func readConstantMethodrefInfo(context []byte) ([]byte, uint32) {
	//实现同Fieldref_INFO
	return readConstantFieldrefInfo(context)
}

/******************************************************************
    功能:读取InterfaceMethodref_Info
	入参:文件内容
    返回值:1、转化后的码流，即常量Index值
	      2、消耗的码流数量
******************************************************************/
func readConstantInterfaceMethodrefInfo(context []byte) ([]byte, uint32) {
	//实现同Fieldref_INFO
	return readConstantFieldrefInfo(context)
}

/******************************************************************
    功能:读取NameAndType_Info
	入参:文件内容
    返回值:1、转化后的码流，即常量Index值
	      2、消耗的码流数量
******************************************************************/
func readConstantNameAndTypeInfo(context []byte) ([]byte, uint32) {
	//实现同Fieldref_INFO
	return readConstantFieldrefInfo(context)
}

/******************************************************************
    功能:读取class信息
	入参:1、文件内容
	    2、常量池
    返回值:1、可访问性标志
	      2、类名在常量池中的地址
	      3、超类的地址
	      4、error
******************************************************************/
func readClassInfo(context, constPool []byte) ([]byte, uint16, uint32, uint32, error) {

	//可访问性
	accessFlag := comFunc.BytesToUint16(context[0:2])
	//类名在常量池中为位置
	classNameIndex := uint32(comFunc.BytesToUint16(context[2:4]))
	//类名在符号表中的位置
	classSymbol, err := GetClassFromConstPool(constPool, classNameIndex)
	if err != nil {
		return nil, 0, memCtrl.INVALID_MEM, memCtrl.INVALID_MEM, err
	}
	//超类名在常量池中为位置
	superClassNameIndex := uint32(comFunc.BytesToUint16(context[4:6]))
	superClassAddr := uint32(0)
	//为0则意味着该类是Object,没有超类
	if superClassNameIndex != 0 {
		//超类名在符号表中的位置
		superClassSymbol, err := GetClassFromConstPool(constPool, superClassNameIndex)
		if err != nil {
			return nil, 0, memCtrl.INVALID_MEM, memCtrl.INVALID_MEM, err
		}
		superClassAddr = memCtrl.GetClassMemAddr(superClassSymbol)
		//如果获取不到，则说明不在内存中，需要去加载
		if superClassAddr == memCtrl.INVALID_MEM {
			//获取类名(string)
			className := string(memCtrl.GetSymbol(superClassSymbol))
			superClassAddr, err = LoadClass(className)
			if err != nil {
				return nil, 0, memCtrl.INVALID_MEM, memCtrl.INVALID_MEM, err
			}
		}
	}

	return context[6:], accessFlag, classSymbol, superClassAddr, nil
}

/******************************************************************
    功能:读取interface信息
	入参:文件内容
    返回值:1、读取后的context
	      2、解析后的类信息码流
******************************************************************/
func readInterfaces(context []byte, constPool []byte) ([]byte, uint32, []byte, error) {

	interfaceNum := comFunc.BytesToUint16(context[0:2])
	result := make([]byte, interfaceNum*4)

	//接口数量
	num := uint32(interfaceNum)

	for i := uint32(0); i < num; i++ {

		adr := (*uint32)(comFunc.BytesToUnsafePointer(result[i*4 : i*4+4]))
		index := uint32(comFunc.BytesToUint16(context[2*i+2 : 2*i+4]))
		//接口在符号表中的位置
		interfaceSymbol, err := GetClassFromConstPool(constPool, index)
		if err != nil {
			return nil, 0, nil, err
		}
		*adr = memCtrl.GetClassMemAddr(interfaceSymbol)
		//如果获取不到，则说明不在内存中，需要去加载
		if *adr == memCtrl.INVALID_MEM {
			//获取接口名(string)
			interfaceName := string(memCtrl.GetSymbol(interfaceSymbol))

			*adr, err = LoadClass(interfaceName)
			if err != nil {
				return nil, 0, nil, err
			}
		}
	}

	return context[interfaceNum*2+2:], num, result, nil
}

/******************************************************************
    功能:读取Filed信息
	入参:文件内容
    返回值:1、读取后的context
	      2、解析后的类信息码流
******************************************************************/
func readFields(context []byte) ([]byte, []byte) {
	return nil, nil
}
