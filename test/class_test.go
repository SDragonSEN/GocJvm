package test

import (
	"testing"

	"../class"
)

/*********************************************************
测试对象:DirEntry
测试内容:读取存在的.class文件成功
*********************************************************/
func Test_DirEntry_Case1(t *testing.T) {
	dir, err := class.NewDirEntry("./stub/class/DirEntry", false)

	if err != nil {
		t.Error("地址转化错误")
		return
	}

	data, err := dir.ReadClass("package1.Main")
	if data == nil {
		t.Error("读取失败，内容为空")
	}

	if err != nil {
		t.Error("读取失败，err不为空")
	}
}

/*********************************************************
测试对象:DirEntry
测试内容:读取不存在的.class文件失败
*********************************************************/
func Test_DirEntry_Case2(t *testing.T) {
	dir, err := class.NewDirEntry("./stub/class/DirEntry", false)

	if err != nil {
		t.Error("地址转化错误")
		return
	}

	data, err := dir.ReadClass("package1.MainError")
	if data != nil {
		t.Fail()
	}

	if err == nil {
		t.Error()
	}
}

/*********************************************************
测试对象:DirEntry
测试内容:路径不存在
*********************************************************/
func Test_DirEntry_Case3(t *testing.T) {
	_, err := class.NewDirEntry("./stub/class/DirEntry/error", false)

	if err == nil {
		t.Error("地址转化错误")
	}
}

/*********************************************************
测试对象:DirEntry
测试内容:读取存在的.class文件成功
*********************************************************/
func Test_DirEntry_Case4(t *testing.T) {
	dir, err := class.NewDirEntry("./stub/class/JarEntry", true)

	if err != nil {
		t.Error("地址转化错误")
		return
	}

	data, err := dir.ReadClass("package2.Class2")
	if data == nil {
		t.Error("读取失败，内容为空")
	}

	if err != nil {
		t.Error("读取失败，err不为空")
	}
}
