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

//类内存节点
type ClassItem struct {
	Next      uint32
	ClassName uint32 //对应符号表中的地址
}

const MEM_HEADER_SIZE = 13   //内存块大小(4) + 前一个内存块地址(4) + 后一个内存块地址(4) + 内存类型(1)
const SYMBOL_HEADER_SIZE = 8 //next指针(4) + 符号长度(4)
const CLASS_HEADER_SIZE = 8  //next指针(4) + 类名(4)

const INVALID_MEM = 0xFFFFFFFF //无效内存值

const (
	HEADER_NODE = iota //头节点
	SYMBOL_NODE        //符号表节点
	CLASS_NODE         //类表节点
	ACCESS_NODE        //引用类型节点
	ARRAY_NODE         //数组节点
	DATA_NODE          //数据节点
)
