package symbol

import (
	. "basic/memCtrl"
	"bytes"
	"errors"
)

//符号表信息
type SymbolItem struct {
	Next   uint32
	Length uint32
}

const SYMBOL_HEADER_SIZE = 8 //next指针(4) + 符号长度(4)

var SYM_J uint32 //long型的符号
var SYM_D uint32 //double型的符号

var SYM_Kjava_lang_String uint32
var SYM_KZ uint32 //boolean
var SYM_KB uint32
var SYM_KC uint32
var SYM_KS uint32
var SYM_KI uint32
var SYM_KJ uint32 //long
var SYM_KF uint32
var SYM_KD uint32

var SYM_java_io_PrintStream uint32
var SYM_println uint32
var SYM_print uint32
var SYM_Ljava_lang_String_V uint32
var SYM_S_V uint32
var SYM_CINIT uint32
var SYM_JAVA_LANG_CLASS uint32
var symHeaderAdr uint32

/***********************************
 初始化函数
************************************/
func init() {
	//初始化符号表头结点
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
	SYM_print, err = PutSymbol([]byte("print"))
	if err != nil {
		panic("")
	}

	SYM_Ljava_lang_String_V, err = PutSymbol([]byte("(Ljava/lang/String;)V"))
	if err != nil {
		panic("")
	}
	SYM_Kjava_lang_String, err = PutSymbol([]byte("[java/lang/String"))
	if err != nil {
		panic("")
	}
	SYM_KZ, err = PutSymbol([]byte("[Z"))
	if err != nil {
		panic("")
	}
	SYM_KB, err = PutSymbol([]byte("[B"))
	if err != nil {
		panic("")
	}
	SYM_KC, err = PutSymbol([]byte("[C"))
	if err != nil {
		panic("")
	}
	SYM_KS, err = PutSymbol([]byte("[S"))
	if err != nil {
		panic("")
	}
	SYM_KI, err = PutSymbol([]byte("[I"))
	if err != nil {
		panic("")
	}
	SYM_KJ, err = PutSymbol([]byte("[J"))
	if err != nil {
		panic("")
	}
	SYM_KF, err = PutSymbol([]byte("[F"))
	if err != nil {
		panic("")
	}
	SYM_KD, err = PutSymbol([]byte("[D"))
	if err != nil {
		panic("")
	}
	SYM_S_V, err = PutSymbol([]byte("()V"))
	if err != nil {
		panic("")
	}
	SYM_CINIT, err = PutSymbol([]byte("<clinit>"))
	if err != nil {
		panic("")
	}
	SYM_JAVA_LANG_CLASS, err = PutSymbol([]byte("java/lang/Class"))
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
