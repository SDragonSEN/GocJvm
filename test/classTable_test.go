package test

import (
	"testing"

	. "basic/memCtrl"
	. "basic/symbol"
	. "class/classTable"
)

/*********************************************************
测试对象:符号表功能
测试内容:根据地址获取符号
*********************************************************/
func Test_ClassTable_Case1(t *testing.T) {
	memAdr1, _ := PutSymbol([]byte("Object"))
	memAdr2, _ := PutSymbol([]byte("String"))
	class1, err1 := PutClass(memAdr1, []byte{})
	class2, err2 := PutClass(memAdr2, []byte{})
	if err1 != nil || err2 != nil {
		t.Error("put失败")
	}
	if class1 == INVALID_MEM || class2 == INVALID_MEM {
		t.Error("分配失败!")
	}
	if class1 != GetClassMemAddr(memAdr1) || class2 != GetClassMemAddr(memAdr2) {
		t.Error("获取出错1!", class1, GetClassMemAddr(memAdr1), class2, GetClassMemAddr(memAdr2))
	}
	if class1 == class2 {
		t.Error("获取出错2!")
	}
}
