package test

import (
	"testing"
	"unsafe"

	"accessOp"
	"classAnaly"
	"methodStack"
)

/*********************************************************
测试对象:DirEntry
测试内容:读取存在的.class文件成功
*********************************************************/
func Test_StructSize(t *testing.T) {
	var accessInfo access.ACCESS_INFO
	if unsafe.Sizeof(accessInfo) != access.ACCESS_INFO_SIZE {
		t.Error("ACCESS_INFO_SIZE")
	}
	var arrayInfo access.ARRAY_INFO
	if unsafe.Sizeof(arrayInfo) != access.ARRAY_INFO_SIZE {
		t.Error("ARRAY_INFO_SIZE")
	}
	var stringInfo access.STRING
	if unsafe.Sizeof(stringInfo) != access.STRING_SIZE {
		t.Error("STRING_SIZE")
	}
	var hashItm access.HASH_ITEM
	if unsafe.Sizeof(hashItm) != access.HASH_ITEM_SIZE {
		t.Error("HASH_ITEM_SIZE")
	}
	var classInfo classAnaly.CLASS_INFO
	if unsafe.Sizeof(classInfo) != classAnaly.CLASS_INFO_SIZE {
		t.Error("CLASS_INFO_SIZE")
	}
	var filedItemInfo classAnaly.FILED_ITEM
	if unsafe.Sizeof(filedItemInfo) != classAnaly.FILED_ITEM_SIZE {
		t.Error("FILED_ITEM_SIZE")
	}
	var filedInfo classAnaly.FILED_INFO
	if unsafe.Sizeof(filedInfo) != classAnaly.FILED_INFO_SIZE {
		t.Error("FILED_INFO_SIZE")
	}
	var attriInfo classAnaly.ATTRI_INFO
	if unsafe.Sizeof(attriInfo) != classAnaly.ATTRI_INFO_SIZE {
		t.Error("ATTRI_INFO_SIZE")
	}
	var methodInfo classAnaly.METHOD
	if unsafe.Sizeof(methodInfo) != classAnaly.METHOD_SIZE {
		t.Error("METHOD_SIZE")
	}
	var condeAttri classAnaly.CODE_ATTRI
	if unsafe.Sizeof(condeAttri) != classAnaly.CODE_ATTRI_SIZE {
		t.Error("CODE_ATTRI_SIZE")
	}
	var methodStack method.METHOD_STACK
	if unsafe.Sizeof(methodStack) != method.METHOD_STACK_SIZE {
		t.Error("METHOD_STACK_SIZE")
	}
	var methodFrame method.METHOD_FRAME
	if unsafe.Sizeof(methodFrame) != method.METHOD_FRAME_SIZE {
		t.Error("methodFrame")
	}
}
