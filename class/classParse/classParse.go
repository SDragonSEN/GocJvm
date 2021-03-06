package classParse

import (
	"bytes"
	"errors"
	"fmt"

	. "access/string"
	. "basic/com"
	. "basic/memCtrl"
	. "basic/symbol"
	. "class/classFind"
	. "class/classInterface"
	. "class/classTable"
)

/******************************************************************
    功能:加载类
	入参:类名
    返回值:1、类地址
	      2、error
******************************************************************/
func LoadClass(className string) (*CLASS_INFO, error) {

	//读取字节码文件内容
	context, err := ReadClass(className)
	if err != nil {
		return nil, err
	}

	//解析结果定义
	classInfoMem := make([]byte, CLASS_INFO_SIZE)
	result := make([]byte, CLASS_INFO_SIZE)

	//Class Info信息定义
	classInfo := (*CLASS_INFO)(BytesToUnsafePointer(classInfoMem[0:]))
	classInfo.IsCInit = false
	//读取魔数
	context, err = readMagicNum(context)
	if err != nil {
		return nil, err
	}

	//读取版本号，暂时不使用
	context, _, _ = readVersion(context)
	//读取常量池
	context, constPool, clazConst, num, err := readConstantPool(context)
	classInfo.ConstNum = num
	if err != nil {
		return nil, err
	}
	result = append(result, constPool...)
	classInfo.ClassConstDev = uint32(len(result))
	result = append(result, clazConst...)
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
		fileditem := (*FILED_ITEM)(BytesToUnsafePointer(unstatic[i : i+FILED_ITEM_SIZE]))
		fileditem.FiledInfoDev += classInfo.FiledInfoDev
	}
	for i := uint32(0); i < uint32(len(static)); i += FILED_ITEM_SIZE {
		fileditem := (*FILED_ITEM)(BytesToUnsafePointer(static[i : i+FILED_ITEM_SIZE]))
		fileditem.FiledInfoDev += classInfo.FiledInfoDev
	}
	result = append(result, attriInfo...)
	//静态字段
	classInfo.StaticParaDev = uint32(len(result))
	classInfo.StaticParaSize = staticSize * 4
	classInfo.StaticParaNum = uint32(len(static) / FILED_ITEM_SIZE)
	result = append(result, static...)

	//非静态字段
	classInfo.UnstaticParaDev = uint32(len(result))
	classInfo.UnstaticParaSize = unstaticSize * 4
	classInfo.UnstaticParaNum = uint32(len(unstatic) / FILED_ITEM_SIZE)
	result = append(result, unstatic...)

	//静态常量初始化
	if staticSize != 0 {
		clazInstAdr, err := Malloc(classInfo.StaticParaSize, CLASS_INSTANCE_NODE)
		classInfo.StaticMem = clazInstAdr
		if err != nil {
			return nil, err
		}
		for _, pair := range constPair {
			v := (*uint32)(GetPointer(clazInstAdr+pair.StaticFiledIndex*4, 4))
			*v = GetUint32FromConstPool(constPool, pair.ConstIndex)
			if err != nil {
				return nil, err
			}
			if pair.IsLongOrDouble {
				v := (*uint32)(GetPointer(clazInstAdr+pair.StaticFiledIndex*4+4, 4))
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
	methodInfos := *(*[]METHOD)(BytesToArray(methods, METHOD_SIZE))
	for i := 0; i < len(methodInfos); i++ {
		methodInfos[i].Attribute += attriDev
		if methodInfos[i].CodeAddr != INVALID_MEM {
			methodInfos[i].CodeAddr += attriDev
		}
	}
	result = append(result, methods...)
	result = append(result, attris...)

	//属性暂不解析
	result = append(result, context...)

	copy(result[0:CLASS_INFO_SIZE], classInfoMem)
	//保存到内存中
	memAdr, err := PutClass(classInfo.ClassName, result)
	classInfo = (*CLASS_INFO)(GetPointer(memAdr, CLASS_INFO_SIZE))
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
	if bytes.Compare(context[0:4], MagicNum) != 0 {
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
	minor_version = BytesToUint16(context[0:2])
	major_version = BytesToUint16(context[2:4])
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
func readConstantPool(context []byte) ([]byte, []byte, []byte, uint32, error) {
	//符号数量从1到size-1
	size := BytesToUint16(context[0:2])
	//消耗的码流数量
	var count uint32
	count = 2
	//结果
	result := make([]byte, 0)
	strs := make([]uint32, 0)
	clzs := make([]byte, 0)
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
				return nil, nil, nil, 0, err
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
			clzIndexMem := make([]byte, 4)
			p := (*uint32)(BytesToUnsafePointer(clzIndexMem))
			*p = uint32(i)
			clzs = append(clzs, clzIndexMem...)
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
			return nil, nil, nil, 0, errors.New("常量池解析错误")
		}
		count += consume
		result = append(result, constantBytes...)
	}
	//将String常量的值换成字符串地址(字符串常量池中的地址)
	for _, v := range strs {
		str := (*CONSTANT_TYPE_32)(BytesToUnsafePointer(result[v : v+4]))
		strAdr := GetUtf8FromConstPool(result, str.Param)
		if err != nil {
			return nil, nil, nil, 0, err
		}
		str.Param, err = PutString(BytesToUtf16(GetSymbol(strAdr)))
		if err != nil {
			return nil, nil, nil, 0, err
		}
	}
	return context[count:], result, clzs, uint32(size), nil
}

/******************************************************************
    功能:读取UTF8_INFO
	入参:文件内容
    返回值:1、转化后的码流，即符号表中的地址
	      2、消耗的码流数量
******************************************************************/
func readConstantUtf8Info(context []byte) ([]byte, uint32, error) {
	//获取utf8长度
	length := BytesToUint16(context[0:2])
	var count uint32
	var err error
	count = 2
	result := [4]byte{}
	constant32 := (*CONSTANT_TYPE_32)(BytesToUnsafePointer(result[:]))
	//将utf8码流加到符号表中
	constant32.Param, err = PutSymbol(context[count : count+uint32(length)])
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
	integer := BytesToUint32(context[0:4])
	result := [4]byte{}
	constantInt := (*CONSTANT_TYPE_32)(BytesToUnsafePointer(result[:]))
	constantInt.Param = integer
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
	long := BytesToUint32(context[0:4])

	constant64 := (*CONSTANT_TYPE_32)(BytesToUnsafePointer(result[0:4]))
	constant64.Param = long
	//获取低位
	long = BytesToUint32(context[4:8])
	constant64 = (*CONSTANT_TYPE_32)(BytesToUnsafePointer(result[4:8]))
	constant64.Param = long
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
	index := BytesToUint16(context[0:2])
	result := [4]byte{}
	constantInt := (*CONSTANT_TYPE_32)(BytesToUnsafePointer(result[:]))
	constantInt.Param = uint32(index)
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
	constantInt := (*CONSTANT_TYPE_16)(BytesToUnsafePointer(result[:]))
	//获取class index值
	index := BytesToUint16(context[0:2])
	constantInt.Param1 = index
	//获取name and type index值
	index = BytesToUint16(context[2:4])
	constantInt.Param2 = index
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
	accessFlag := BytesToUint16(context[0:2])
	//类名在常量池中为位置
	classNameIndex := uint32(BytesToUint16(context[2:4]))
	//类名在符号表中的位置
	classSymbol := GetClassNameFromConstPool(constPool, classNameIndex)

	//超类名在常量池中为位置
	superClassNameIndex := uint32(BytesToUint16(context[4:6]))
	var superClassAdr uint32 = INVALID_MEM
	//为0则意味着该类是Object,没有超类
	if superClassNameIndex != 0 {
		//超类名在符号表中的位置
		superClassSymbol := GetClassNameFromConstPool(constPool, superClassNameIndex)

		superClassAdr = GetClassMemAddr(superClassSymbol)
		//如果获取不到，则说明不在内存中，需要去加载
		if superClassAdr == INVALID_MEM {
			//获取类名(string)
			className := string(GetSymbol(superClassSymbol))
			superClass, err := LoadClass(className)
			if err != nil {
				return nil, 0, INVALID_MEM, INVALID_MEM, err
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

	interfaceNum := BytesToUint16(context[0:2])
	result := make([]byte, interfaceNum*4)

	//接口数量
	num := uint32(interfaceNum)

	for i := uint32(0); i < num; i++ {

		adr := (*uint32)(BytesToUnsafePointer(result[i*4 : i*4+4]))
		index := uint32(BytesToUint16(context[2*i+2 : 2*i+4]))
		//接口在符号表中的位置
		interfaceSymbol := GetClassNameFromConstPool(constPool, index)

		*adr = GetClassMemAddr(interfaceSymbol)
		//如果获取不到，则说明不在内存中，需要去加载
		if *adr == INVALID_MEM {
			//获取接口名(string)
			interfaceName := string(GetSymbol(interfaceSymbol))

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
	filedNum := uint32(BytesToUint16(context[0:2]))
	context = context[2:]

	//字段信息
	filedInfos := make([]byte, 0)
	//非静态字段
	unstaticFileds := make([]byte, 0)
	//静态字段
	staticFileds := make([]byte, 0)
	//常量对
	constPairs := make([]CONST_PAIR, 0)

	constSymbol, err := PutSymbol([]byte("ConstantValue"))
	if err != nil {
		return nil, nil, nil, nil, 0, 0, nil, err
	}
	for i := uint32(0); i < filedNum; i++ {
		filed := make([]byte, FILED_INFO_SIZE)
		filedInfo := (*FILED_INFO)(BytesToUnsafePointer(filed))
		//可访问性
		filedInfo.AccessFlag = BytesToUint16(context[0:2])
		//字段名
		filedName := GetUtf8FromConstPool(constPool, uint32(BytesToUint16(context[2:4])))

		//描述符
		filedInfo.Descriptor = GetUtf8FromConstPool(constPool, uint32(BytesToUint16(context[4:6])))

		//属性数量
		filedInfo.AttriCount = uint32(BytesToUint16(context[6:8]))

		context = context[8:]
		for j := uint32(0); j < filedInfo.AttriCount; j++ {
			attriMem := make([]byte, ATTRI_INFO_SIZE)
			attri := (*ATTRI_INFO)(BytesToUnsafePointer(attriMem))
			//属性名
			attri.AttriName = GetUtf8FromConstPool(constPool, uint32(BytesToUint16(context[0:2])))

			//属性长度
			attri.Length = BytesToUint32(context[2:6])

			//判断是否是常量属性
			//静态常量的处理
			if attri.AttriName == constSymbol &&
				(filedInfo.AccessFlag&FILED_ACC_STATIC == FILED_ACC_STATIC) &&
				(filedInfo.AccessFlag&FILED_ACC_FINAL == FILED_ACC_FINAL) {
				constValue := uint32(BytesToUint16(context[6:8]))

				if SYM_J == filedInfo.Descriptor ||
					SYM_D == filedInfo.Descriptor {
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
		filedItem := (*FILED_ITEM)(BytesToUnsafePointer(filedMem))
		filedItem.FiledInfoDev = uint32(len(filedInfos))
		filedItem.FiledName = filedName

		//判断是否是long或double
		if filedInfo.AccessFlag&FILED_ACC_STATIC == FILED_ACC_STATIC {
			//静态字段的处理
			filedItem.Index = staticNum
			staticFileds = append(staticFileds, filedMem...)
			staticNum++
			if SYM_J == filedInfo.Descriptor ||
				SYM_D == filedInfo.Descriptor {
				staticNum++
			}
		} else {
			//非静态字段的处理
			filedItem.Index = unstaticNum
			unstaticFileds = append(unstaticFileds, filedMem...)
			unstaticNum++
			if SYM_J == filedInfo.Descriptor ||
				SYM_D == filedInfo.Descriptor {
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
	methodsNum := uint32(BytesToUint16(context[0:2]))
	context = context[2:]
	methods := make([]byte, 0)
	attris := make([]byte, 0)
	var err error
	for i := uint32(0); i < methodsNum; i++ {
		methodInfo_mem := make([]byte, METHOD_SIZE)
		methodInfo := (*METHOD)(BytesToUnsafePointer(methodInfo_mem))
		methodInfo.CodeAddr = INVALID_MEM
		//方法可访问性
		methodInfo.AccessFlag = BytesToUint16(context[0:2])
		//方法名
		methodInfo.MethodName = GetUtf8FromConstPool(constPool, uint32(BytesToUint16(context[2:4])))
		if err != nil {
			return nil, nil, nil, 0, err
		}
		//方法描述符
		methodInfo.Descriptor = GetUtf8FromConstPool(constPool, uint32(BytesToUint16(context[4:6])))

		//属性数量
		attriNum := uint32(BytesToUint16(context[6:8]))
		methodInfo.AttriNum = attriNum
		//属性地址
		methodInfo.Attribute = uint32(len(attris))
		context = context[8:]
		//Code符号表中的地址
		codeSymbol, err := PutSymbol([]byte("Code"))
		if err != nil {
			return nil, nil, nil, 0, err
		}
		for j := uint32(0); j < attriNum; j++ {
			attri_mem := make([]byte, ATTRI_INFO_SIZE)
			attri := (*ATTRI_INFO)(BytesToUnsafePointer(attri_mem))
			//属性名
			attri.AttriName = GetUtf8FromConstPool(constPool, uint32(BytesToUint16(context[0:2])))

			if attri.AttriName == codeSymbol {
				//Code属性的处理
				methodInfo.CodeAddr = uint32(len(attris))
				attriLength := uint32(BytesToUint32(context[2:6]))
				codeMem := readCode(context[uint32(6) : uint32(6)+attriLength])
				attris = append(attris, codeMem...)
				context = context[uint32(6)+attriLength:]
			} else {
				//非Code属性的处理
				attri.Length = uint32(BytesToUint32(context[2:6]))
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
	code := (*CODE_ATTRI)(BytesToUnsafePointer(code_mem))
	code.MaxStack = uint32(BytesToUint16(context[0:2]))
	code.MaxLocal = uint32(BytesToUint16(context[2:4]))
	code.CodeLength = uint32(BytesToUint32(context[4:8]))
	context = context[8:]
	codeOp := make([]byte, 0, code.CodeLength)
	for i := uint32(0); i < code.CodeLength; {
		op := context[i]
		codeOp = append(codeOp, op)
		i++
		switch op {
		case NOP:
		case ACONST_NULL:
		case ICONST_M1:
		case ICONST_0:
		case ICONST_1:
		case ICONST_2:
		case ICONST_3:
		case ICONST_4:
		case ICONST_5:
		case LCONST_0:
		case LCONST_1:
		case FCONST_0:
		case FCONST_1:
		case FCONST_2:
		case DCONST_0:
		case DCONST_1:
		case ILOAD_0:
		case ILOAD_1:
		case ILOAD_2:
		case ILOAD_3:
		case LLOAD_0:
		case LLOAD_1:
		case LLOAD_2:
		case LLOAD_3:
		case FLOAD_0:
		case FLOAD_1:
		case FLOAD_2:
		case FLOAD_3:
		case DLOAD_0:
		case DLOAD_1:
		case DLOAD_2:
		case DLOAD_3:
		case ALOAD_0:
		case ALOAD_1:
		case ALOAD_2:
		case ALOAD_3:
		case IALOAD:
		case LALOAD:
		case FALOAD:
		case DALOAD:
		case AALOAD:
		case BALOAD:
		case CALOAD:
		case SALOAD:
		case ISTORE_0:
		case ISTORE_1:
		case ISTORE_2:
		case ISTORE_3:
		case LSTORE_0:
		case LSTORE_1:
		case LSTORE_2:
		case LSTORE_3:
		case FSTORE_0:
		case FSTORE_1:
		case FSTORE_2:
		case FSTORE_3:
		case DSTORE_0:
		case DSTORE_1:
		case DSTORE_2:
		case DSTORE_3:
		case ASTORE_0:
		case ASTORE_1:
		case ASTORE_2:
		case ASTORE_3:
		case IASTORE:
		case LASTORE:
		case FASTORE:
		case DASTORE:
		case AASTORE:
		case BASTORE:
		case CASTORE:
		case SASTORE:
		case POP:
		case POP2:
		case DUP:
		case DUP_X1:
		case DUP_X2:
		case DUP2:
		case DUP2_X1:
		case DUP2_X2:
		case SWAP:
		case IADD:
		case LADD:
		case FADD:
		case DADD:
		case ISUB:
		case LSUB:
		case FSUB:
		case DSUB:
		case IMUL:
		case LMUL:
		case FMUL:
		case DMUL:
		case IDIV:
		case LDIV:
		case FDIV:
		case DDIV:
		case IREM:
		case LREM:
		case FREM:
		case DREM:
		case INEG:
		case LENG:
		case FNEG:
		case DNEG:
		case ISHL:
		case LSHL:
		case ISHR:
		case LSHR:
		case IUSHR:
		case LUSHR:
		case IAND:
		case LAND:
		case IOR:
		case LOR:
		case IXOR:
		case LXOR:
		case I2L:
		case I2F:
		case I2D:
		case L2I:
		case L2F:
		case L2D:
		case F2I:
		case F2L:
		case F2D:
		case D2I:
		case D2L:
		case D2F:
		case I2B:
		case I2C:
		case I2S:
		case LCMP:
		case FCMPL:
		case FCMPG:
		case DCMPL:
		case DCMPG:
		case IRETURN:
		case LRETURN:
		case FRETURN:
		case DRETURN:
		case ARETURN:
		case RETURN:
		case ARRAYLENGTH:
		case ATHROW:
		case MONITORENTER:
		case MONITOREXIT:
			//无操作数
		case BIPUSH:
			fallthrough
		case LDC:
			fallthrough
		case ILOAD:
			fallthrough
		case LLOAD:
			fallthrough
		case FLOAD:
			fallthrough
		case DLOAD:
			fallthrough
		case ALOAD:
			fallthrough
		case ISTORE:
			fallthrough
		case LSTORE:
			fallthrough
		case FSTORE:
			fallthrough
		case DSTORE:
			fallthrough
		case ASTORE:
			fallthrough
		case NEWARRAY:
			codeOp = append(codeOp, context[i])
			i++

		case SIPUSH:
			b := [2]byte{}
			s := (*int16)(BytesToUnsafePointer(b[:]))
			*s = BytesToInt16(context[i : i+2])
			codeOp = append(codeOp, b[:]...)
			i += 2

		case LDC_W:
			fallthrough
		case LDC2_W:
			fallthrough
		case IFEQ:
			fallthrough
		case IFNE:
			fallthrough
		case IFLT:
			fallthrough
		case IFGE:
			fallthrough
		case IFGT:
			fallthrough
		case IFLE:
			fallthrough
		case IF_ICMPEQ:
			fallthrough
		case IF_ICMPNE:
			fallthrough
		case IF_ICMPLT:
			fallthrough
		case IF_ICMPGE:
			fallthrough
		case IF_ICMPGT:
			fallthrough
		case IF_ICMPLE:
			fallthrough
		case IF_ACMPEQ:
			fallthrough
		case IF_ACMPNE:
			fallthrough
		case GOTO:
			fallthrough
		case GETSTATIC:
			fallthrough
		case PUTSTATIC:
			fallthrough
		case GETFIELD:
			fallthrough
		case PUTFIELD:
			fallthrough
		case INVOKEVIRTUAL:
			fallthrough
		case INVOKESPECIAL:
			fallthrough
		case INVOKESTATIC:
			fallthrough
		case INVOKEINTERFACE:
			fallthrough
		case NEW:
			fallthrough
		case ANEWARRAY:
			fallthrough
		case CHECKCAST:
			fallthrough
		case INSTANCEOF:
			fallthrough
		case IFNULL:
			fallthrough
		case IFNONNULL:
			b := [2]byte{}
			s := (*uint16)(BytesToUnsafePointer(b[:]))
			*s = BytesToUint16(context[i : i+2])
			codeOp = append(codeOp, b[:]...)
			i += 2

		case IINC:
			codeOp = append(codeOp, context[i], context[i+1])
			i += 2

		case TABLESWITCH:
			//3个填充位
			codeOp = append(codeOp, context[i:i+3]...)
			i += 3
			//default
			def := [4]byte{}
			s := (*uint32)(BytesToUnsafePointer(def[:]))
			*s = BytesToUint32(context[i : i+4])
			codeOp = append(codeOp, def[:]...)
			i += 4
			//low
			low_b := [4]byte{}
			low := (*uint32)(BytesToUnsafePointer(low_b[:]))
			*low = BytesToUint32(context[i : i+4])
			codeOp = append(codeOp, low_b[:]...)
			i += 4
			//high
			high_b := [4]byte{}
			high := (*uint32)(BytesToUnsafePointer(high_b[:]))
			*high = BytesToUint32(context[i : i+4])
			codeOp = append(codeOp, high_b[:]...)
			i += 4
			//offset
			for j := *low; j <= *high; j++ {
				offset_b := [4]byte{}
				offset := (*uint32)(BytesToUnsafePointer(offset_b[:]))
				*offset = BytesToUint32(context[i : i+4])
				codeOp = append(codeOp, offset_b[:]...)
				i += 4
			}

		case LOOKUPSWITCH:
			//3个填充位
			codeOp = append(codeOp, context[i:i+3]...)
			i += 3
			//default
			def := [4]byte{}
			s := (*uint32)(BytesToUnsafePointer(def[:]))
			*s = BytesToUint32(context[i : i+4])
			codeOp = append(codeOp, def[:]...)
			i += 4
			//pair的数量
			pairs_b := [4]byte{}
			pairs := (*uint32)(BytesToUnsafePointer(pairs_b[:]))
			*pairs = BytesToUint32(context[i : i+4])
			codeOp = append(codeOp, pairs_b[:]...)
			i += 4
			//各pair
			for j := uint32(0); j <= (*pairs)*2; j++ {
				pair_b := [4]byte{}
				pair := (*uint32)(BytesToUnsafePointer(pair_b[:]))
				*pair = BytesToUint32(context[i : i+4])
				codeOp = append(codeOp, pair_b[:]...)
				i += 4
			}
		case WIDE:
			op := context[i]
			codeOp = append(codeOp, op)
			i++
			switch op {
			case ILOAD:
				fallthrough
			case LLOAD:
				fallthrough
			case FLOAD:
				fallthrough
			case DLOAD:
				fallthrough
			case ALOAD:
				fallthrough
			case ISTORE:
				fallthrough
			case LSTORE:
				fallthrough
			case FSTORE:
				fallthrough
			case DSTORE:
				fallthrough
			case ASTORE:
				b := [2]byte{}
				s := (*uint16)(BytesToUnsafePointer(b[:]))
				*s = BytesToUint16(context[i : i+2])
				codeOp = append(codeOp, b[:]...)
				i += 2
			case IINC:
				b := [2]byte{}
				s := (*uint16)(BytesToUnsafePointer(b[:]))
				*s = BytesToUint16(context[i : i+2])
				codeOp = append(codeOp, b[:]...)
				codeOp = append(codeOp, context[i+2])
				i += 3
			case RET:
				fmt.Println(op)
				panic("readCode():该指令不支持（wide)")
			default:
				fmt.Println(op)
				panic("readCode():该指令不支持wide")
			}

		case MULTIANEWARRAY:
			b := [2]byte{}
			s := (*uint16)(BytesToUnsafePointer(b[:]))
			*s = BytesToUint16(context[i : i+2])
			codeOp = append(codeOp, b[:]...)
			codeOp = append(codeOp, context[i+2])
			i += 3

		case GOTO_W:
			b := [4]byte{}
			s := (*uint32)(BytesToUnsafePointer(b[:]))
			*s = BytesToUint32(context[i : i+4])
			codeOp = append(codeOp, b[:]...)
			i += 4

		case JSR:
			fallthrough
		case RET:
			fallthrough
		case JSR_W:
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
