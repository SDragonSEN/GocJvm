package test

import (
	"testing"
	"unsafe"

	. "access/access"
	. "access/array"
	. "access/string"
	. "class/classInterface"
	. "methodStack"
)

/*********************************************************
测试对象:DirEntry
测试内容:读取存在的.class文件成功
*********************************************************/
func Test_StructSize(t *testing.T) {
	var accessInfo ACCESS_INFO
	if unsafe.Sizeof(accessInfo) != ACCESS_INFO_SIZE {
		t.Error("ACCESS_INFO_SIZE")
	}
	var arrayInfo ARRAY_INFO
	if unsafe.Sizeof(arrayInfo) != ARRAY_INFO_SIZE {
		t.Error("ARRAY_INFO_SIZE")
	}
	var stringInfo STRING
	if unsafe.Sizeof(stringInfo) != STRING_SIZE {
		t.Error("STRING_SIZE")
	}
	var hashItm HASH_ITEM
	if unsafe.Sizeof(hashItm) != HASH_ITEM_SIZE {
		t.Error("HASH_ITEM_SIZE")
	}
	var classInfo CLASS_INFO
	if unsafe.Sizeof(classInfo) != CLASS_INFO_SIZE {
		t.Error("CLASS_INFO_SIZE")
	}
	var filedItemInfo FILED_ITEM
	if unsafe.Sizeof(filedItemInfo) != FILED_ITEM_SIZE {
		t.Error("FILED_ITEM_SIZE")
	}
	var filedInfo FILED_INFO
	if unsafe.Sizeof(filedInfo) != FILED_INFO_SIZE {
		t.Error("FILED_INFO_SIZE")
	}
	var attriInfo ATTRI_INFO
	if unsafe.Sizeof(attriInfo) != ATTRI_INFO_SIZE {
		t.Error("ATTRI_INFO_SIZE")
	}
	var methodInfo METHOD
	if unsafe.Sizeof(methodInfo) != METHOD_SIZE {
		t.Error("METHOD_SIZE")
	}
	var condeAttri CODE_ATTRI
	if unsafe.Sizeof(condeAttri) != CODE_ATTRI_SIZE {
		t.Error("CODE_ATTRI_SIZE")
	}
	var methodStack METHOD_STACK
	if unsafe.Sizeof(methodStack) != METHOD_STACK_SIZE {
		t.Error("METHOD_STACK_SIZE")
	}
	var methodFrame METHOD_FRAME
	if unsafe.Sizeof(methodFrame) != METHOD_FRAME_SIZE {
		t.Error("methodFrame")
	}
}
