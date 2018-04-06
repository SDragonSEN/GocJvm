package method

import (
	"../memoryControl"
)

type METHOD_STACK struct {
	MaxStackSize uint32
	StackNum     uint32
	TopFrame     uint32 //地址
	PC           uint32
}

const METHOD_STACK_SIZE = 3 * 4

type METHOD_FRAME struct {
	LowFrame        uint32 //地址
	Claz            uint32 //类地址
	ReturnPc        uint32 //返回后的PC地址
	VarSize         uint32 //变量区大小
	OpStackSize     uint32 //操作数栈大小
	CurOpStackIndex uint32 //当前操作数栈索引
}

const METHOD_FRAME_SIZE = 6 * 4

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
func (self *METHOD_STACK) PushFrame(varSize, opStackSize, clazAdr, returnPc uint32) (*METHOD_FRAME, uint32) {
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
	return methodFrame, adr
}
