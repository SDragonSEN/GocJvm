package memCtrl

import (
	"bytes"
)

/***********************************
 添加符号，返回地址
************************************/
func PutSymbol(symbol []byte) uint32 {
	curAddr := symHeaderAdr
	var curSymbol *SymbolItem
	for curAddr != INVALID_MEM {
		curSymbol = (*SymbolItem)(BytesToUnsafePointer(memory[curAddr:]))
		if curSymbol.Length == uint32(len(symbol)) && 0 == bytes.Compare(memory[curAddr+SYMBOL_HEADER_SIZE:curAddr+SYMBOL_HEADER_SIZE+curSymbol.Length], symbol) {
			return curAddr
		}
		curAddr = curSymbol.Next
	}

	newSymAdr, _ := Malloc(uint32(SYMBOL_HEADER_SIZE+len(symbol)), SYMBOL_NODE)
	curSymbol.Next = newSymAdr
	newSym := (*SymbolItem)(BytesToUnsafePointer(memory[newSymAdr:]))
	newSym.Length = uint32(len(symbol))
	newSym.Next = INVALID_MEM
	copy(memory[newSymAdr+SYMBOL_HEADER_SIZE:newSymAdr+SYMBOL_HEADER_SIZE+newSym.Next], symbol)
	return newSymAdr
}

/***********************************
 获取符号
************************************/
func getSymbol(adr uint32) []byte {
	symbol := (*SymbolItem)(BytesToUnsafePointer(memory[adr:]))
	return memory[adr+SYMBOL_HEADER_SIZE : adr+SYMBOL_HEADER_SIZE+symbol.Length]
}
