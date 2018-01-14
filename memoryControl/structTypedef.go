package memCtrl

type NodeHeader struct {
	Size     uint32
	PreNode  uint32
	NextNode uint32
	Type     uint8
}

const MEM_HEADER_SIZE = 13     //内存块大小(4) + 前一个内存块地址(4) + 后一个内存块地址(4) + 内存类型(1)
const INVALID_MEM = 0xFFFFFFFF //无效内存值

const (
	HEADER_NODE = iota
)

const (
	CONSTANT = iota
	RECYCLE
)
