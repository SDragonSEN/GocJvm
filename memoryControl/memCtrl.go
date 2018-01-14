package memCtrl

import (
	"errors"
	"fmt"
	"os"
)

var memory []byte
var recycleMemIndex uint32 //可回收空间起始地址
var memSize uint32

/******************************************************************
    内存空间初始化，参数为内存大小
******************************************************************/
func Init(size uint32) {
	InitEx(size, size/2)
}

/******************************************************************
    内存空间初始化，参数为 1、内存大小；2、可回收空间起始地址
******************************************************************/
func InitEx(size uint32, index uint32) {

	if index > size || size < MEM_HEADER_SIZE*2 || size-index+1 < MEM_HEADER_SIZE || index < MEM_HEADER_SIZE {
		fmt.Println("memCtrl:InitEx():参数错误,size,index:", size, index)
		os.Exit(-1)
	}
	//分配内存
	memory = make([]byte, size)
	//设置可回收内存起始地址
	recycleMemIndex = index
	memSize = size
	//初始化Constant内存的头结点
	constantHeader, addr := getConstantHeader()
	constantHeader.NextNode = INVALID_MEM
	constantHeader.PreNode = INVALID_MEM
	constantHeader.Size = MEM_HEADER_SIZE
	constantHeader.Type = HEADER_NODE
	WriteHeader(constantHeader, memory[addr:addr+MEM_HEADER_SIZE])
	//初始化Recycle内存的头结点
	recycleHeader, addr := getRecycleHeader()
	recycleHeader.NextNode = INVALID_MEM
	recycleHeader.PreNode = INVALID_MEM
	recycleHeader.Size = MEM_HEADER_SIZE
	recycleHeader.Type = HEADER_NODE
	WriteHeader(recycleHeader, memory[addr:addr+MEM_HEADER_SIZE])
}

/******************************************************************
    获取不回收地址的HeaderNode
******************************************************************/
func getConstantHeader() (NodeHeader, uint32) {
	return FormatHeader(memory[0:MEM_HEADER_SIZE]), 0
}

/******************************************************************
    获取可回收地址的HeaderNode
******************************************************************/
func getRecycleHeader() (NodeHeader, uint32) {
	return FormatHeader(memory[recycleMemIndex : recycleMemIndex+MEM_HEADER_SIZE]), recycleMemIndex
}

/******************************************************************
    分配内存，返回地址
******************************************************************/
func Malloc(size uint32, memType uint8) (uint32, error) {
	return MallocEX(size, memType, RECYCLE)
}

/******************************************************************
    分配内存(指定类型)，返回地址
******************************************************************/
func MallocEX(size uint32, memType uint8, memAttr uint8) (uint32, error) {
	var header NodeHeader
	var addr uint32
	var end uint32
	if memAttr == CONSTANT {
		header, addr = getConstantHeader()
		end = recycleMemIndex
	} else {
		header, addr = getRecycleHeader()
		end = memSize
	}

	for header.NextNode != INVALID_MEM {
		/* 两个节点之间是否有足够大小 */
		if (header.NextNode - header.Size - addr) >= (size + MEM_HEADER_SIZE) {
			/* 分配新的节点,并初始化Header信息 */
			newAddr := addr + header.Size
			newNode := FormatHeader(memory[newAddr : newAddr+MEM_HEADER_SIZE])
			newNode.PreNode = addr
			newNode.NextNode = header.NextNode
			newNode.Size = size + MEM_HEADER_SIZE
			newNode.Type = memType
			WriteHeader(newNode, memory[newAddr:newAddr+MEM_HEADER_SIZE])
			/* 将分配的内存刷成全0 */
			for i := newAddr + MEM_HEADER_SIZE; i < newAddr+MEM_HEADER_SIZE+size; i++ {
				memory[i] = 0
			}
			/* 修改下一个节点前指针 */
			nextNode := FormatHeader(memory[header.NextNode : header.NextNode+MEM_HEADER_SIZE])
			nextNode.PreNode = newAddr
			WriteHeader(nextNode, memory[header.NextNode:header.NextNode+MEM_HEADER_SIZE])
			/* 修改上一个节点后指针 */
			header.NextNode = newAddr
			WriteHeader(header, memory[addr:addr+MEM_HEADER_SIZE])
			return newAddr, nil
		}
		/* 指向下一个节点 */
		addr = header.NextNode
		header = FormatHeader(memory[header.NextNode : header.NextNode+MEM_HEADER_SIZE])
	}
	/* 最后一个节点之后有没有足够的内存 */
	if (end - header.Size - addr) >= (size + MEM_HEADER_SIZE) {
		newAddr := addr + header.Size
		newNode := FormatHeader(memory[newAddr : newAddr+MEM_HEADER_SIZE])
		newNode.PreNode = addr
		newNode.NextNode = INVALID_MEM
		newNode.Size = size + MEM_HEADER_SIZE
		newNode.Type = memType
		WriteHeader(newNode, memory[newAddr:newAddr+MEM_HEADER_SIZE])

		for i := newAddr + MEM_HEADER_SIZE; i < newAddr+MEM_HEADER_SIZE+size; i++ {
			memory[i] = 0
		}

		header.NextNode = newAddr
		WriteHeader(header, memory[addr:addr+MEM_HEADER_SIZE])
		return newAddr, nil
	}

	return 0, errors.New("No Enough Memory!")
}

/******************************************************************
    释放内存
******************************************************************/
func MemFree(addr int) error {
	return nil
}

/******************************************************************
    重新分配内存大小
******************************************************************/
func ReAlloc(size int) (int, error) {
	return 0, nil
}

/******************************************************************
    Log内存信息
******************************************************************/
func LogMem() {
	fmt.Println("可回收内存起始索引:", recycleMemIndex)
	fmt.Println(memory)
}
