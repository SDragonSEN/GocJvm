package memCtrl

import (
	"errors"
)

/***********************************
 添加class，返回地址
************************************/
func PutClass(className uint32, context []byte) (uint32, error) {
	curAddr := classHeaderAdr
	var curClass *ClassItem
	for curAddr != INVALID_MEM {
		curClass = (*ClassItem)(GetPointer(curAddr, CLASS_HEADER_SIZE))
		curAddr = curClass.Next
	}
	newClassAdr, err := Malloc(uint32(CLASS_HEADER_SIZE+len(context)), CLASS_NODE)
	if err != nil {
		return INVALID_MEM, errors.New("PutClassMemAddr():内存不足")
	}
	curClass.Next = newClassAdr
	newClass := (*ClassItem)(GetPointer(newClassAdr, CLASS_HEADER_SIZE))
	newClass.ClassName = className
	newClass.Next = INVALID_MEM
	copy(Memory[newClassAdr+CLASS_HEADER_SIZE:newClassAdr+SYMBOL_HEADER_SIZE+uint32(len(context))], context)
	return newClassAdr + CLASS_HEADER_SIZE, nil
}

/***********************************
 获取Class地址
************************************/
func GetClassMemAddr(className uint32) uint32 {
	curAddr := classHeaderAdr
	var curClass *ClassItem
	for curAddr != INVALID_MEM {
		curClass = (*ClassItem)(GetPointer(curAddr, CLASS_HEADER_SIZE))
		if curClass.ClassName == className {
			return curAddr + CLASS_HEADER_SIZE
		}
		curAddr = curClass.Next
	}
	return INVALID_MEM
}
