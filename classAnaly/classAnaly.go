package classAnaly

import (
	"bytes"
	"errors"
	"fmt"

	"../accessOp"
	"../class"
	"../comFunc"
	"../comValue"
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
	rsv                   [2]uint8
	ConstNum              uint32 //常量数量
	FiledInfoDev          uint32 //参数信息偏移
	UnstaticParaDev       uint32 //非static参数地址
	UnstaticParaSize      uint32 //非static参数大小
	UnstaticParaTotalSize uint32 //非static参数内存总大小(即，分配实例的大小)
	StaticParaDev         uint32 //static参数地址
	StaticParaSize        uint32 //static参数大小
	StaticMem             uint32 //类实例地址
	InterfaceDev          uint32 //接口定义偏移
	InterfaceNum          uint32 //接口数量
	MethodDev             uint32 //方法定义偏移
	MethodNum             uint32 //方法数量
	LocalAdr              uint32 //该类的地址
}

const CLASS_INFO_SIZE = 16 * 4

type FILED_ITEM struct {
	FiledName    uint32 //字段名(符号表索引)
	Index        uint32 //实例(包括类实例)中的索引值,从0开始，遇到long和double则跳1
	FiledInfoDev uint32 //字段描述偏移
}

const FILED_ITEM_SIZE = 3 * 4

type FILED_INFO struct {
	AccessFlag uint16 //可访问性
	rsv        [2]uint8
	Descriptor uint32 //描述符(符号表索引)
	AttriCount uint32 //属性数量
}

const FILED_INFO_SIZE = 3 * 4

type ATTRI_INFO struct {
	AttriName uint32 //属性名(符号表中的地址)
	Length    uint32 //长度
}

const ATTRI_INFO_SIZE = 8

type CONST_PAIR struct {
	StaticFiledIndex uint32 //static字段索引
	ConstIndex       uint32 //常量索引
	IsLongOrDouble   bool
}

type METHOD struct {
	AccessFlag uint16 //可访问属性
	rsv        [2]uint8
	MethodName uint32 //方法名
	Descriptor uint32 //描述符
	CodeAddr   uint32 //code地址,code属性里没有属性和长度，直接就是Code结构体开始
	Attribute  uint32 //属性地址
	AttriNum   uint32 //属性数量
}

const METHOD_SIZE = 6 * 4

type CODE_ATTRI struct {
	MaxStack       uint32 //方法栈
	MaxLocal       uint32 //局部变量大小
	CodeLength     uint32
	ExceptionCount uint32
	AttriNum       uint32
}

const CODE_ATTRI_SIZE = 20

var magicNum = []byte{0xCA, 0xFE, 0xBA, 0xBE}

const FILED_ACC_PUBLIC = 0x0001
const FILED_ACC_PRIVATE = 0x0002
const FILED_ACC_PROTECTED = 0x0004
const FILED_ACC_STATIC = 0x0008
const FILED_ACC_FINAL = 0x0010
const FILED_ACC_VOILATIE = 0x0040
const FILED_ACC_TRANSIENT = 0x0080
const FILED_ACC_SYNTHETIC = 0x1000
const FILED_ACC_ENUM = 0x4000

const CLASS_ACC_PUBLIC = 0x0001
const CLASS_ACC_FINAL = 0x0010
const CLASS_ACC_SUPER = 0x0020 //必选
const CLASS_ACC_INTERFACE = 0x0200
const CLASS_ACC_ABSTRACT = 0x0400
const CLASS_ACC_SYNTHETIC = 0x1000
const CLASS_ACC_ANNOTATION = 0x2000
const CLASS_ACC_ENUM = 0x4000

/******************************************************************
    功能:加载类
	入参:类名
    返回值:1、类地址
	      2、error
******************************************************************/
func LoadClass(className string) (*CLASS_INFO, error) {

	//读取字节码文件内容
	context, err := class.ReadClass(className)
	if err != nil {
		return nil, err
	}

	//解析结果定义
	classInfoMem := make([]byte, CLASS_INFO_SIZE)
	result := make([]byte, CLASS_INFO_SIZE)

	//Class Info信息定义
	classInfo := (*CLASS_INFO)(comFunc.BytesToUnsafePointer(classInfoMem[0:]))

	//读取魔数
	context, err = readMagicNum(context)
	if err != nil {
		return nil, err
	}

	//读取版本号，暂时不使用
	context, _, _ = readVersion(context)
	//读取常量池
	context, constPool, num, err := readConstantPool(context)
	classInfo.ConstNum = num
	if err != nil {
		return nil, err
	}
	result = append(result, constPool...)

	//读取类信息
	context, classInfo.AccessFlag, classInfo.ClassName, classInfo.SuperClassAddr, err = readClassInfo(context, constPool)
	if err != nil {
		return nil, err
	}

	//读取接口信息
	classInfo.InterfaceDev = uint32(len(result))
	context, num, interfaceInfo, err := readInterfaces(context, constPool)
	if err != nil {
		return nil, err
	}
	classInfo.InterfaceNum = num
	result = append(result, interfaceInfo...)

	//读取字段信息
	context, attriInfo, unstatic, static, unstaticSize, staticSize, constPair, err := readFields(context, constPool)
	classInfo.FiledInfoDev = uint32(len(result))

	//刷新字段的偏移信息
	for i := uint32(0); i < uint32(len(unstatic)); i += FILED_ITEM_SIZE {
		fileditem := (*FILED_ITEM)(comFunc.BytesToUnsafePointer(unstatic[i : i+FILED_ITEM_SIZE]))
		fileditem.FiledInfoDev += classInfo.FiledInfoDev
	}
	for i := uint32(0); i < uint32(len(static)); i += FILED_ITEM_SIZE {
		fileditem := (*FILED_ITEM)(comFunc.BytesToUnsafePointer(static[i : i+FILED_ITEM_SIZE]))
		fileditem.FiledInfoDev += classInfo.FiledInfoDev
	}
	result = append(result, attriInfo...)
	//静态字段
	classInfo.StaticParaDev = uint32(len(result))
	classInfo.StaticParaSize = staticSize * 4
	result = append(result, static...)

	//非静态字段
	classInfo.UnstaticParaDev = uint32(len(result))
	classInfo.UnstaticParaSize = unstaticSize * 4
	result = append(result, unstatic...)

	//静态常量初始化
	if staticSize != 0 {
		clazInstAdr, err := memCtrl.Malloc(classInfo.StaticParaSize, memCtrl.CLASS_INSTANCE_NODE)
		classInfo.StaticMem = clazInstAdr
		if err != nil {
			return nil, err
		}
		for _, pair := range constPair {
			v := (*uint32)(memCtrl.GetPointer(clazInstAdr+pair.StaticFiledIndex*4, 4))
			*v = GetUint32FromConstPool(constPool, pair.ConstIndex)
			if err != nil {
				return nil, err
			}
			if pair.IsLongOrDouble {
				v := (*uint32)(memCtrl.GetPointer(clazInstAdr+pair.StaticFiledIndex*4+4, 4))
				*v = GetUint32FromConstPool(constPool, pair.ConstIndex)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	//UnstaticParaTotalSize字段计算
	superClass := classInfo.GetSuperClass()
	if superClass == nil {
		classInfo.UnstaticParaTotalSize = classInfo.UnstaticParaSize
	} else {
		classInfo.UnstaticParaTotalSize = classInfo.UnstaticParaSize + superClass.UnstaticParaTotalSize
	}

	//读取Method
	context, methods, attris, methodNum, err := readMethods(context, constPool)
	if err != nil {
		return nil, err
	}
	classInfo.MethodDev = uint32(len(result))
	classInfo.MethodNum = methodNum

	//刷新属性的偏移值
	attriDev := uint32(len(methods) + len(result))
	methodInfos := *(*[]METHOD)(comFunc.BytesToArray(methods, METHOD_SIZE))
	for i := 0; i < len(methodInfos); i++ {
		methodInfos[i].Attribute += attriDev
		methodInfos[i].CodeAddr += attriDev
	}
	result = append(result, methods...)
	result = append(result, attris...)

	//属性暂不解析
	result = append(result, context...)

	copy(result[0:CLASS_INFO_SIZE], classInfoMem)
	//保存到内存中
	memAdr, err := memCtrl.PutClass(classInfo.ClassName, result)
	classInfo = (*CLASS_INFO)(memCtrl.GetPointer(memAdr, CLASS_INFO_SIZE))
	classInfo.LocalAdr = memAdr

	if err != nil {
		return nil, err
	}
	//to do,执行static代码块
	return classInfo, nil
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
	strs := make([]uint32, 0)
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
			strs = append(strs, uint32(len(result)))
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
	//将String常量的值换成字符串地址(字符串常量池中的地址)
	for _, v := range strs {
		str := (*CONSTANT_TYPE_32)(comFunc.BytesToUnsafePointer(result[v : v+4]))
		strAdr := GetUtf8FromConstPool(result, str.param)
		if err != nil {
			return nil, nil, 0, err
		}
		str.param, err = access.PutString(access.BytesToUint16(memCtrl.GetSymbol(strAdr)))
		if err != nil {
			return nil, nil, 0, err
		}
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
	classSymbol := GetClassNameFromConstPool(constPool, classNameIndex)

	//超类名在常量池中为位置
	superClassNameIndex := uint32(comFunc.BytesToUint16(context[4:6]))
	var superClassAdr uint32 = memCtrl.INVALID_MEM
	//为0则意味着该类是Object,没有超类
	if superClassNameIndex != 0 {
		//超类名在符号表中的位置
		superClassSymbol := GetClassNameFromConstPool(constPool, superClassNameIndex)

		superClassAdr = memCtrl.GetClassMemAddr(superClassSymbol)
		//如果获取不到，则说明不在内存中，需要去加载
		if superClassAdr == memCtrl.INVALID_MEM {
			//获取类名(string)
			className := string(memCtrl.GetSymbol(superClassSymbol))
			superClass, err := LoadClass(className)
			if err != nil {
				return nil, 0, memCtrl.INVALID_MEM, memCtrl.INVALID_MEM, err
			}
			superClassAdr = superClass.LocalAdr
		}
	}
	return context[6:], accessFlag, classSymbol, superClassAdr, nil
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
		interfaceSymbol := GetClassNameFromConstPool(constPool, index)

		*adr = memCtrl.GetClassMemAddr(interfaceSymbol)
		//如果获取不到，则说明不在内存中，需要去加载
		if *adr == memCtrl.INVALID_MEM {
			//获取接口名(string)
			interfaceName := string(memCtrl.GetSymbol(interfaceSymbol))

			superClass, err := LoadClass(interfaceName)
			if err != nil {
				return nil, 0, nil, err
			}
			*adr = superClass.LocalAdr
		}
	}

	return context[interfaceNum*2+2:], num, result, nil
}

/******************************************************************
    功能:读取Filed信息
	入参:1、文件内容
	    2、常量池
    返回值:1、读取后的context
          2、字段描述
          3、非静态字段信息
          4、静态字段信息
          5、非静态字段大小
          6、静态字段大小
          7、常量对
          8、error
******************************************************************/
func readFields(context, constPool []byte) ([]byte, []byte, []byte, []byte, uint32, uint32, []CONST_PAIR, error) {
	unstaticNum := uint32(0)
	staticNum := uint32(0)

	//字段数量
	filedNum := uint32(comFunc.BytesToUint16(context[0:2]))
	context = context[2:]

	//字段信息
	filedInfos := make([]byte, 0)
	//非静态字段
	unstaticFileds := make([]byte, 0)
	//静态字段
	staticFileds := make([]byte, 0)
	//常量对
	constPairs := make([]CONST_PAIR, 0)

	constSymbol, err := memCtrl.PutSymbol([]byte("ConstantValue"))
	if err != nil {
		return nil, nil, nil, nil, 0, 0, nil, err
	}
	for i := uint32(0); i < filedNum; i++ {
		filed := make([]byte, FILED_INFO_SIZE)
		filedInfo := (*FILED_INFO)(comFunc.BytesToUnsafePointer(filed))
		//可访问性
		filedInfo.AccessFlag = comFunc.BytesToUint16(context[0:2])
		//字段名
		filedName := GetUtf8FromConstPool(constPool, uint32(comFunc.BytesToUint16(context[2:4])))

		//描述符
		filedInfo.Descriptor = GetUtf8FromConstPool(constPool, uint32(comFunc.BytesToUint16(context[4:6])))

		//属性数量
		filedInfo.AttriCount = uint32(comFunc.BytesToUint16(context[6:8]))

		context = context[8:]
		for j := uint32(0); j < filedInfo.AttriCount; j++ {
			attriMem := make([]byte, ATTRI_INFO_SIZE)
			attri := (*ATTRI_INFO)(comFunc.BytesToUnsafePointer(attriMem))
			//属性名
			attri.AttriName = GetUtf8FromConstPool(constPool, uint32(comFunc.BytesToUint16(context[0:2])))

			//属性长度
			attri.Length = comFunc.BytesToUint32(context[2:6])

			//判断是否是常量属性
			//静态常量的处理
			if attri.AttriName == constSymbol &&
				(filedInfo.AccessFlag&FILED_ACC_STATIC == FILED_ACC_STATIC) &&
				(filedInfo.AccessFlag&FILED_ACC_FINAL == FILED_ACC_FINAL) {
				constValue := uint32(comFunc.BytesToUint16(context[6:8]))

				if memCtrl.SYM_J == filedInfo.Descriptor ||
					memCtrl.SYM_D == filedInfo.Descriptor {
					constPairs = append(constPairs, CONST_PAIR{staticNum, constValue, true})
				} else {
					constPairs = append(constPairs, CONST_PAIR{staticNum, constValue, false})
				}
			}
			filed = append(filed, attriMem...)
			//属性内容暂不解析
			filed = append(filed, context[6:6+attri.Length]...)
			context = context[6+attri.Length:]

		}

		//字段Item
		filedMem := make([]byte, FILED_ITEM_SIZE)
		filedItem := (*FILED_ITEM)(comFunc.BytesToUnsafePointer(filedMem))
		filedItem.FiledInfoDev = uint32(len(filedInfos))
		filedItem.FiledName = filedName

		//判断是否是long或double
		if filedInfo.AccessFlag&FILED_ACC_STATIC == FILED_ACC_STATIC {
			//静态字段的处理
			filedItem.Index = staticNum
			staticFileds = append(staticFileds, filedMem...)
			staticNum++
			if memCtrl.SYM_J == filedInfo.Descriptor ||
				memCtrl.SYM_D == filedInfo.Descriptor {
				staticNum++
			}
		} else {
			//非静态字段的处理
			filedItem.Index = unstaticNum
			unstaticFileds = append(unstaticFileds, filedMem...)
			unstaticNum++
			if memCtrl.SYM_J == filedInfo.Descriptor ||
				memCtrl.SYM_D == filedInfo.Descriptor {
				unstaticNum++
			}
		}
		filedInfos = append(filedInfos, filed...)
	}
	return context, filedInfos, unstaticFileds, staticFileds, unstaticNum, staticNum, constPairs, nil
}

/******************************************************************
    功能:读取Method信息
	入参:1、文件内容
	    2、常量池
    返回值:1、读取后的context
          2、方法描述
          3、属性描述
          4、方法数量
          5、error
******************************************************************/
func readMethods(context, constPool []byte) ([]byte, []byte, []byte, uint32, error) {
	//方法数量
	methodsNum := uint32(comFunc.BytesToUint16(context[0:2]))
	context = context[2:]
	methods := make([]byte, 0)
	attris := make([]byte, 0)
	var err error
	for i := uint32(0); i < methodsNum; i++ {
		methodInfo_mem := make([]byte, METHOD_SIZE)
		methodInfo := (*METHOD)(comFunc.BytesToUnsafePointer(methodInfo_mem))
		methodInfo.CodeAddr = memCtrl.INVALID_MEM
		//方法可访问性
		methodInfo.AccessFlag = comFunc.BytesToUint16(context[0:2])
		//方法名
		methodInfo.MethodName = GetUtf8FromConstPool(constPool, uint32(comFunc.BytesToUint16(context[2:4])))
		if err != nil {
			return nil, nil, nil, 0, err
		}
		//方法描述符
		methodInfo.Descriptor = GetUtf8FromConstPool(constPool, uint32(comFunc.BytesToUint16(context[4:6])))

		//属性数量
		attriNum := uint32(comFunc.BytesToUint16(context[6:8]))
		methodInfo.AttriNum = attriNum
		//属性地址
		methodInfo.Attribute = uint32(len(attris))
		context = context[8:]
		//Code符号表中的地址
		codeSymbol, err := memCtrl.PutSymbol([]byte("Code"))
		if err != nil {
			return nil, nil, nil, 0, err
		}
		for j := uint32(0); j < attriNum; j++ {
			attri_mem := make([]byte, ATTRI_INFO_SIZE)
			attri := (*ATTRI_INFO)(comFunc.BytesToUnsafePointer(attri_mem))
			//属性名
			attri.AttriName = GetUtf8FromConstPool(constPool, uint32(comFunc.BytesToUint16(context[0:2])))

			if attri.AttriName == codeSymbol {
				//Code属性的处理
				methodInfo.CodeAddr = uint32(len(attris))
				attriLength := uint32(comFunc.BytesToUint32(context[2:6]))
				codeMem := readCode(context[uint32(6) : uint32(6)+attriLength])
				attris = append(attris, codeMem...)
				context = context[uint32(6)+attriLength:]
			} else {
				//非Code属性的处理
				attri.Length = uint32(comFunc.BytesToUint32(context[2:6]))
				attris = append(attris, attri_mem...)
				attris = append(attris, context[6:uint32(6)+attri.Length]...)
				context = context[uint32(6)+attri.Length:]
			}
		}
		methods = append(methods, methodInfo_mem...)
	}
	return context, methods, attris, methodsNum, nil
}

/******************************************************************
    功能:读取Code信息
	入参:1、文件内容

    返回值:1、格式化后的字节码
******************************************************************/
func readCode(context []byte) []byte {
	code_mem := make([]byte, CODE_ATTRI_SIZE)
	code := (*CODE_ATTRI)(comFunc.BytesToUnsafePointer(code_mem))
	code.MaxStack = uint32(comFunc.BytesToUint16(context[0:2]))
	code.MaxLocal = uint32(comFunc.BytesToUint16(context[2:4]))
	code.CodeLength = uint32(comFunc.BytesToUint32(context[4:8]))
	context = context[8:]
	codeOp := make([]byte, 0, code.CodeLength)
	for i := uint32(0); i < code.CodeLength; {
		op := context[i]
		codeOp = append(codeOp, op)
		i++
		switch op {
		case comValue.NOP:
		case comValue.ACONST_NULL:
		case comValue.ICONST_M1:
		case comValue.ICONST_0:
		case comValue.ICONST_1:
		case comValue.ICONST_2:
		case comValue.ICONST_3:
		case comValue.ICONST_4:
		case comValue.ICONST_5:
		case comValue.LCONST_0:
		case comValue.LCONST_1:
		case comValue.FCONST_0:
		case comValue.FCONST_1:
		case comValue.FCONST_2:
		case comValue.DCONST_0:
		case comValue.DCONST_1:
		case comValue.ILOAD_0:
		case comValue.ILOAD_1:
		case comValue.ILOAD_2:
		case comValue.ILOAD_3:
		case comValue.LLOAD_0:
		case comValue.LLOAD_1:
		case comValue.LLOAD_2:
		case comValue.LLOAD_3:
		case comValue.FLOAD_0:
		case comValue.FLOAD_1:
		case comValue.FLOAD_2:
		case comValue.FLOAD_3:
		case comValue.DLOAD_0:
		case comValue.DLOAD_1:
		case comValue.DLOAD_2:
		case comValue.DLOAD_3:
		case comValue.ALOAD_0:
		case comValue.ALOAD_1:
		case comValue.ALOAD_2:
		case comValue.ALOAD_3:
		case comValue.IALOAD:
		case comValue.LALOAD:
		case comValue.FALOAD:
		case comValue.DALOAD:
		case comValue.AALOAD:
		case comValue.BALOAD:
		case comValue.CALOAD:
		case comValue.SALOAD:
		case comValue.ISTORE_0:
		case comValue.ISTORE_1:
		case comValue.ISTORE_2:
		case comValue.ISTORE_3:
		case comValue.LSTORE_0:
		case comValue.LSTORE_1:
		case comValue.LSTORE_2:
		case comValue.LSTORE_3:
		case comValue.FSTORE_0:
		case comValue.FSTORE_1:
		case comValue.FSTORE_2:
		case comValue.FSTORE_3:
		case comValue.DSTORE_0:
		case comValue.DSTORE_1:
		case comValue.DSTORE_2:
		case comValue.DSTORE_3:
		case comValue.ASTORE_0:
		case comValue.ASTORE_1:
		case comValue.ASTORE_2:
		case comValue.ASTORE_3:
		case comValue.IASTORE:
		case comValue.LASTORE:
		case comValue.FASTORE:
		case comValue.DASTORE:
		case comValue.AASTORE:
		case comValue.BASTORE:
		case comValue.CASTORE:
		case comValue.SASTORE:
		case comValue.POP:
		case comValue.POP2:
		case comValue.DUP:
		case comValue.DUP_X1:
		case comValue.DUP_X2:
		case comValue.DUP2:
		case comValue.DUP2_X1:
		case comValue.DUP2_X2:
		case comValue.SWAP:
		case comValue.IADD:
		case comValue.LADD:
		case comValue.FADD:
		case comValue.DADD:
		case comValue.ISUB:
		case comValue.LSUB:
		case comValue.FSUB:
		case comValue.DSUB:
		case comValue.IMUL:
		case comValue.LMUL:
		case comValue.FMUL:
		case comValue.DMUL:
		case comValue.IDIV:
		case comValue.LDIV:
		case comValue.FDIV:
		case comValue.DDIV:
		case comValue.IREM:
		case comValue.LREM:
		case comValue.FREM:
		case comValue.DREM:
		case comValue.INEG:
		case comValue.LENG:
		case comValue.FNEG:
		case comValue.DNEG:
		case comValue.ISHL:
		case comValue.LSHL:
		case comValue.ISHR:
		case comValue.LSHR:
		case comValue.IUSHR:
		case comValue.LUSHR:
		case comValue.IAND:
		case comValue.LAND:
		case comValue.IOR:
		case comValue.LOR:
		case comValue.IXOR:
		case comValue.LXOR:
		case comValue.I2L:
		case comValue.I2F:
		case comValue.I2D:
		case comValue.L2I:
		case comValue.L2F:
		case comValue.L2D:
		case comValue.F2I:
		case comValue.F2L:
		case comValue.F2D:
		case comValue.D2I:
		case comValue.D2L:
		case comValue.D2F:
		case comValue.I2B:
		case comValue.I2C:
		case comValue.I2S:
		case comValue.LCMP:
		case comValue.FCMPL:
		case comValue.FCMPG:
		case comValue.DCMPL:
		case comValue.DCMPG:
		case comValue.IRETURN:
		case comValue.LRETURN:
		case comValue.FRETURN:
		case comValue.DRETURN:
		case comValue.ARETURN:
		case comValue.RETURN:
		case comValue.ARRAYLENGTH:
		case comValue.ATHROW:
		case comValue.MONITORENTER:
		case comValue.MONITOREXIT:
			//无操作数
		case comValue.BIPUSH:
			fallthrough
		case comValue.LDC:
			fallthrough
		case comValue.ILOAD:
			fallthrough
		case comValue.LLOAD:
			fallthrough
		case comValue.FLOAD:
			fallthrough
		case comValue.DLOAD:
			fallthrough
		case comValue.ALOAD:
			fallthrough
		case comValue.ISTORE:
			fallthrough
		case comValue.LSTORE:
			fallthrough
		case comValue.FSTORE:
			fallthrough
		case comValue.DSTORE:
			fallthrough
		case comValue.ASTORE:
			fallthrough
		case comValue.NEWARRAY:
			codeOp = append(codeOp, context[i])
			i++

		case comValue.SIPUSH:
			b := [2]byte{}
			s := (*int16)(comFunc.BytesToUnsafePointer(b[:]))
			*s = comFunc.BytesToInt16(context[i : i+2])
			codeOp = append(codeOp, b[:]...)
			i += 2

		case comValue.LDC_W:
			fallthrough
		case comValue.LDC2_W:
			fallthrough
		case comValue.IFEQ:
			fallthrough
		case comValue.IFNE:
			fallthrough
		case comValue.IFLT:
			fallthrough
		case comValue.IFGE:
			fallthrough
		case comValue.IFGT:
			fallthrough
		case comValue.IFLE:
			fallthrough
		case comValue.IF_ICMPEQ:
			fallthrough
		case comValue.IF_ICMPNE:
			fallthrough
		case comValue.IF_ICMPLT:
			fallthrough
		case comValue.IF_ICMPGE:
			fallthrough
		case comValue.IF_ICMPGT:
			fallthrough
		case comValue.IF_ICMPLE:
			fallthrough
		case comValue.IF_ACMPEQ:
			fallthrough
		case comValue.IF_ACMPNE:
			fallthrough
		case comValue.GOTO:
			fallthrough
		case comValue.GETSTATIC:
			fallthrough
		case comValue.PUTSTATIC:
			fallthrough
		case comValue.GETFIELD:
			fallthrough
		case comValue.PUTFIELD:
			fallthrough
		case comValue.INVOKEVIRTUAL:
			fallthrough
		case comValue.INVOKESPECIAL:
			fallthrough
		case comValue.INVOKESTATIC:
			fallthrough
		case comValue.INVOKEINTERFACE:
			fallthrough
		case comValue.NEW:
			fallthrough
		case comValue.ANEWARRAY:
			fallthrough
		case comValue.CHECKCAST:
			fallthrough
		case comValue.INSTANCEOF:
			fallthrough
		case comValue.IFNULL:
			fallthrough
		case comValue.IFNONNULL:
			b := [2]byte{}
			s := (*uint16)(comFunc.BytesToUnsafePointer(b[:]))
			*s = comFunc.BytesToUint16(context[i : i+2])
			codeOp = append(codeOp, b[:]...)
			i += 2

		case comValue.IINC:
			codeOp = append(codeOp, context[i], context[i+1])
			i += 2

		case comValue.TABLESWITCH:
			//3个填充位
			codeOp = append(codeOp, context[i:i+3]...)
			i += 3
			//default
			def := [4]byte{}
			s := (*uint32)(comFunc.BytesToUnsafePointer(def[:]))
			*s = comFunc.BytesToUint32(context[i : i+4])
			codeOp = append(codeOp, def[:]...)
			i += 4
			//low
			low_b := [4]byte{}
			low := (*uint32)(comFunc.BytesToUnsafePointer(low_b[:]))
			*low = comFunc.BytesToUint32(context[i : i+4])
			codeOp = append(codeOp, low_b[:]...)
			i += 4
			//high
			high_b := [4]byte{}
			high := (*uint32)(comFunc.BytesToUnsafePointer(high_b[:]))
			*high = comFunc.BytesToUint32(context[i : i+4])
			codeOp = append(codeOp, high_b[:]...)
			i += 4
			//offset
			for j := *low; j <= *high; j++ {
				offset_b := [4]byte{}
				offset := (*uint32)(comFunc.BytesToUnsafePointer(offset_b[:]))
				*offset = comFunc.BytesToUint32(context[i : i+4])
				codeOp = append(codeOp, offset_b[:]...)
				i += 4
			}

		case comValue.LOOKUPSWITCH:
			//3个填充位
			codeOp = append(codeOp, context[i:i+3]...)
			i += 3
			//default
			def := [4]byte{}
			s := (*uint32)(comFunc.BytesToUnsafePointer(def[:]))
			*s = comFunc.BytesToUint32(context[i : i+4])
			codeOp = append(codeOp, def[:]...)
			i += 4
			//pair的数量
			pairs_b := [4]byte{}
			pairs := (*uint32)(comFunc.BytesToUnsafePointer(pairs_b[:]))
			*pairs = comFunc.BytesToUint32(context[i : i+4])
			codeOp = append(codeOp, pairs_b[:]...)
			i += 4
			//各pair
			for j := uint32(0); j <= (*pairs)*2; j++ {
				pair_b := [4]byte{}
				pair := (*uint32)(comFunc.BytesToUnsafePointer(pair_b[:]))
				*pair = comFunc.BytesToUint32(context[i : i+4])
				codeOp = append(codeOp, pair_b[:]...)
				i += 4
			}
		case comValue.WIDE:
			op := context[i]
			codeOp = append(codeOp, op)
			i++
			switch op {
			case comValue.ILOAD:
				fallthrough
			case comValue.LLOAD:
				fallthrough
			case comValue.FLOAD:
				fallthrough
			case comValue.DLOAD:
				fallthrough
			case comValue.ALOAD:
				fallthrough
			case comValue.ISTORE:
				fallthrough
			case comValue.LSTORE:
				fallthrough
			case comValue.FSTORE:
				fallthrough
			case comValue.DSTORE:
				fallthrough
			case comValue.ASTORE:
				b := [2]byte{}
				s := (*uint16)(comFunc.BytesToUnsafePointer(b[:]))
				*s = comFunc.BytesToUint16(context[i : i+2])
				codeOp = append(codeOp, b[:]...)
				i += 2
			case comValue.IINC:
				b := [2]byte{}
				s := (*uint16)(comFunc.BytesToUnsafePointer(b[:]))
				*s = comFunc.BytesToUint16(context[i : i+2])
				codeOp = append(codeOp, b[:]...)
				codeOp = append(codeOp, context[i+2])
				i += 3
			case comValue.RET:
				fmt.Println(op)
				panic("readCode():该指令不支持（wide)")
			default:
				fmt.Println(op)
				panic("readCode():该指令不支持wide")
			}

		case comValue.MULTIANEWARRAY:
			b := [2]byte{}
			s := (*uint16)(comFunc.BytesToUnsafePointer(b[:]))
			*s = comFunc.BytesToUint16(context[i : i+2])
			codeOp = append(codeOp, b[:]...)
			codeOp = append(codeOp, context[i+2])
			i += 3

		case comValue.GOTO_W:
			b := [4]byte{}
			s := (*uint32)(comFunc.BytesToUnsafePointer(b[:]))
			*s = comFunc.BytesToUint32(context[i : i+4])
			codeOp = append(codeOp, b[:]...)
			i += 4

		case comValue.JSR:
			fallthrough
		case comValue.RET:
			fallthrough
		case comValue.JSR_W:
			fmt.Println(op)
			panic("readCode():该指令不支持")
		default:
			fmt.Println(op)
			panic("readCode():该指令不存在")
		}
	}

	//异常处理表，code属性(),暂缺
	code_mem = append(code_mem, codeOp...)
	return append(code_mem, context[code.CodeLength:]...)
}
