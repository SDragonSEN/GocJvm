package method

import (
	"../classAnaly"
	"../comValue"
	"../memoryControl"
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
	self.OpStackSize++
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
	self.OpStackSize--
	return *p
}
func (self *METHOD_STACK) Excute() {
	frame := (*METHOD_FRAME)(memCtrl.GetPointer(self.TopFrame, METHOD_FRAME_SIZE))
	switch memCtrl.Memory[self.PC] {
	case comValue.GETSTATIC:
		self.GetStatic(frame)
	}
}
func (self *METHOD_STACK) GetStatic(frame *METHOD_FRAME) {
	self.PC++
	p := (*uint16)(memCtrl.GetPointer(self.PC, 2))
	filedInfo := classAnaly.GetStaticFiledInfo(classAnaly.GetConstantPoolSlice(frame.Claz), uint32(*p))
	//fmt.Println(string(memCtrl.GetSymbol(filedInfo.ClassName)), string(memCtrl.GetSymbol(filedInfo.FiledName)), string(memCtrl.GetSymbol(filedInfo.FiledType)))
	self.PC += 2

	var superClassAdr uint32
	superClassAdr = memCtrl.GetClassMemAddr(filedInfo.ClassName)
	//如果获取不到，则说明不在内存中，需要去加载
	if superClassAdr == memCtrl.INVALID_MEM {
		//获取类名(string)
		className := string(memCtrl.GetSymbol(filedInfo.ClassName))
		superClass, err := classAnaly.LoadClass(className)
		if err != nil {
			panic("GetStatic()")
		}
		superClassAdr = superClass.LocalAdr
	}

	superClass := (*classAnaly.CLASS_INFO)(memCtrl.GetPointer(superClassAdr, classAnaly.CLASS_INFO_SIZE))
	//判断是否是long或double型
	if filedInfo.FiledType == memCtrl.SYM_J ||
		filedInfo.FiledType == memCtrl.SYM_D {
		v := superClass.GetStaticData64(filedInfo.FiledName, filedInfo.FiledType)
		frame.Push(v[0])
		frame.Push(v[1])
	} else {
		v := superClass.GetStaticData32(filedInfo.FiledName, filedInfo.FiledType)
		frame.Push(v)
	}
}
