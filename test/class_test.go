package test

import (
	"testing"

	"../class"
	"../classAnaly"
	"../memoryControl"
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

/*********************************************************
测试对象:Object
测试内容:从classpath中读取Object类
*********************************************************/
func Test_ReadClass_Case1(t *testing.T) {
	class.InitClassPath("")

	data, err := class.ReadClass("java.lang.Object")

	if data == nil {
		t.Error("读取失败，内容为空")
	}

	if err != nil {
		t.Error("读取失败，err不为空")
	}
}

/*********************************************************
测试对象:类加载
测试内容:加载Object类
*********************************************************/
func Test_AnalyClass_Case1(t *testing.T) {
	class.InitClassPath("")
	memCtrl.Init(25535)

	_, err := classAnaly.LoadClass("java.lang.Object")

	if err != nil {
		t.Error("解析失败，err不为空")
	}
}
