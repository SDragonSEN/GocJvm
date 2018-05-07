package method

import (
	"fmt"
	"unicode/utf16"

	. "access/access"
	. "access/array"
	. "access/string"
	. "basic/com"
	. "basic/memCtrl"
	. "basic/symbol"
	. "class/classInterface"
	. "class/classParse"
	. "class/classTable"
)

type METHOD_STACK struct {
	MaxStackSize uint32
	StackNum     uint32
	TopFrame     uint32 //地址
	PC           uint32
}

const METHOD_STACK_SIZE = 4 * 4

type METHOD_FRAME struct {
	LowFrame        uint32 //地址
	Claz            uint32 //类地址
	ReturnPc        uint32 //返回后的PC地址
	VarSize         uint32 //变量区大小
	OpStackSize     uint32 //操作数栈大小
	CurOpStackIndex uint32 //当前操作数栈索引
	LocalAdr        uint32 //Frame的地址
}

const METHOD_FRAME_SIZE = 7 * 4

/******************************************************************
    功能:新建方法栈
	入参:无
    返回值:1、*METHOD_STACK
	注:函数栈大小默认为100
******************************************************************/
func NewMethodStack() *METHOD_STACK {
	adr, err := Malloc(METHOD_STACK_SIZE, METHOD_STACK_NODE)
	if err != nil {
		panic("NewMethodStack()")
	}
	methodStack := (*METHOD_STACK)(GetPointer(adr, METHOD_STACK_SIZE))
	methodStack.MaxStackSize = 100 //目前默认100
	methodStack.StackNum = 0
	methodStack.TopFrame = INVALID_MEM
	return methodStack
}

/******************************************************************
    功能:方法压栈
	入参:无
    返回值:1、*METHOD_FRAME
	      2、地址
******************************************************************/
func (self *METHOD_STACK) PushFrame(varSize, opStackSize, clazAdr, returnPc uint32) *METHOD_FRAME {
	adr, err := Malloc(METHOD_FRAME_SIZE+(varSize+opStackSize)*4, METHOD_FRAME_NODE)
	if err != nil {
		panic("PushFrame()")
	}
	methodFrame := (*METHOD_FRAME)(GetPointer(adr, METHOD_FRAME_SIZE))
	methodFrame.VarSize = varSize
	methodFrame.OpStackSize = opStackSize
	methodFrame.Claz = clazAdr
	methodFrame.CurOpStackIndex = 0
	methodFrame.ReturnPc = returnPc
	methodFrame.LowFrame = self.TopFrame
	methodFrame.LocalAdr = adr
	self.StackNum++
	self.TopFrame = adr
	if self.StackNum > self.MaxStackSize {
		panic("stack overflow！")
	}
	return methodFrame
}

/******************************************************************
    功能:方法压栈
	入参:无
    返回值:1、*METHOD_FRAME
	      2、地址
******************************************************************/
func (self *METHOD_STACK) PopFrame() *METHOD_FRAME {
	if self.StackNum == 0 {
		panic("PopFrame(): stack is empty!")
	}
	curFrameAdr := self.TopFrame
	curFrame := (*METHOD_FRAME)(GetPointer(curFrameAdr, METHOD_FRAME_SIZE))

	self.TopFrame = curFrame.LowFrame
	self.PC = curFrame.ReturnPc
	self.StackNum--
	MemFree(curFrameAdr)
	if self.StackNum == 0 {
		return nil
	}
	return (*METHOD_FRAME)(GetPointer(self.TopFrame, METHOD_FRAME_SIZE))
}

/******************************************************************
    功能:方法帧设置变量区值
	入参:1、index
	    2、value
    返回值:无
******************************************************************/
func (self *METHOD_FRAME) SetVar(index, value uint32) {
	if index >= self.VarSize {
		panic("SetVar()方法变量异常")
	}
	p := (*uint32)(GetPointer(self.LocalAdr+METHOD_FRAME_SIZE+(index*4), 4))
	*p = value
}

/******************************************************************
    功能:获取方法帧变量区值
	入参:index
    返回值:value
******************************************************************/
func (self *METHOD_FRAME) GetVar(index uint32) uint32 {
	if index >= self.VarSize {
		panic("GetVar()方法变量异常")
	}
	p := (*uint32)(GetPointer(self.LocalAdr+METHOD_FRAME_SIZE+(index*4), 4))
	return *p
}

/******************************************************************
    功能:操作数栈压栈
	入参:value
    返回值:无
******************************************************************/
func (self *METHOD_FRAME) Push(value uint32) {
	if self.CurOpStackIndex >= self.OpStackSize {
		panic("Push()操作数栈异常")
	}
	p := (*uint32)(GetPointer(self.LocalAdr+METHOD_FRAME_SIZE+(self.VarSize+self.CurOpStackIndex)*4, 4))
	*p = value
	self.CurOpStackIndex++
}

/******************************************************************
    功能:操作数栈弹栈
	入参:无
    返回值:value
******************************************************************/
func (self *METHOD_FRAME) Pop() uint32 {
	if self.CurOpStackIndex == 0 {
		panic("Push()操作数栈异常")
	}
	p := (*uint32)(GetPointer(self.LocalAdr+METHOD_FRAME_SIZE+(self.VarSize+self.CurOpStackIndex-1)*4, 4))
	self.CurOpStackIndex--
	return *p
}

/******************************************************************
    功能:操作数栈弹栈
	入参:无
    返回值:value
******************************************************************/
func (self *METHOD_FRAME) PopInt64() int64 {
	v1 := self.Pop()
	v0 := self.Pop()

	v := int64((uint64(v0) << 32) | uint64(v1))
	return v
}

/******************************************************************
    功能:获取栈顶元素
	入参:无
    返回值:value
******************************************************************/
func (self *METHOD_FRAME) Peek() uint32 {
	if self.CurOpStackIndex == 0 {
		panic("Push()操作数栈异常")
	}
	p := (*uint32)(GetPointer(self.LocalAdr+METHOD_FRAME_SIZE+(self.VarSize+self.CurOpStackIndex-1)*4, 4))
	return *p
}

/******************************************************************
    功能:执行函数
	入参:无
    返回值:无
******************************************************************/
func (self *METHOD_STACK) Excute() {
	i := 0
	frame := (*METHOD_FRAME)(GetPointer(self.TopFrame, METHOD_FRAME_SIZE))
	for {
		self.Log(frame)
		fmt.Println(Format(Memory[self.PC]))
		fmt.Println(self.PC)
		switch Memory[self.PC] {
		case NOP:
			if i > 20 {
				panic("NOP")
			}
			i++
			self.PC++
		case ACONST_NULL:
			self.Aconst_Null(frame)
		case ICONST_M1:
			self.IConst(frame, -1)
		case ICONST_0:
			self.IConst(frame, 0)
		case ICONST_1:
			self.IConst(frame, 1)
		case ICONST_2:
			self.IConst(frame, 2)
		case ICONST_3:
			self.IConst(frame, 3)
		case ICONST_4:
			self.IConst(frame, 4)
		case ICONST_5:
			self.IConst(frame, 5)
		case BIPUSH:
			self.BIPush(frame)
		case LDC:
			self.Ldc(frame)
		case LDC_W:
			self.Ldc_w(frame)
		case LDC2_W:
			self.Ldc2_w(frame)
		case ILOAD:
			self.ILoad(frame)
		case ALOAD:
			self.ILoad(frame)
		case ILOAD_0:
			self.Load32(frame, 0)
		case ILOAD_1:
			self.Load32(frame, 1)
		case ILOAD_2:
			self.Load32(frame, 2)
		case ILOAD_3:
			self.Load32(frame, 3)
		case LLOAD_2:
			self.Load64(frame, 2)
		case FLOAD_0:
			self.Load32(frame, 0)
		case DLOAD_0:
			self.Load64(frame, 0)
		case ALOAD_0:
			self.Load32(frame, 0)
		case ALOAD_1:
			self.Load32(frame, 1)
		case ALOAD_2:
			self.Load32(frame, 2)
		case IALOAD:
			self.IALoad(frame)
		case BALOAD:
			self.BALoad(frame)
		case ASTORE:
			self.Store32Byte(frame)
		case ISTORE:
			self.IStore(frame)
		case ISTORE_1:
			self.Store32(frame, 1)
		case ISTORE_2:
			self.Store32(frame, 2)
		case ISTORE_3:
			self.Store32(frame, 3)
		case LSTORE_2:
			self.Store64(frame, 2)
		case ASTORE_1:
			self.Store32(frame, 1)
		case ASTORE_2:
			self.Store32(frame, 2)
		case IASTORE:
			self.IAStore(frame)
		case BASTORE:
			self.BAStore(frame)
		case POP:
			frame.Pop()
			self.PC++
		case DUP:
			self.Dup(frame)
		case IADD:
			self.Iadd(frame)
		case ISUB:
			self.Isub(frame)
		case IMUL:
			self.Imul(frame)
		case IDIV:
			self.Idiv(frame)
		case IREM:
			self.Irem(frame)
		case INEG:
			self.Ineg(frame)
		case IAND:
			self.Iand(frame)
		case LAND:
			self.Land(frame)
		case I2L:
			self.I2L(frame)
		case LCMP:
			self.Lcmp(frame)
		case IFGE:
			self.Ifge(frame)
		case IFNE:
			self.Ifne(frame)
		case IF_ICMPNE:
			self.Icmpne(frame)
		case IF_ICMPLT:
			self.Icmplt(frame)
		case IF_ICMPGT:
			self.Icmpgt(frame)
		case IF_ICMPLE:
			self.Icmple(frame)

		case IRETURN:
			frame = self.IReturn(frame)
			if frame == nil {
				goto label
			}
		case LRETURN:
			frame = self.LReturn(frame)
			if frame == nil {
				goto label
			}
		case GOTO:
			self.Goto(frame)
		case GETSTATIC:
			self.GetStatic(frame)
		case PUTSTATIC:
			self.PutStatic(frame)
		case GETFIELD:
			self.GetFiled(frame)
		case PUTFIELD:
			self.PutFiled(frame)
		case INVOKEVIRTUAL:
			newFrame := self.InvokeVirtual(frame)
			if newFrame != nil {
				frame = newFrame
			}
		case INVOKESTATIC:
			frame = self.InvokeStatic(frame)
		case INVOKESPECIAL:
			frame = self.InvokeSpecial(frame)
		case RETURN:
			frame = self.PopFrame()
			if frame == nil {
				goto label
			}
		case NEW:
			self.New(frame)
		case ARRAYLENGTH:
			self.ArrayLength(frame)
		case NEWARRAY:
			self.NewArray(frame)
		case ANEWARRAY:
			self.ANewArray(frame)
		case IFNULL:
			self.IfNull(frame)
		default:
			fmt.Printf("Memory[self.PC]:%x ", Memory[self.PC])
			panic(Format(Memory[self.PC]))
			goto label
		}
	}
label:
}

/******************************************************************
    功能:getstatic指令
	入参:无
    返回值:无
******************************************************************/
func (self *METHOD_STACK) GetStatic(frame *METHOD_FRAME) {
	self.PC++
	p := (*uint16)(GetPointer(self.PC, 2))
	filedInfo := GetFiledInfo(GetConstantPoolSlice(frame.Claz), uint32(*p))
	//fmt.Println(string(GetSymbol(filedInfo.ClassName)), string(GetSymbol(filedInfo.FiledName)), string(GetSymbol(filedInfo.FiledType)))
	self.PC += 2

	var classAdr uint32
	classAdr = GetClassMemAddr(filedInfo.ClassName)
	//如果获取不到，则说明不在内存中，需要去加载
	if classAdr == INVALID_MEM {
		//获取类名(string)
		className := string(GetSymbol(filedInfo.ClassName))
		classInfo, err := LoadClass(className)
		if err != nil {
			panic("GetStatic()")
		}
		classAdr = classInfo.LocalAdr
		CInit(classInfo.LocalAdr)
	}

	classInfo := (*CLASS_INFO)(GetPointer(classAdr, CLASS_INFO_SIZE))
	//判断是否是long或double型
	if filedInfo.FiledType == SYM_J ||
		filedInfo.FiledType == SYM_D {
		v := classInfo.GetStaticData64(filedInfo.FiledName, filedInfo.FiledType)
		frame.Push(v[0])
		frame.Push(v[1])
	} else {
		v := classInfo.GetStaticData32(filedInfo.FiledName, filedInfo.FiledType)
		frame.Push(v)
	}
}

/******************************************************************
    功能:putstatic指令
	入参:无
    返回值:无
******************************************************************/
func (self *METHOD_STACK) PutStatic(frame *METHOD_FRAME) {
	self.PC++
	p := (*uint16)(GetPointer(self.PC, 2))
	filedInfo := GetFiledInfo(GetConstantPoolSlice(frame.Claz), uint32(*p))
	//fmt.Println(string(GetSymbol(filedInfo.ClassName)), string(GetSymbol(filedInfo.FiledName)), string(GetSymbol(filedInfo.FiledType)))
	self.PC += 2

	var classAdr uint32
	classAdr = GetClassMemAddr(filedInfo.ClassName)
	//如果获取不到，则说明不在内存中，需要去加载
	if classAdr == INVALID_MEM {
		//获取类名(string)
		className := string(GetSymbol(filedInfo.ClassName))
		classInfo, err := LoadClass(className)
		if err != nil {
			panic("PutStatic()")
		}
		classAdr = classInfo.LocalAdr
		CInit(classInfo.LocalAdr)
	}

	classInfo := (*CLASS_INFO)(GetPointer(classAdr, CLASS_INFO_SIZE))
	//判断是否是long或double型
	if filedInfo.FiledType == SYM_J ||
		filedInfo.FiledType == SYM_D {
		v1 := frame.Pop()
		v0 := frame.Pop()
		classInfo.PutStaticData64(filedInfo.FiledName, filedInfo.FiledType, v0, v1)
	} else {
		v := frame.Pop()
		classInfo.PutStaticData32(filedInfo.FiledName, filedInfo.FiledType, v)
	}
}

/******************************************************************
    功能:GetFiled
	入参:1、*METHOD_FRAME
    返回值:无
******************************************************************/
func (self *METHOD_STACK) GetFiled(frame *METHOD_FRAME) {
	self.PC++
	p := (*uint16)(GetPointer(self.PC, 2))
	filedInfo := GetFiledInfo(GetConstantPoolSlice(frame.Claz), uint32(*p))
	//	fmt.Println(string(GetSymbol(filedInfo.FiledName)))
	self.PC += 2
	accessAdr := frame.Pop()
	this := (*ACCESS_INFO)(GetPointer(accessAdr, ACCESS_INFO_SIZE))
	thisClass := (*CLASS_INFO)(GetPointer(this.TypeAddr, CLASS_INFO_SIZE))
	index := thisClass.GetUnstaticDataIndex(filedInfo.FiledName, filedInfo.FiledType)
	data := GetData(accessAdr)
	v0 := (*uint32)(BytesToUnsafePointer(data[index*4 : index*4+4]))
	frame.Push(*v0)
	if filedInfo.FiledType == SYM_J ||
		filedInfo.FiledType == SYM_D {
		v1 := (*uint32)(BytesToUnsafePointer(data[index*4+4 : index*4+8]))
		frame.Push(*v1)
	}
}

/******************************************************************
    功能:PutFiled
	入参:1、*METHOD_FRAME
    返回值:无
******************************************************************/
func (self *METHOD_STACK) PutFiled(frame *METHOD_FRAME) {
	self.PC++
	p := (*uint16)(GetPointer(self.PC, 2))
	filedInfo := GetFiledInfo(GetConstantPoolSlice(frame.Claz), uint32(*p))
	self.PC += 2
	if filedInfo.FiledType == SYM_J ||
		filedInfo.FiledType == SYM_D {
		v1 := frame.Pop()
		v0 := frame.Pop()
		accessAdr := frame.Pop()
		this := (*ACCESS_INFO)(GetPointer(accessAdr, ACCESS_INFO_SIZE))
		thisClass := (*CLASS_INFO)(GetPointer(this.TypeAddr, CLASS_INFO_SIZE))
		index := thisClass.GetUnstaticDataIndex(filedInfo.FiledName, filedInfo.FiledType)
		data := GetData(accessAdr)
		p0 := (*uint32)(BytesToUnsafePointer(data[index*4 : index*4+4]))
		*p0 = v0
		p1 := (*uint32)(BytesToUnsafePointer(data[index*4+4 : index*4+8]))
		*p1 = v1
	} else {
		v := frame.Pop()
		accessAdr := frame.Pop()
		this := (*ACCESS_INFO)(GetPointer(accessAdr, ACCESS_INFO_SIZE))
		thisClass := (*CLASS_INFO)(GetPointer(this.TypeAddr, CLASS_INFO_SIZE))
		index := thisClass.GetUnstaticDataIndex(filedInfo.FiledName, filedInfo.FiledType)
		data := GetData(accessAdr)
		p := (*uint32)(BytesToUnsafePointer(data[index*4 : index*4+4]))
		*p = v
	}
}

/******************************************************************
    功能:BIPush指令
	入参:无
    返回值:无
******************************************************************/
func (self *METHOD_STACK) BIPush(frame *METHOD_FRAME) {
	self.PC++
	v := (*int8)(GetPointer(self.PC, 1))
	self.PC++
	frame.Push(uint32(*v))
}

/******************************************************************
    功能:ldc指令
	入参:无
    返回值:无
******************************************************************/
func (self *METHOD_STACK) Ldc(frame *METHOD_FRAME) {
	self.PC++

	classIndex := uint32(Memory[self.PC])
	fmt.Println(classIndex)
	if IsClassConstant(frame.Claz, classIndex) {
		fmt.Println("Ldc0")
		className := GetClassFromConstPool(GetConstantPoolSlice(frame.Claz), classIndex)
		classInstant := GetClassMemAddr(className)
		if classInstant == INVALID_MEM {
			classInfo, err := LoadClass(string(GetSymbol(className)))
			if err != nil {
				panic("Ldc()")
			}
			CInit(classInfo.LocalAdr)
			classInstant = classInfo.LocalAdr
		}
		frame.Push(classInstant)
		self.PC++
		return
	}
	fmt.Println("Ldc1")
	v := GetUint32FromConstPool(GetConstantPoolSlice(frame.Claz), classIndex)
	frame.Push(v)
	self.PC++
}

/******************************************************************
    功能:Ldc_w指令
	入参:无
    返回值:无
******************************************************************/
func (self *METHOD_STACK) Ldc_w(frame *METHOD_FRAME) {
	self.PC++
	p := uint32(*(*uint16)(GetPointer(self.PC, 2)))
	if IsClassConstant(frame.Claz, p) {
		className := GetClassFromConstPool(GetConstantPoolSlice(frame.Claz), p)
		classInstant := GetClassMemAddr(className)
		if classInstant == INVALID_MEM {
			classInfo, err := LoadClass(string(GetSymbol(className)))
			if err != nil {
				panic("Ldc()")
			}
			CInit(classInfo.LocalAdr)
			classInstant = classInfo.LocalAdr
		}
		frame.Push(classInstant)
		self.PC += 2
		return
	}
	v := GetUint32FromConstPool(GetConstantPoolSlice(frame.Claz), p)
	frame.Push(v)
	self.PC += 2
}

/******************************************************************
    功能:Ldc2_w指令
	入参:无
    返回值:无
******************************************************************/
func (self *METHOD_STACK) Ldc2_w(frame *METHOD_FRAME) {
	self.PC++
	p := (*uint16)(GetPointer(self.PC, 2))
	v0, v1 := GetUint64FromConstPool(GetConstantPoolSlice(frame.Claz), uint32(*p))
	frame.Push(v0)
	frame.Push(v1)
	self.PC += 2
}

/******************************************************************
    功能:InvokeVirtual指令
	入参:无
    返回值:无
******************************************************************/
func (self *METHOD_STACK) InvokeVirtual(frame *METHOD_FRAME) *METHOD_FRAME {
	self.PC++
	p := (*uint16)(GetPointer(self.PC, 2))
	self.PC += 2
	methodRef := GetMethodInfo(GetConstantPoolSlice(frame.Claz), uint32(*p))
	if StubInvokeFunc(frame, methodRef) {
		return nil
	}
	fmt.Println(string(GetSymbol(methodRef.ClassName)),
		string(GetSymbol(methodRef.MethodName)),
		string(GetSymbol(methodRef.MethodDesp)))
	//反射函数调用
	if methodRef.ClassName == SYM_JAVA_LANG_CLASS {
		return self.ReflectMethod(frame, &methodRef)
	}
	num := CalParaSize(string(GetSymbol(methodRef.MethodDesp)))
	param := make([]uint32, num)
	//将上一个栈帧中的值弹出，保存到新的栈帧中的局部变量中
	for i := num; i > 0; i-- {
		param[i-1] = frame.Pop()
	}

	//获取this中的类
	this := frame.Pop()
	classInfo := (*CLASS_INFO)(GetPointer(GetClassInfo(this), CLASS_INFO_SIZE))

	if StubInvokeFunc(frame, methodRef) {
		return nil
	}

	//查找方法
	methodInfo, methodCLass, codeAdr := classInfo.FindMethodEx(methodRef.MethodName, methodRef.MethodDesp)
	if methodInfo == nil || codeAdr == INVALID_MEM {
		fmt.Println(methodInfo, codeAdr)
		fmt.Println(string(GetSymbol(methodRef.ClassName)),
			string(GetSymbol(methodRef.MethodName)),
			string(GetSymbol(methodRef.MethodDesp)))
		panic("InvokeVirtual()6")
	}

	codeAttri := (*CODE_ATTRI)(GetPointer(codeAdr, CODE_ATTRI_SIZE))
	//创建方法栈
	newFrame := self.PushFrame(codeAttri.MaxLocal, codeAttri.MaxStack, methodCLass.LocalAdr, self.PC)
	self.PC = codeAdr + CODE_ATTRI_SIZE
	newFrame.SetVar(0, this)
	for i, k := range param {
		newFrame.SetVar(uint32(i+1), k)
	}

	return newFrame
}

/******************************************************************
    功能:InvokeStatic指令
	入参:无
    返回值:无
******************************************************************/
func (self *METHOD_STACK) InvokeStatic(frame *METHOD_FRAME) *METHOD_FRAME {
	self.PC++
	p := (*uint16)(GetPointer(self.PC, 2))
	self.PC += 2
	//获取方法描述
	methodRef := GetMethodInfo(GetConstantPoolSlice(frame.Claz), uint32(*p))

	//获取方法中的类
	var classInfo *CLASS_INFO
	var err error
	classAdr := GetClassMemAddr(methodRef.ClassName)
	if classAdr == INVALID_MEM {
		classInfo, err = LoadClass(string(GetSymbol(methodRef.ClassName)))
		if err != nil {
			panic("InvokeStatic()")
		}
		CInit(classInfo.LocalAdr)
	} else {
		classInfo = (*CLASS_INFO)(GetPointer(classAdr, CLASS_INFO_SIZE))
	}
	//查找方法
	methodInfo, codeAdr := classInfo.FindMethod(methodRef.MethodName, methodRef.MethodDesp)
	if methodInfo.AccessFlag&METHOD_ACC_NATIVE == METHOD_ACC_NATIVE {
		ExcuteLocalMethodAdp(methodRef.ClassName, methodRef.MethodName, methodRef.MethodDesp, frame)

		return frame
	} else if codeAdr == INVALID_MEM {
		panic("方法没找到!")
	}
	fmt.Println(string(GetSymbol(methodRef.ClassName)),
		string(GetSymbol(methodRef.MethodName)),
		string(GetSymbol(methodRef.MethodDesp)))

	codeAttri := (*CODE_ATTRI)(GetPointer(codeAdr, CODE_ATTRI_SIZE))
	//创建方法栈
	newFrame := self.PushFrame(codeAttri.MaxLocal, codeAttri.MaxStack, classInfo.LocalAdr, self.PC)
	self.PC = codeAdr + CODE_ATTRI_SIZE

	//计算需要弹出的参数个数
	num := CalParaSize(string(GetSymbol(methodRef.MethodDesp)))

	//将上一个栈帧中的值弹出，保存到新的栈帧中的局部变量中
	for i := num; i > 0; i-- {
		newFrame.SetVar(uint32(i-1), frame.Pop())
	}

	return newFrame
}

/******************************************************************
    功能:InvokeSpecial指令
	入参:无
    返回值:无
******************************************************************/
func (self *METHOD_STACK) InvokeSpecial(frame *METHOD_FRAME) *METHOD_FRAME {
	self.PC++
	p := (*uint16)(GetPointer(self.PC, 2))
	self.PC += 2
	//获取方法描述
	methodRef := GetMethodInfo(GetConstantPoolSlice(frame.Claz), uint32(*p))

	//获取方法中的类
	var classInfo *CLASS_INFO
	var err error
	classAdr := GetClassMemAddr(methodRef.ClassName)
	if classAdr == INVALID_MEM {
		classInfo, err = LoadClass(string(GetSymbol(methodRef.ClassName)))
		if err != nil {
			panic("InvokeSpecial()")
		}
		CInit(classInfo.LocalAdr)
	} else {
		classInfo = (*CLASS_INFO)(GetPointer(classAdr, CLASS_INFO_SIZE))
	}
	fmt.Println(string(GetSymbol(classInfo.ClassName)),
		string(GetSymbol(methodRef.MethodName)),
		string(GetSymbol(methodRef.MethodDesp)))
	//查找方法
	methodInfo, methodCLass, codeAdr := classInfo.FindMethodEx(methodRef.MethodName, methodRef.MethodDesp)
	if methodInfo == nil || codeAdr == INVALID_MEM {
		fmt.Println(methodInfo, codeAdr)
		panic("InvokeSpecial()6")
	}
	codeAttri := (*CODE_ATTRI)(GetPointer(codeAdr, CODE_ATTRI_SIZE))
	//创建方法栈
	newFrame := self.PushFrame(codeAttri.MaxLocal, codeAttri.MaxStack, methodCLass.LocalAdr, self.PC)
	self.PC = codeAdr + CODE_ATTRI_SIZE

	//计算需要弹出的参数个数
	num := CalParaSize(string(GetSymbol(methodRef.MethodDesp)))
	//加上this
	num++
	//将上一个栈帧中的值弹出，保存到新的栈帧中的局部变量中
	for i := num; i > 0; i-- {
		newFrame.SetVar(uint32(i-1), frame.Pop())
	}

	return newFrame
}

/******************************************************************
    功能:New
	入参:*METHOD_FRAME
    返回值:无
******************************************************************/
func (self *METHOD_STACK) New(frame *METHOD_FRAME) {
	self.PC++
	p := (*uint16)(GetPointer(self.PC, 2))
	self.PC += 2
	className := GetClassFromConstPool(GetConstantPoolSlice(frame.Claz), uint32(*p))
	var classAdr uint32
	var classInfo *CLASS_INFO
	classAdr = GetClassMemAddr(className)
	//如果获取不到，则说明不在内存中，需要去加载
	if classAdr == INVALID_MEM {
		//获取类名(string)
		classNameStr := string(GetSymbol(className))
		classInfo, err := LoadClass(classNameStr)
		if err != nil {
			panic("GetStatic()")
		}
		classAdr = classInfo.LocalAdr
		CInit(classInfo.LocalAdr)
	}
	accessInfo, accessAdr, err := NewAccessInfo()
	if err != nil {
		panic(err)
	}
	classInfo = (*CLASS_INFO)(GetPointer(classAdr, CLASS_INFO_SIZE))
	accessInfo.TypeAddr = classAdr
	accessInfo.DataAddr, err = Malloc(classInfo.UnstaticParaTotalSize, INSTANCE_NODE)
	if err != nil {
		panic(err)
	}
	frame.Push(accessAdr)
}

/******************************************************************
    功能:ArrayLength
	入参:*METHOD_FRAME
    返回值:无
******************************************************************/
func (self *METHOD_STACK) ArrayLength(frame *METHOD_FRAME) {
	self.PC++
	arrAdr := frame.Pop()
	if arrAdr == INVALID_MEM {
		panic("null pointer")
	}
	arrInfo, _ := GetArrayInfo(arrAdr)
	frame.Push(arrInfo.Length)
}

/******************************************************************
    功能:NewArray
	入参:*METHOD_FRAME
    返回值:无
******************************************************************/
func (self *METHOD_STACK) NewArray(frame *METHOD_FRAME) {
	self.PC++
	p := (*uint8)(GetPointer(self.PC, 1))
	var symbol uint32
	var width uint32
	switch *p {
	case AT_BOOLEAN:
		symbol = SYM_KZ
		width = 1
	case AT_BYTE:
		symbol = SYM_KB
		width = 1
	case AT_CHAR:
		symbol = SYM_KC
		width = 2
	case AT_FLOAT:
		symbol = SYM_KF
		width = 4
	case AT_DOUBLE:
		symbol = SYM_KD
		width = 8
	case AT_SHORT:
		symbol = SYM_KS
		width = 2
	case AT_INT:
		symbol = SYM_KI
		width = 4
	case AT_LONG:
		symbol = SYM_KJ
		width = 8
	}
	self.PC++

	_, arrAdr, err := NewArray(symbol, width, frame.Pop())
	if err != nil {
		panic("NewArray()8")
	}
	frame.Push(arrAdr)
}

/******************************************************************
    功能:ANewArray
	入参:*METHOD_FRAME
    返回值:无
******************************************************************/
func (self *METHOD_STACK) ANewArray(frame *METHOD_FRAME) {
	self.PC++
	p := (*uint8)(GetPointer(self.PC, 2))
	self.PC += 2

	_, arrAdr, err := NewArray(GetUtf8FromConstPool(GetConstantPoolSlice(frame.Claz), uint32(*p)), 4, frame.Pop())
	if err != nil {
		panic("NewArray()8")
	}
	frame.Push(arrAdr)
}

/******************************************************************
    功能:Dup
	入参:*METHOD_FRAME
    返回值:无
******************************************************************/
func (self *METHOD_STACK) Dup(frame *METHOD_FRAME) {
	self.PC++
	frame.Push(frame.Peek())
}

/******************************************************************
    功能:Iadd
	入参:*METHOD_FRAME
    返回值:无
******************************************************************/
func (self *METHOD_STACK) Iadd(frame *METHOD_FRAME) {
	self.PC++
	v0 := int32(frame.Pop())
	v1 := int32(frame.Pop())

	frame.Push(uint32(v0 + v1))
}

/******************************************************************
    功能:Isub
	入参:*METHOD_FRAME
    返回值:无
******************************************************************/
func (self *METHOD_STACK) Isub(frame *METHOD_FRAME) {
	self.PC++
	v1 := int32(frame.Pop())
	v0 := int32(frame.Pop())

	frame.Push(uint32(v0 - v1))
}

/******************************************************************
    功能:Imul
	入参:*METHOD_FRAME
    返回值:无
******************************************************************/
func (self *METHOD_STACK) Imul(frame *METHOD_FRAME) {
	self.PC++
	v1 := int32(frame.Pop())
	v0 := int32(frame.Pop())

	frame.Push(uint32(v0 * v1))
}

/******************************************************************
    功能:Idiv
	入参:*METHOD_FRAME
    返回值:无
******************************************************************/
func (self *METHOD_STACK) Idiv(frame *METHOD_FRAME) {
	self.PC++
	v1 := int32(frame.Pop())
	v0 := int32(frame.Pop())
	if v1 == 0 {
		panic("Idiv")
	}
	frame.Push(uint32(v0 / v1))
}

/******************************************************************
    功能:Irem
	入参:*METHOD_FRAME
    返回值:无
******************************************************************/
func (self *METHOD_STACK) Irem(frame *METHOD_FRAME) {
	self.PC++
	v1 := int32(frame.Pop())
	v0 := int32(frame.Pop())
	if v1 == 0 {
		panic("Irem")
	}
	frame.Push(uint32(v0 % v1))
}

/******************************************************************
    功能:Iand
	入参:*METHOD_FRAME
    返回值:无
******************************************************************/
func (self *METHOD_STACK) Iand(frame *METHOD_FRAME) {
	self.PC++
	v1 := uint32(frame.Pop())
	v0 := uint32(frame.Pop())

	frame.Push(v0 & v1)
}

/******************************************************************
    功能:Iand
	入参:*METHOD_FRAME
    返回值:无
******************************************************************/
func (self *METHOD_STACK) Land(frame *METHOD_FRAME) {
	self.PC++
	v1 := uint32(frame.Pop())
	v0 := uint32(frame.Pop())

	w1 := uint32(frame.Pop())
	w0 := uint32(frame.Pop())

	frame.Push(v0 & w0)
	frame.Push(v1 & w1)
}

/******************************************************************
    功能:Iand
	入参:*METHOD_FRAME
    返回值:无
******************************************************************/
func (self *METHOD_STACK) I2L(frame *METHOD_FRAME) {
	self.PC++
	v1 := int32(frame.Pop())
	n := uint64(v1)
	frame.Push(uint32(n >> 32))
	frame.Push(uint32(n & 0xffffffff))
}

/******************************************************************
    功能:Ineg
	入参:*METHOD_FRAME
    返回值:无
******************************************************************/
func (self *METHOD_STACK) Ineg(frame *METHOD_FRAME) {
	self.PC++
	v0 := int32(frame.Pop())
	frame.Push(uint32(-v0))
}

/******************************************************************
    功能:Lcmp
	入参:*METHOD_FRAME
    返回值:无
******************************************************************/
func (self *METHOD_STACK) Lcmp(frame *METHOD_FRAME) {
	self.PC++
	v0 := frame.PopInt64()
	v1 := frame.PopInt64()
	v := int32(-1)
	if v0 > v1 {
		frame.Push(uint32(1))
	} else if v0 < v1 {
		frame.Push(uint32(v))
	} else {
		frame.Push(uint32(0))
	}
}

/******************************************************************
    功能:Ifge
	入参:*METHOD_FRAME
    返回值:无
******************************************************************/
func (self *METHOD_STACK) Ifge(frame *METHOD_FRAME) {
	p := (*int16)(GetPointer(self.PC+1, 2))
	v := int32(frame.Pop())
	if v >= 0 {
		self.PC = uint32(int32(self.PC) + int32(*p))
	} else {
		self.PC += 3
	}
}

/******************************************************************
    功能:Ifne
	入参:*METHOD_FRAME
    返回值:无
******************************************************************/
func (self *METHOD_STACK) Ifne(frame *METHOD_FRAME) {
	p := (*int16)(GetPointer(self.PC+1, 2))
	v := int32(frame.Pop())
	if v != 0 {
		self.PC = uint32(int32(self.PC) + int32(*p))
	} else {
		self.PC += 3
	}
}

/******************************************************************
    功能:Icmpne
	入参:*METHOD_FRAME
    返回值:无
******************************************************************/
func (self *METHOD_STACK) Icmpne(frame *METHOD_FRAME) {
	p := (*int16)(GetPointer(self.PC+1, 2))
	v0 := int32(frame.Pop())
	v1 := int32(frame.Pop())
	if v0 != v1 {
		self.PC = uint32(int32(self.PC) + int32(*p))
	} else {
		self.PC += 3
	}
}

/******************************************************************
    功能:Icmplt
	入参:*METHOD_FRAME
    返回值:无
******************************************************************/
func (self *METHOD_STACK) Icmplt(frame *METHOD_FRAME) {
	p := (*int16)(GetPointer(self.PC+1, 2))
	v0 := int32(frame.Pop())
	v1 := int32(frame.Pop())
	if v0 > v1 {
		self.PC = uint32(int32(self.PC) + int32(*p))
	} else {
		self.PC += 3
	}
}

/******************************************************************
    功能:Icmpgt
	入参:*METHOD_FRAME
    返回值:无
******************************************************************/
func (self *METHOD_STACK) Icmpgt(frame *METHOD_FRAME) {
	p := (*int16)(GetPointer(self.PC+1, 2))
	v0 := int32(frame.Pop())
	v1 := int32(frame.Pop())
	if v0 < v1 {
		self.PC = uint32(int32(self.PC) + int32(*p))
	} else {
		self.PC += 3
	}
}

/******************************************************************
    功能:Icmple
	入参:*METHOD_FRAME
    返回值:无
******************************************************************/
func (self *METHOD_STACK) Icmple(frame *METHOD_FRAME) {
	p := (*int16)(GetPointer(self.PC+1, 2))
	v0 := int32(frame.Pop())
	v1 := int32(frame.Pop())
	if v0 >= v1 {
		self.PC = uint32(int32(self.PC) + int32(*p))
	} else {
		self.PC += 3
	}
}

/******************************************************************
    功能:IReturn
	入参:*METHOD_FRAME
    返回值:无
******************************************************************/
func (self *METHOD_STACK) IReturn(frame *METHOD_FRAME) *METHOD_FRAME {
	v := frame.Pop()
	newFrame := self.PopFrame()
	if newFrame == nil {
		return nil
	}
	newFrame.Push(v)
	return newFrame
}

/******************************************************************
    功能:LReturn
	入参:*METHOD_FRAME
    返回值:无
******************************************************************/
func (self *METHOD_STACK) LReturn(frame *METHOD_FRAME) *METHOD_FRAME {
	v1 := frame.Pop()
	v0 := frame.Pop()
	newFrame := self.PopFrame()
	if newFrame == nil {
		return nil
	}
	newFrame.Push(v0)
	newFrame.Push(v1)
	return newFrame
}

/******************************************************************
    功能:Icmpgt
	入参:*METHOD_FRAME
    返回值:无
******************************************************************/
func (self *METHOD_STACK) Goto(frame *METHOD_FRAME) {
	p := (*int16)(GetPointer(self.PC+1, 2))

	self.PC = uint32(int32(self.PC) + int32(*p))
}

/******************************************************************
    功能:ILoad
	入参:1、*METHOD_FRAME
    返回值:无
******************************************************************/
func (self *METHOD_STACK) ILoad(frame *METHOD_FRAME) {
	self.PC++
	frame.Push(frame.GetVar(uint32(Memory[self.PC])))
	self.PC++
}

/******************************************************************
    功能:Load32
	入参:1、*METHOD_FRAME
	    2、局部变量索引
    返回值:无
******************************************************************/
func (self *METHOD_STACK) Load32(frame *METHOD_FRAME, index uint32) {
	self.PC++
	frame.Push(frame.GetVar(index))
}

/******************************************************************
    功能:Load32
	入参:1、*METHOD_FRAME
	    2、局部变量索引
    返回值:无
******************************************************************/
func (self *METHOD_STACK) Load64(frame *METHOD_FRAME, index uint32) {
	self.PC++
	frame.Push(frame.GetVar(index))
	frame.Push(frame.GetVar(index + 1))
}

/******************************************************************
    功能:IAStore
	入参:1、*METHOD_FRAME
    返回值:无
******************************************************************/
func (self *METHOD_STACK) IAStore(frame *METHOD_FRAME) {
	self.PC++
	v := frame.Pop()
	index := int32(frame.Pop())
	arrRef := frame.Pop()
	if arrRef == 0 {
		panic("IAStore() null")
	}
	arrInfo, context := GetArrayInfo(arrRef)
	if index < 0 || index >= int32(arrInfo.Length) {
		panic("IAStore()")
	}

	p := (*uint32)(BytesToUnsafePointer(context[index*4 : index*4+4]))
	*p = v
}

/******************************************************************
    功能:BAStore
	入参:1、*METHOD_FRAME
    返回值:无
******************************************************************/
func (self *METHOD_STACK) BAStore(frame *METHOD_FRAME) {
	self.PC++
	v := frame.Pop()
	index := int32(frame.Pop())
	arrRef := frame.Pop()
	if arrRef == 0 {
		panic("IAStore() null")
	}
	arrInfo, context := GetArrayInfo(arrRef)
	if index < 0 || index >= int32(arrInfo.Length) {
		panic("IAStore()")
	}

	context[index] = uint8(v)
}

/******************************************************************
    功能:IALoad
	入参:1、*METHOD_FRAME
    返回值:无
******************************************************************/
func (self *METHOD_STACK) IALoad(frame *METHOD_FRAME) {
	self.PC++
	index := int32(frame.Pop())
	arrRef := frame.Pop()
	if arrRef == 0 {
		panic("IALoad() null")
	}
	arrInfo, context := GetArrayInfo(arrRef)
	if index < 0 || index >= int32(arrInfo.Length) {
		panic("IALoad()")
	}
	p := (*int32)(BytesToUnsafePointer(context[index*4 : index*4+4]))
	frame.Push(uint32(*p))
}

/******************************************************************
    功能:BALoad
	入参:1、*METHOD_FRAME
	    2、局部变量索引
    返回值:无
******************************************************************/
func (self *METHOD_STACK) BALoad(frame *METHOD_FRAME) {
	self.PC++
	index := int32(frame.Pop())
	arrRef := frame.Pop()
	if arrRef == 0 {
		panic("BALoad() null")
	}
	arrInfo, context := GetArrayInfo(arrRef)
	if index < 0 || index >= int32(arrInfo.Length) {
		panic("BALoad()")
	}
	frame.Push(uint32(context[index]))
}

/******************************************************************
    功能:Store32Byte
	入参:1、*METHOD_FRAME
	    2、局部变量索引
    返回值:无
******************************************************************/
func (self *METHOD_STACK) Store32Byte(frame *METHOD_FRAME) {
	self.PC++
	frame.SetVar(uint32(Memory[self.PC]), frame.Pop())
	self.PC++
}

/******************************************************************
    功能:Store32
	入参:1、*METHOD_FRAME
	    2、局部变量索引
    返回值:无
******************************************************************/
func (self *METHOD_STACK) IStore(frame *METHOD_FRAME) {
	self.PC++
	frame.SetVar(uint32(Memory[self.PC]), frame.Pop())
	self.PC++
}

/******************************************************************
    功能:Store32
	入参:1、*METHOD_FRAME
	    2、局部变量索引
    返回值:无
******************************************************************/
func (self *METHOD_STACK) Store32(frame *METHOD_FRAME, index uint32) {
	self.PC++
	frame.SetVar(index, frame.Pop())
}

/******************************************************************
    功能:Store64
	入参:1、*METHOD_FRAME
	    2、局部变量索引
    返回值:无
******************************************************************/
func (self *METHOD_STACK) Store64(frame *METHOD_FRAME, index uint32) {
	self.PC++
	frame.SetVar(index+1, frame.Pop())
	frame.SetVar(index, frame.Pop())
}

/******************************************************************
    功能:Aconst_Null
	入参:1、*METHOD_FRAME
	    2、值
    返回值:无
******************************************************************/
func (self *METHOD_STACK) Aconst_Null(frame *METHOD_FRAME) {
	self.PC++
	frame.Push(uint32(0))
}

/******************************************************************
    功能:IConst
	入参:1、*METHOD_FRAME
	    2、值
    返回值:无
******************************************************************/
func (self *METHOD_STACK) IConst(frame *METHOD_FRAME, value int32) {
	self.PC++
	frame.Push(uint32(value))
}

/******************************************************************
    功能:IfNull
	入参:*METHOD_FRAME
    返回值:无
******************************************************************/
func (self *METHOD_STACK) IfNull(frame *METHOD_FRAME) {
	p := (*int16)(GetPointer(self.PC+1, 2))
	v := uint32(frame.Pop())
	if v == 0 {
		self.PC = uint32(int32(self.PC) + int32(*p))
	} else {
		self.PC += 3
	}
}

/******************************************************************
    功能:计算参数大小，不包括this
	入参:函数描述符
    返回值:参数大小
******************************************************************/
func CalParaSize(desp string) uint32 {
	runes := []rune(desp)
	var num uint32 = 0
	isAccess := false
	for i := 0; i < len(runes); i++ {
		if runes[i] == '(' {
			continue
		}
		if runes[i] == ')' {
			break
		}
		if isAccess {
			if runes[i] == ';' {
				isAccess = false
			}
			continue
		}
		switch runes[i] {
		case 'B', 'C', 'F', 'I', 'S', 'Z':
			num++
		case 'D', 'J':
			num += 2
		case '[', 'L':
			num++
			isAccess = true
		}
	}
	return num
}

/******************************************************************
    功能:桩代码
	入参:无
    返回值:无
******************************************************************/
func StubInvokeFunc(frame *METHOD_FRAME, methodRef MethodInfo) bool {
	//System.out.println函数打桩
	if methodRef.ClassName == SYM_java_io_PrintStream &&
		methodRef.MethodName == SYM_println &&
		methodRef.MethodDesp == SYM_Ljava_lang_String_V {
		strAccess := frame.Pop()
		strInst := (*STRING)(BytesToUnsafePointer(GetData(strAccess)))
		_, context := GetArrayInfo(strInst.ArrAdr)
		utf16_str := *(*[]uint16)(BytesToArray(context, 2))
		fmt.Println(string(utf16.Decode(utf16_str)))

		fmt.Println("STUB:", string(GetSymbol(methodRef.ClassName)),
			string(GetSymbol(methodRef.MethodName)),
			string(GetSymbol(methodRef.MethodDesp)))

		return true
	}
	if methodRef.ClassName == SYM_java_io_PrintStream &&
		methodRef.MethodName == SYM_println &&
		methodRef.MethodDesp == SYM_S_V {
		fmt.Println()

		fmt.Println("STUB:", string(GetSymbol(methodRef.ClassName)),
			string(GetSymbol(methodRef.MethodName)),
			string(GetSymbol(methodRef.MethodDesp)))

		return true
	}
	return false
}

/******************************************************************
    功能:Log
	入参:无
    返回值:无
******************************************************************/
func (self *METHOD_STACK) Log(frame *METHOD_FRAME) {
	fmt.Println("---", self.TopFrame, self.StackNum, "---")
	fmt.Print("Stack: ")
	for i := uint32(0); i < frame.CurOpStackIndex; i++ {
		p := (*int32)(GetPointer(frame.LocalAdr+METHOD_FRAME_SIZE+(frame.VarSize+i)*4, 4))
		fmt.Print(*p, " ")
	}
	fmt.Println()
	fmt.Printf("LocalVer: ")
	for i := uint32(0); i < frame.VarSize; i++ {
		p := (*int32)(GetPointer(frame.LocalAdr+METHOD_FRAME_SIZE+i*4, 4))
		fmt.Print(*p, " ")
	}
	fmt.Println()
}

/******************************************************************
    功能:执行static代码块
	入参:无
    返回值:无
******************************************************************/
func CInit(adr uint32) {
	classInfo := (*CLASS_INFO)(GetPointer(adr, CLASS_INFO_SIZE))
	if classInfo.IsCInit {
		return
	}
	classInfo.IsCInit = true
	if classInfo.SuperClassAddr != INVALID_MEM {
		CInit(classInfo.SuperClassAddr)
	}
	methodInfo, codeAdr := classInfo.FindMethod(SYM_CINIT, SYM_S_V)
	if methodInfo == nil || codeAdr == INVALID_MEM {
		return
	}
	fmt.Println(string(GetSymbol(classInfo.ClassName)), "static{}")
	codeAttri := (*CODE_ATTRI)(GetPointer(codeAdr, CODE_ATTRI_SIZE))
	//创建方法栈
	methodStack := NewMethodStack()
	methodStack.PushFrame(codeAttri.MaxLocal, codeAttri.MaxStack, classInfo.LocalAdr, 0)
	methodStack.PC = codeAdr + CODE_ATTRI_SIZE
	methodStack.Excute()
}

/******************************************************************
    功能:执行static代码块
	入参:无
    返回值:无
******************************************************************/
func (self *METHOD_STACK) ReflectMethod(frame *METHOD_FRAME, methodRef *MethodInfo) *METHOD_FRAME {
	//获取this中的类
	classAdr := GetClassMemAddr(methodRef.ClassName)
	var classInfo *CLASS_INFO
	var err error
	if classAdr == INVALID_MEM {
		classInfo, err = LoadClass(string(GetSymbol(methodRef.ClassName)))
		if err != nil {
			panic("ReflectMethod()")
		}
		CInit(classInfo.LocalAdr)
	} else {
		classInfo = (*CLASS_INFO)(GetPointer(classAdr, CLASS_INFO_SIZE))
	}
	//查找方法
	methodInfo, codeAdr := classInfo.FindMethod(methodRef.MethodName, methodRef.MethodDesp)
	if methodInfo.AccessFlag&METHOD_ACC_NATIVE == METHOD_ACC_NATIVE {
		ExcuteLocalMethodAdp(methodRef.ClassName, methodRef.MethodName, methodRef.MethodDesp, frame)

		return frame
	} else if codeAdr == INVALID_MEM {
		panic("方法没找到!")
	}
	fmt.Println(string(GetSymbol(methodRef.ClassName)),
		string(GetSymbol(methodRef.MethodName)),
		string(GetSymbol(methodRef.MethodDesp)))

	codeAttri := (*CODE_ATTRI)(GetPointer(codeAdr, CODE_ATTRI_SIZE))
	//创建方法栈
	newFrame := self.PushFrame(codeAttri.MaxLocal, codeAttri.MaxStack, classInfo.LocalAdr, self.PC)
	self.PC = codeAdr + CODE_ATTRI_SIZE

	//计算需要弹出的参数个数
	num := CalParaSize(string(GetSymbol(methodRef.MethodDesp)))
	param := make([]uint32, num)
	//将上一个栈帧中的值弹出，保存到新的栈帧中的局部变量中
	for i := num; i > 0; i-- {
		param[i-1] = frame.Pop()
	}
	//获取this中的类
	newFrame.SetVar(0, frame.Pop())
	for i, k := range param {
		newFrame.SetVar(uint32(i+1), k)
	}
	return newFrame
}
