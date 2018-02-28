package memCtrl

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
