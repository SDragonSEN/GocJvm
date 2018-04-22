package method

import (
	"fmt"
	"unicode/utf16"

	"accessOp"
	"classAnaly"
	"comFunc"
	"comValue"
	"memoryControl"
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
	adr, err := memCtrl.Malloc(METHOD_STACK_SIZE, memCtrl.METHOD_STACK_NODE)
	if err != nil {
		panic("NewMethodStack()")
	}
	methodStack := (*METHOD_STACK)(memCtrl.GetPointer(adr, METHOD_STACK_SIZE))
	methodStack.MaxStackSize = 100 //目前默认100
	methodStack.StackNum = 0
	methodStack.TopFrame = memCtrl.INVALID_MEM
	return methodStack
}

/******************************************************************
    功能:方法压栈
	入参:无
    返回值:1、*METHOD_FRAME
	      2、地址
******************************************************************/
func (self *METHOD_STACK) PushFrame(varSize, opStackSize, clazAdr, returnPc uint32) *METHOD_FRAME {
	adr, err := memCtrl.Malloc(METHOD_FRAME_SIZE+(varSize+opStackSize)*4, memCtrl.METHOD_FRAME_NODE)
	if err != nil {
		panic("PushFrame()")
	}
	methodFrame := (*METHOD_FRAME)(memCtrl.GetPointer(adr, METHOD_FRAME_SIZE))
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
	curFrame := (*METHOD_FRAME)(memCtrl.GetPointer(curFrameAdr, METHOD_FRAME_SIZE))

	self.TopFrame = curFrame.LowFrame
	self.PC = curFrame.ReturnPc
	self.StackNum--
	memCtrl.MemFree(curFrameAdr)
	if self.StackNum == 0 {
		return nil
	}
	return (*METHOD_FRAME)(memCtrl.GetPointer(self.TopFrame, METHOD_FRAME_SIZE))
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
	p := (*uint32)(memCtrl.GetPointer(self.LocalAdr+METHOD_FRAME_SIZE+(index*4), 4))
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
	p := (*uint32)(memCtrl.GetPointer(self.LocalAdr+METHOD_FRAME_SIZE+(index*4), 4))
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
	p := (*uint32)(memCtrl.GetPointer(self.LocalAdr+METHOD_FRAME_SIZE+(self.VarSize+self.CurOpStackIndex)*4, 4))
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
	p := (*uint32)(memCtrl.GetPointer(self.LocalAdr+METHOD_FRAME_SIZE+(self.VarSize+self.CurOpStackIndex-1)*4, 4))
	self.CurOpStackIndex--
	return *p
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
	p := (*uint32)(memCtrl.GetPointer(self.LocalAdr+METHOD_FRAME_SIZE+(self.VarSize+self.CurOpStackIndex-1)*4, 4))
	return *p
}

/******************************************************************
    功能:执行函数
	入参:无
    返回值:无
******************************************************************/
func (self *METHOD_STACK) Excute() {
	i := 0
	frame := (*METHOD_FRAME)(memCtrl.GetPointer(self.TopFrame, METHOD_FRAME_SIZE))
	for {
		self.Log(frame)
		fmt.Println(comValue.Format(memCtrl.Memory[self.PC]))
		fmt.Println(self.PC)
		switch memCtrl.Memory[self.PC] {
		case comValue.NOP:
			if i > 20 {
				return
			}
			i++
			self.PC++
		case comValue.ICONST_M1:
			self.IConst(frame, -1)
		case comValue.ICONST_0:
			self.IConst(frame, 0)
		case comValue.ICONST_1:
			self.IConst(frame, 1)
		case comValue.ICONST_2:
			self.IConst(frame, 2)
		case comValue.ICONST_3:
			self.IConst(frame, 3)
		case comValue.ICONST_4:
			self.IConst(frame, 4)
		case comValue.ICONST_5:
			self.IConst(frame, 5)
		case comValue.BIPUSH:
			self.BIPush(frame)
		case comValue.LDC:
			self.Ldc(frame)
		case comValue.ILOAD:
			self.ILoad(frame)
		case comValue.ALOAD:
			self.ILoad(frame)
		case comValue.ILOAD_0:
			self.Load32(frame, 0)
		case comValue.ILOAD_1:
			self.Load32(frame, 1)
		case comValue.ILOAD_2:
			self.Load32(frame, 2)
		case comValue.ILOAD_3:
			self.Load32(frame, 3)
		case comValue.ALOAD_0:
			self.Load32(frame, 0)
		case comValue.ALOAD_1:
			self.Load32(frame, 1)
		case comValue.ALOAD_2:
			self.Load32(frame, 2)
		case comValue.IALOAD:
			self.IALoad(frame)
		case comValue.BALOAD:
			self.BALoad(frame)
		case comValue.ASTORE:
			self.Store32Byte(frame)
		case comValue.ISTORE:
			self.IStore(frame)
		case comValue.ISTORE_1:
			self.Store32(frame, 1)
		case comValue.ISTORE_2:
			self.Store32(frame, 2)
		case comValue.ISTORE_3:
			self.Store32(frame, 3)
		case comValue.ASTORE_1:
			self.Store32(frame, 1)
		case comValue.ASTORE_2:
			self.Store32(frame, 2)
		case comValue.IASTORE:
			self.IAStore(frame)
		case comValue.BASTORE:
			self.BAStore(frame)
		case comValue.POP:
			frame.Pop()
			self.PC++
		case comValue.DUP:
			self.Dup(frame)
		case comValue.IADD:
			self.Iadd(frame)
		case comValue.ISUB:
			self.Isub(frame)
		case comValue.IMUL:
			self.Imul(frame)
		case comValue.IDIV:
			self.Idiv(frame)
		case comValue.IREM:
			self.Irem(frame)
		case comValue.INEG:
			self.Ineg(frame)
		case comValue.IFGE:
			self.Ifge(frame)
		case comValue.IFNE:
			self.Ifne(frame)
		case comValue.IF_ICMPNE:
			self.Icmpne(frame)
		case comValue.IF_ICMPLT:
			self.Icmplt(frame)
		case comValue.IF_ICMPGT:
			self.Icmpgt(frame)
		case comValue.IF_ICMPLE:
			self.Icmple(frame)

		case comValue.IRETURN:
			frame = self.IReturn(frame)
			if frame == nil {
				goto label
			}
		case comValue.GOTO:
			self.Goto(frame)
		case comValue.GETSTATIC:
			self.GetStatic(frame)
		case comValue.GETFIELD:
			self.GetFiled(frame)
		case comValue.PUTFIELD:
			self.PutFiled(frame)
		case comValue.INVOKEVIRTUAL:
			newFrame := self.InvokeVirtual(frame)
			if newFrame != nil {
				frame = newFrame
			}
		case comValue.INVOKESTATIC:
			frame = self.InvokeStatic(frame)
		case comValue.INVOKESPECIAL:
			frame = self.InvokeSpecial(frame)
		case comValue.RETURN:
			frame = self.PopFrame()
			if frame == nil {
				goto label
			}
		case comValue.NEW:
			self.New(frame)
		case comValue.ARRAYLENGTH:
			self.ArrayLength(frame)
		case comValue.NEWARRAY:
			self.NewArray(frame)
		default:
			fmt.Printf("memCtrl.Memory[self.PC]:%x ", memCtrl.Memory[self.PC])
			fmt.Println(comValue.Format(memCtrl.Memory[self.PC]))
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
	p := (*uint16)(memCtrl.GetPointer(self.PC, 2))
	filedInfo := classAnaly.GetFiledInfo(classAnaly.GetConstantPoolSlice(frame.Claz), uint32(*p))
	//fmt.Println(string(memCtrl.GetSymbol(filedInfo.ClassName)), string(memCtrl.GetSymbol(filedInfo.FiledName)), string(memCtrl.GetSymbol(filedInfo.FiledType)))
	self.PC += 2

	var classAdr uint32
	classAdr = memCtrl.GetClassMemAddr(filedInfo.ClassName)
	//如果获取不到，则说明不在内存中，需要去加载
	if classAdr == memCtrl.INVALID_MEM {
		//获取类名(string)
		className := string(memCtrl.GetSymbol(filedInfo.ClassName))
		classInfo, err := classAnaly.LoadClass(className)
		if err != nil {
			panic("GetStatic()")
		}
		classAdr = classInfo.LocalAdr
	}

	classInfo := (*classAnaly.CLASS_INFO)(memCtrl.GetPointer(classAdr, classAnaly.CLASS_INFO_SIZE))
	//判断是否是long或double型
	if filedInfo.FiledType == memCtrl.SYM_J ||
		filedInfo.FiledType == memCtrl.SYM_D {
		v := classInfo.GetStaticData64(filedInfo.FiledName, filedInfo.FiledType)
		frame.Push(v[0])
		frame.Push(v[1])
	} else {
		v := classInfo.GetStaticData32(filedInfo.FiledName, filedInfo.FiledType)
		frame.Push(v)
	}
}

/******************************************************************
    功能:GetFiled
	入参:1、*METHOD_FRAME
    返回值:无
******************************************************************/
func (self *METHOD_STACK) GetFiled(frame *METHOD_FRAME) {
	self.PC++
	p := (*uint16)(memCtrl.GetPointer(self.PC, 2))
	filedInfo := classAnaly.GetFiledInfo(classAnaly.GetConstantPoolSlice(frame.Claz), uint32(*p))
	//	fmt.Println(string(memCtrl.GetSymbol(filedInfo.FiledName)))
	self.PC += 2
	accessAdr := frame.Pop()
	this := (*access.ACCESS_INFO)(memCtrl.GetPointer(accessAdr, access.ACCESS_INFO_SIZE))
	thisClass := (*classAnaly.CLASS_INFO)(memCtrl.GetPointer(this.TypeAddr, classAnaly.CLASS_INFO_SIZE))
	index := thisClass.GetUnstaticDataIndex(filedInfo.FiledName, filedInfo.FiledType)
	data := access.GetData(accessAdr)
	v0 := (*uint32)(comFunc.BytesToUnsafePointer(data[index*4 : index*4+4]))
	frame.Push(*v0)
	if filedInfo.FiledType == memCtrl.SYM_J ||
		filedInfo.FiledType == memCtrl.SYM_D {
		v1 := (*uint32)(comFunc.BytesToUnsafePointer(data[index*4+4 : index*4+8]))
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
	p := (*uint16)(memCtrl.GetPointer(self.PC, 2))
	filedInfo := classAnaly.GetFiledInfo(classAnaly.GetConstantPoolSlice(frame.Claz), uint32(*p))
	self.PC += 2
	if filedInfo.FiledType == memCtrl.SYM_J ||
		filedInfo.FiledType == memCtrl.SYM_D {
		v1 := frame.Pop()
		v0 := frame.Pop()
		accessAdr := frame.Pop()
		this := (*access.ACCESS_INFO)(memCtrl.GetPointer(accessAdr, access.ACCESS_INFO_SIZE))
		thisClass := (*classAnaly.CLASS_INFO)(memCtrl.GetPointer(this.TypeAddr, classAnaly.CLASS_INFO_SIZE))
		index := thisClass.GetUnstaticDataIndex(filedInfo.FiledName, filedInfo.FiledType)
		data := access.GetData(accessAdr)
		p0 := (*uint32)(comFunc.BytesToUnsafePointer(data[index*4 : index*4+4]))
		*p0 = v0
		p1 := (*uint32)(comFunc.BytesToUnsafePointer(data[index*4+4 : index*4+8]))
		*p1 = v1
	} else {
		v := frame.Pop()
		accessAdr := frame.Pop()
		this := (*access.ACCESS_INFO)(memCtrl.GetPointer(accessAdr, access.ACCESS_INFO_SIZE))
		thisClass := (*classAnaly.CLASS_INFO)(memCtrl.GetPointer(this.TypeAddr, classAnaly.CLASS_INFO_SIZE))
		index := thisClass.GetUnstaticDataIndex(filedInfo.FiledName, filedInfo.FiledType)
		data := access.GetData(accessAdr)
		p := (*uint32)(comFunc.BytesToUnsafePointer(data[index*4 : index*4+4]))
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
	v := (*int8)(memCtrl.GetPointer(self.PC, 1))
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
	v := classAnaly.GetStringFromConstPool(classAnaly.GetConstantPoolSlice(frame.Claz), uint32(memCtrl.Memory[self.PC]))
	frame.Push(v)
	self.PC++
}

/******************************************************************
    功能:InvokeVirtual指令
	入参:无
    返回值:无
******************************************************************/
func (self *METHOD_STACK) InvokeVirtual(frame *METHOD_FRAME) *METHOD_FRAME {
	self.PC++
	p := (*uint16)(memCtrl.GetPointer(self.PC, 2))
	self.PC += 2
	methodRef := classAnaly.GetMethodInfo(classAnaly.GetConstantPoolSlice(frame.Claz), uint32(*p))
	if StubInvokeFunc(frame, methodRef) {
		return nil
	}

	num := CalParaSize(string(memCtrl.GetSymbol(methodRef.MethodDesp)))
	param := make([]uint32, num)
	//将上一个栈帧中的值弹出，保存到新的栈帧中的局部变量中
	for i := num; i > 0; i-- {
		param[i-1] = frame.Pop()
	}
	//获取this中的类
	this := frame.Pop()
	classInfo := (*classAnaly.CLASS_INFO)(memCtrl.GetPointer(access.GetClassInfo(this), classAnaly.CLASS_INFO_SIZE))

	fmt.Println(string(memCtrl.GetSymbol(methodRef.ClassName)),
		string(memCtrl.GetSymbol(methodRef.MethodName)),
		string(memCtrl.GetSymbol(methodRef.MethodDesp)))

	if StubInvokeFunc(frame, methodRef) {
		return nil
	}

	//查找方法
	methodInfo, methodCLass, codeAdr := classInfo.FindMethodEx(methodRef.MethodName, methodRef.MethodDesp)
	if methodInfo == nil || codeAdr == memCtrl.INVALID_MEM {
		fmt.Println(methodInfo, codeAdr)
		fmt.Println(string(memCtrl.GetSymbol(methodRef.ClassName)),
			string(memCtrl.GetSymbol(methodRef.MethodName)),
			string(memCtrl.GetSymbol(methodRef.MethodDesp)))
		panic("InvokeVirtual()6")
	}

	codeAttri := (*classAnaly.CODE_ATTRI)(memCtrl.GetPointer(codeAdr, classAnaly.CODE_ATTRI_SIZE))
	//创建方法栈
	newFrame := self.PushFrame(codeAttri.MaxLocal, codeAttri.MaxStack, methodCLass.LocalAdr, self.PC)
	self.PC = codeAdr + classAnaly.CODE_ATTRI_SIZE
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
	p := (*uint16)(memCtrl.GetPointer(self.PC, 2))
	self.PC += 2
	//获取方法描述
	methodRef := classAnaly.GetMethodInfo(classAnaly.GetConstantPoolSlice(frame.Claz), uint32(*p))

	//获取方法中的类
	var classInfo *classAnaly.CLASS_INFO
	var err error
	classAdr := memCtrl.GetClassMemAddr(methodRef.ClassName)
	if classAdr == memCtrl.INVALID_MEM {
		classInfo, err = classAnaly.LoadClass(string(memCtrl.GetSymbol(methodRef.ClassName)))
		if err != nil {
			panic("InvokeSpecial()")
		}
	} else {
		classInfo = (*classAnaly.CLASS_INFO)(memCtrl.GetPointer(classAdr, classAnaly.CLASS_INFO_SIZE))
	}
	//查找方法
	methodInfo, codeAdr := classInfo.FindMethod(methodRef.MethodName, methodRef.MethodDesp)
	if methodInfo == nil || codeAdr == memCtrl.INVALID_MEM {
		fmt.Println(methodInfo, codeAdr)
		fmt.Println(string(memCtrl.GetSymbol(classInfo.ClassName)),
			string(memCtrl.GetSymbol(methodRef.MethodName)),
			string(memCtrl.GetSymbol(methodRef.MethodDesp)))
		panic("InvokeSpecial()6")
	}
	fmt.Println(string(memCtrl.GetSymbol(methodRef.ClassName)),
		string(memCtrl.GetSymbol(methodRef.MethodName)),
		string(memCtrl.GetSymbol(methodRef.MethodDesp)))

	codeAttri := (*classAnaly.CODE_ATTRI)(memCtrl.GetPointer(codeAdr, classAnaly.CODE_ATTRI_SIZE))
	//创建方法栈
	newFrame := self.PushFrame(codeAttri.MaxLocal, codeAttri.MaxStack, classInfo.LocalAdr, self.PC)
	self.PC = codeAdr + classAnaly.CODE_ATTRI_SIZE

	//计算需要弹出的参数个数
	num := CalParaSize(string(memCtrl.GetSymbol(methodRef.MethodDesp)))

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
	p := (*uint16)(memCtrl.GetPointer(self.PC, 2))
	self.PC += 2
	//获取方法描述
	methodRef := classAnaly.GetMethodInfo(classAnaly.GetConstantPoolSlice(frame.Claz), uint32(*p))

	//获取方法中的类
	var classInfo *classAnaly.CLASS_INFO
	var err error
	classAdr := memCtrl.GetClassMemAddr(methodRef.ClassName)
	if classAdr == memCtrl.INVALID_MEM {
		classInfo, err = classAnaly.LoadClass(string(memCtrl.GetSymbol(methodRef.ClassName)))
		if err != nil {
			panic("InvokeSpecial()")
		}
	} else {
		classInfo = (*classAnaly.CLASS_INFO)(memCtrl.GetPointer(classAdr, classAnaly.CLASS_INFO_SIZE))
	}
	fmt.Println(string(memCtrl.GetSymbol(classInfo.ClassName)),
		string(memCtrl.GetSymbol(methodRef.MethodName)),
		string(memCtrl.GetSymbol(methodRef.MethodDesp)))
	//查找方法
	methodInfo, methodCLass, codeAdr := classInfo.FindMethodEx(methodRef.MethodName, methodRef.MethodDesp)
	if methodInfo == nil || codeAdr == memCtrl.INVALID_MEM {
		fmt.Println(methodInfo, codeAdr)
		panic("InvokeSpecial()6")
	}
	codeAttri := (*classAnaly.CODE_ATTRI)(memCtrl.GetPointer(codeAdr, classAnaly.CODE_ATTRI_SIZE))
	//创建方法栈
	newFrame := self.PushFrame(codeAttri.MaxLocal, codeAttri.MaxStack, methodCLass.LocalAdr, self.PC)
	self.PC = codeAdr + classAnaly.CODE_ATTRI_SIZE

	//计算需要弹出的参数个数
	num := CalParaSize(string(memCtrl.GetSymbol(methodRef.MethodDesp)))
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
	p := (*uint16)(memCtrl.GetPointer(self.PC, 2))
	self.PC += 2
	className := classAnaly.GetClassFromConstPool(classAnaly.GetConstantPoolSlice(frame.Claz), uint32(*p))
	var classAdr uint32
	var classInfo *classAnaly.CLASS_INFO
	classAdr = memCtrl.GetClassMemAddr(className)
	//如果获取不到，则说明不在内存中，需要去加载
	if classAdr == memCtrl.INVALID_MEM {
		//获取类名(string)
		classNameStr := string(memCtrl.GetSymbol(className))
		classInfo, err := classAnaly.LoadClass(classNameStr)
		if err != nil {
			panic("GetStatic()")
		}
		classAdr = classInfo.LocalAdr
	}
	accessInfo, accessAdr, err := access.NewAccessInfo()
	if err != nil {
		panic(err)
	}
	classInfo = (*classAnaly.CLASS_INFO)(memCtrl.GetPointer(classAdr, classAnaly.CLASS_INFO_SIZE))
	accessInfo.TypeAddr = classAdr
	accessInfo.DataAddr, err = memCtrl.Malloc(classInfo.UnstaticParaTotalSize, memCtrl.INSTANCE_NODE)
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
	if arrAdr == memCtrl.INVALID_MEM {
		panic("null pointer")
	}
	arrInfo, _ := access.GetArrayInfo(arrAdr)
	frame.Push(arrInfo.Length)
}

/******************************************************************
    功能:NewArray
	入参:*METHOD_FRAME
    返回值:无
******************************************************************/
func (self *METHOD_STACK) NewArray(frame *METHOD_FRAME) {
	self.PC++
	p := (*uint8)(memCtrl.GetPointer(self.PC, 1))
	var symbol uint32
	var width uint32
	switch *p {
	case comValue.AT_BOOLEAN:
		symbol = memCtrl.SYM_KZ
		width = 1
	case comValue.AT_BYTE:
		symbol = memCtrl.SYM_KB
		width = 1
	case comValue.AT_CHAR:
		symbol = memCtrl.SYM_KC
		width = 2
	case comValue.AT_FLOAT:
		symbol = memCtrl.SYM_KF
		width = 4
	case comValue.AT_DOUBLE:
		symbol = memCtrl.SYM_KD
		width = 8
	case comValue.AT_SHORT:
		symbol = memCtrl.SYM_KS
		width = 2
	case comValue.AT_INT:
		symbol = memCtrl.SYM_KI
		width = 4
	case comValue.AT_LONG:
		symbol = memCtrl.SYM_KJ
		width = 8
	}
	self.PC++

	_, arrAdr, err := access.NewArray(symbol, width, frame.Pop())
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
    功能:Ifge
	入参:*METHOD_FRAME
    返回值:无
******************************************************************/
func (self *METHOD_STACK) Ifge(frame *METHOD_FRAME) {
	p := (*int16)(memCtrl.GetPointer(self.PC+1, 2))
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
	p := (*int16)(memCtrl.GetPointer(self.PC+1, 2))
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
	p := (*int16)(memCtrl.GetPointer(self.PC+1, 2))
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
	p := (*int16)(memCtrl.GetPointer(self.PC+1, 2))
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
	p := (*int16)(memCtrl.GetPointer(self.PC+1, 2))
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
	p := (*int16)(memCtrl.GetPointer(self.PC+1, 2))
	v0 := int32(frame.Pop())
	v1 := int32(frame.Pop())
	if v0 >= v1 {
		self.PC = uint32(int32(self.PC) + int32(*p))
	} else {
		self.PC += 3
	}
}

/******************************************************************
    功能:Ifge
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
    功能:Icmpgt
	入参:*METHOD_FRAME
    返回值:无
******************************************************************/
func (self *METHOD_STACK) Goto(frame *METHOD_FRAME) {
	p := (*int16)(memCtrl.GetPointer(self.PC+1, 2))

	self.PC = uint32(int32(self.PC) + int32(*p))
}

/******************************************************************
    功能:ILoad
	入参:1、*METHOD_FRAME
    返回值:无
******************************************************************/
func (self *METHOD_STACK) ILoad(frame *METHOD_FRAME) {
	self.PC++
	frame.Push(frame.GetVar(uint32(memCtrl.Memory[self.PC])))
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
	arrInfo, context := access.GetArrayInfo(arrRef)
	if index < 0 || index >= int32(arrInfo.Length) {
		panic("IAStore()")
	}

	p := (*uint32)(comFunc.BytesToUnsafePointer(context[index*4 : index*4+4]))
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
	arrInfo, context := access.GetArrayInfo(arrRef)
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
		panic("BALoad() null")
	}
	arrInfo, context := access.GetArrayInfo(arrRef)
	if index < 0 || index >= int32(arrInfo.Length) {
		panic("BALoad()")
	}
	p := (*int32)(comFunc.BytesToUnsafePointer(context[index*4 : index*4+4]))
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
	arrInfo, context := access.GetArrayInfo(arrRef)
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
	frame.SetVar(uint32(memCtrl.Memory[self.PC]), frame.Pop())
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
	frame.SetVar(uint32(memCtrl.Memory[self.PC]), frame.Pop())
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
func StubInvokeFunc(frame *METHOD_FRAME, methodRef classAnaly.MethodInfo) bool {
	//System.out.println函数打桩
	if methodRef.ClassName == memCtrl.SYM_java_io_PrintStream &&
		methodRef.MethodName == memCtrl.SYM_println &&
		methodRef.MethodDesp == memCtrl.SYM_Ljava_lang_String_V {
		strAccess := frame.Pop()
		strInst := (*access.STRING)(comFunc.BytesToUnsafePointer(access.GetData(strAccess)))
		_, context := access.GetArrayInfo(strInst.ArrAdr)
		utf16_str := *(*[]uint16)(comFunc.BytesToArray(context, 2))
		fmt.Println(string(utf16.Decode(utf16_str)))

		fmt.Println("STUB:", string(memCtrl.GetSymbol(methodRef.ClassName)),
			string(memCtrl.GetSymbol(methodRef.MethodName)),
			string(memCtrl.GetSymbol(methodRef.MethodDesp)))

		return true
	}
	if methodRef.ClassName == memCtrl.SYM_java_io_PrintStream &&
		methodRef.MethodName == memCtrl.SYM_println &&
		methodRef.MethodDesp == memCtrl.SYM_S_V {
		fmt.Println()

		fmt.Println("STUB:", string(memCtrl.GetSymbol(methodRef.ClassName)),
			string(memCtrl.GetSymbol(methodRef.MethodName)),
			string(memCtrl.GetSymbol(methodRef.MethodDesp)))

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
	fmt.Println("---", self.TopFrame, "---")
	fmt.Print("Stack: ")
	for i := uint32(0); i < frame.CurOpStackIndex; i++ {
		p := (*int32)(memCtrl.GetPointer(frame.LocalAdr+METHOD_FRAME_SIZE+(frame.VarSize+i)*4, 4))
		fmt.Print(*p, " ")
	}
	fmt.Println()
	fmt.Printf("LocalVer: ")
	for i := uint32(0); i < frame.VarSize; i++ {
		p := (*int32)(memCtrl.GetPointer(frame.LocalAdr+METHOD_FRAME_SIZE+i*4, 4))
		fmt.Print(*p, " ")
	}
	fmt.Println()
}
