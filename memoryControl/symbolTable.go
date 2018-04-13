package memCtrl

import (
	"bytes"
	"errors"
)

var SYM_J uint32 //long型的符号
var SYM_D uint32 //double型的符号

var SYM_java_io_PrintStream uint32
var SYM_println uint32
var SYM_Ljava_lang_String_V uint32

/***********************************
 初始化符号表
************************************/
func SymbolInit() {
	symHeaderAdr, _ = Malloc(SYMBOL_HEADER_SIZE, SYMBOL_NODE)
	symHeader := (*SymbolItem)(GetPointer(symHeaderAdr, SYMBOL_HEADER_SIZE))
	symHeader.Length = 0
	symHeader.Next = INVALID_MEM
	var err error
	SYM_J, err = PutSymbol([]byte("J"))
	if err != nil {
		panic("")
	}
	SYM_D, err = PutSymbol([]byte("D"))
	if err != nil {
		panic("")
	}
	SYM_java_io_PrintStream, err = PutSymbol([]byte("java/io/PrintStream"))
	if err != nil {
		panic("")
	}
	SYM_println, err = PutSymbol([]byte("println"))
	if err != nil {
		panic("")
	}
	SYM_Ljava_lang_String_V, err = PutSymbol([]byte("(Ljava/lang/String;)V"))
	if err != nil {
		panic("")
	}
}

/***********************************
 添加符号，返回地址
************************************/
func PutSymbol(symbol []byte) (uint32, error) {
	curAddr := symHeaderAdr
	var curSymbol *SymbolItem
	for curAddr != INVALID_MEM {
		curSymbol = (*SymbolItem)(GetPointer(curAddr, SYMBOL_HEADER_SIZE))
		if curSymbol.Length == uint32(len(symbol)) && 0 == bytes.Compare(Memory[curAddr+SYMBOL_HEADER_SIZE:curAddr+SYMBOL_HEADER_SIZE+curSymbol.Length], symbol) {
			return curAddr, nil
		}
		curAddr = curSymbol.Next
	}
	newSymAdr, err := Malloc(uint32(SYMBOL_HEADER_SIZE+len(symbol)), SYMBOL_NODE)
	if err != nil {
		return INVALID_MEM, errors.New("PutSymbol():内存不足")
	}
	curSymbol.Next = newSymAdr
	newSym := (*SymbolItem)(GetPointer(newSymAdr, SYMBOL_HEADER_SIZE))
	newSym.Length = uint32(len(symbol))
	newSym.Next = INVALID_MEM
	copy(Memory[newSymAdr+SYMBOL_HEADER_SIZE:newSymAdr+SYMBOL_HEADER_SIZE+newSym.Length], symbol)
	return newSymAdr, nil
}

/***********************************
 获取符号
************************************/
func GetSymbol(adr uint32) []byte {
	symbol := (*SymbolItem)(GetPointer(adr, SYMBOL_HEADER_SIZE))
	return Memory[adr+SYMBOL_HEADER_SIZE : adr+SYMBOL_HEADER_SIZE+symbol.Length]
}
