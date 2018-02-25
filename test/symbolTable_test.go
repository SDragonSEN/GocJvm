package test

import (
	"bytes"
	"testing"

	"../memoryControl"
)

/*********************************************************
测试对象:符号表功能
测试内容:相同符号存储在同一位置
*********************************************************/
func Test_SymbolTable_Case1(t *testing.T) {
	memCtrl.Init(1024)
	memAdr1 := memCtrl.PutSymbol([]byte("Object"))
	memAdr2 := memCtrl.PutSymbol([]byte("Object"))
	if memAdr1 != memAdr2 {
		t.Error("两次返回的地址不一致:", memAdr1, " ", memAdr2)
	}
}

/*********************************************************
测试对象:符号表功能
测试内容:根据地址获取符号
*********************************************************/
func Test_SymbolTable_Case2(t *testing.T) {
	memCtrl.Init(1024)
	memAdr1 := memCtrl.PutSymbol([]byte("Object"))
	sym := memCtrl.GetSymbol(memAdr1)
	if bytes.Compare([]byte("Object"), sym) != 0 {
		t.Error("获取符号内容错误")
	}
}