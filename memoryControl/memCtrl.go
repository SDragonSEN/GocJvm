package memCtrl

import (
	"errors"
	"fmt"
	"os"

	"../comFunc"
)

var Memory []byte
var memSize uint32
var symHeaderAdr uint32
var classHeaderAdr uint32

/******************************************************************
    内存空间初始化，参数为 1、内存大小；
******************************************************************/
func Init(size uint32) {

	if size < MEM_HEADER_SIZE {
		fmt.Println("memCtrl:InitEx():参数错误,size:", size)
		os.Exit(-1)
	}
	//分配内存
	Memory = make([]byte, size)
	memSize = size

	//初始化Constant内存的头结点
	headerNode, _ := getHeader()
	headerNode.NextNode = INVALID_MEM
	headerNode.PreNode = INVALID_MEM
	headerNode.Size = MEM_HEADER_SIZE
	headerNode.Type = HEADER_NODE

	//初始化符号表头结点
	symHeaderAdr, _ = Malloc(SYMBOL_HEADER_SIZE, SYMBOL_NODE)
	symHeader := (*SymbolItem)(comFunc.BytesToUnsafePointer(Memory[symHeaderAdr:]))
	symHeader.Length = 0
	symHeader.Next = INVALID_MEM

	//初始化类表头结点
	classHeaderAdr, _ = Malloc(CLASS_HEADER_SIZE, 0)
	classHeader := (*ClassItem)(comFunc.BytesToUnsafePointer(Memory[classHeaderAdr:]))
	classHeader.ClassName = INVALID_MEM
	classHeader.Next = INVALID_MEM

}

/******************************************************************
    获取HeaderNode
******************************************************************/
func getHeader() (*NodeHeader, uint32) {
	return (*NodeHeader)(comFunc.BytesToUnsafePointer(Memory[0:MEM_HEADER_SIZE])), 0
}

/******************************************************************
    分配内存(指定类型)，返回地址
******************************************************************/
func Malloc(size uint32, memType uint8) (uint32, error) {
	var header *NodeHeader
	var addr uint32

	header, addr = getHeader()

	for header.NextNode != INVALID_MEM {
		/* 两个节点之间是否有足够大小 */
		if (header.NextNode - header.Size - addr) >= (size + MEM_HEADER_SIZE) {
			/* 分配新的节点,并初始化Header信息 */
			newAddr := addr + header.Size
			newNode := (*NodeHeader)(comFunc.BytesToUnsafePointer(Memory[newAddr : newAddr+MEM_HEADER_SIZE]))
			newNode.PreNode = addr
			newNode.NextNode = header.NextNode
			newNode.Size = size + MEM_HEADER_SIZE
			newNode.Type = memType

			/* 将分配的内存刷成全0 */
			for i := newAddr + MEM_HEADER_SIZE; i < newAddr+MEM_HEADER_SIZE+size; i++ {
				Memory[i] = 0
			}
			/* 修改下一个节点前指针 */
			nextNode := (*NodeHeader)(comFunc.BytesToUnsafePointer(Memory[header.NextNode : header.NextNode+MEM_HEADER_SIZE]))
			nextNode.PreNode = newAddr

			/* 修改上一个节点后指针 */
			header.NextNode = newAddr

			return newAddr, nil
		}
		/* 指向下一个节点 */
		addr = header.NextNode
		header = (*NodeHeader)(comFunc.BytesToUnsafePointer(Memory[header.NextNode : header.NextNode+MEM_HEADER_SIZE]))
	}
	/* 最后一个节点之后有没有足够的内存 */
	if (memSize - header.Size - addr) >= (size + MEM_HEADER_SIZE) {
		newAddr := addr + header.Size
		newNode := (*NodeHeader)(comFunc.BytesToUnsafePointer(Memory[newAddr : newAddr+MEM_HEADER_SIZE]))
		newNode.PreNode = addr
		newNode.NextNode = INVALID_MEM
		newNode.Size = size + MEM_HEADER_SIZE
		newNode.Type = memType

		for i := newAddr + MEM_HEADER_SIZE; i < newAddr+MEM_HEADER_SIZE+size; i++ {
			Memory[i] = 0
		}

		header.NextNode = newAddr

		return newAddr + MEM_HEADER_SIZE, nil
	}

	return 0, errors.New("Malloc():No Enough Memory!")
}

/******************************************************************
    释放内存
******************************************************************/
func MemFree(addr int) error {
	if addr == 0 || addr == MEM_HEADER_SIZE {
		return errors.New("MemFree():Can't Free HeaderNode!")
	}
	deleteNode := (*NodeHeader)(comFunc.BytesToUnsafePointer(Memory[addr-MEM_HEADER_SIZE : addr]))

	/* 修改下一个节点前指针 */
	if deleteNode.NextNode != INVALID_MEM {
		nextNode := (*NodeHeader)(comFunc.BytesToUnsafePointer(Memory[deleteNode.NextNode : deleteNode.NextNode+MEM_HEADER_SIZE]))
		nextNode.PreNode = deleteNode.PreNode

	}
	/* 修改上一个节点后指针 */
	preNode := (*NodeHeader)(comFunc.BytesToUnsafePointer(Memory[deleteNode.PreNode : deleteNode.PreNode+MEM_HEADER_SIZE]))
	preNode.NextNode = deleteNode.NextNode

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
	fmt.Println(Memory)
}
