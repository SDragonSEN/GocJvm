package memCtrl

import (
	"../comFunc"
)

//节点头信息
type NodeHeader struct {
	Size     uint32
	PreNode  uint32
	NextNode uint32
	Type     uint8
}

//符号表信息
type SymbolItem struct {
	Next   uint32
	Length uint32
}

const MEM_HEADER_SIZE = 13   //内存块大小(4) + 前一个内存块地址(4) + 后一个内存块地址(4) + 内存类型(1)
const SYMBOL_HEADER_SIZE = 8 //next指针(4) + 符号长度(4)

const INVALID_MEM = 0xFFFFFFFF //无效内存值

const (
	HEADER_NODE = iota //头节点
	SYMBOL_NODE        //符号表节点
)

/******************************************************************
    []byte转NodeHeader型
******************************************************************/
func FormatHeader(b []byte) NodeHeader {
	return NodeHeader{Size: comFunc.BytesToUint32(b[0:4]), PreNode: comFunc.BytesToUint32(b[4:8]), NextNode: comFunc.BytesToUint32(b[8:12]), Type: b[12]}
}

/******************************************************************
    NodeHeader转[]byte型
******************************************************************/
func WriteHeader(nodeHeadr NodeHeader, b []byte) {
	comFunc.Uint32ToBytes(nodeHeadr.Size, b[0:4])
	comFunc.Uint32ToBytes(nodeHeadr.PreNode, b[4:8])
	comFunc.Uint32ToBytes(nodeHeadr.NextNode, b[8:12])
	nodeHeadr.Type = b[12]
}
