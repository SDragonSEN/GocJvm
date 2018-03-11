package memCtrl

import (
	"bytes"
	"errors"

	"../comFunc"
)

/***********************************
 添加符号，返回地址
************************************/
func PutSymbol(symbol []byte) (uint32, error) {
	curAddr := symHeaderAdr
	var curSymbol *SymbolItem
	for curAddr != INVALID_MEM {
		curSymbol = (*SymbolItem)(comFunc.BytesToUnsafePointer(Memory[curAddr:]))
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
	newSym := (*SymbolItem)(comFunc.BytesToUnsafePointer(Memory[newSymAdr:]))
	newSym.Length = uint32(len(symbol))
	newSym.Next = INVALID_MEM
	copy(Memory[newSymAdr+SYMBOL_HEADER_SIZE:newSymAdr+SYMBOL_HEADER_SIZE+newSym.Length], symbol)
	return newSymAdr, nil
}

/***********************************
 获取符号
************************************/
func GetSymbol(adr uint32) []byte {
	symbol := (*SymbolItem)(comFunc.BytesToUnsafePointer(Memory[adr:]))
	return Memory[adr+SYMBOL_HEADER_SIZE : adr+SYMBOL_HEADER_SIZE+symbol.Length]
}
