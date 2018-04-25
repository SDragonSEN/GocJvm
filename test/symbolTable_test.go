package test

import (
	"bytes"
	"testing"

	. "basic/memCtrl"
	. "basic/symbol"
)

/*********************************************************
测试对象:符号表功能
测试内容:相同符号存储在同一位置
*********************************************************/
func Test_SymbolTable_Case1(t *testing.T) {
	memAdr1, err1 := PutSymbol([]byte("Object"))
	memAdr2, err2 := PutSymbol([]byte("Object"))
	if err1 != nil || err2 != nil {
		t.Error("分配失败1!")
	}
	if memAdr1 == INVALID_MEM || memAdr2 == INVALID_MEM {
		t.Error("分配失败2!")
	}
	if memAdr1 != memAdr2 {
		t.Error("两次返回的地址不一致:", memAdr1, " ", memAdr2)
	}
}

/*********************************************************
测试对象:符号表功能
测试内容:根据地址获取符号
*********************************************************/
func Test_SymbolTable_Case2(t *testing.T) {
	memAdr1, _ := PutSymbol([]byte("Object"))
	sym := GetSymbol(memAdr1)
	if bytes.Compare([]byte("Object"), sym) != 0 {
		t.Error("获取符号内容错误")
	}
}
