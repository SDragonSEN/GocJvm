package classFind

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type DirEntry struct {
	Path        string  //绝对路径
	OtherEntrys []Class //该路径下的jar或zip包
}

/******************************************************************
    根据路径生成NewDirEntry
******************************************************************/
func NewDirEntry(path string, isContainJar bool) (*DirEntry, error) {
	absPath, err := filepath.Abs(path)

	if err != nil {
		return nil, err
	}
	f, err := os.Stat(absPath)

	if err != nil {
		return nil, err
	}
	if !f.IsDir() {
		return nil, errors.New("不是文件夹！")
	}
	if isContainJar == true {
		jarEntrys, err := getJarEntry(absPath)
		if err == nil && len(jarEntrys) > 0 {
			return &DirEntry{absPath, jarEntrys}, nil
		}
	}
	return &DirEntry{absPath, nil}, nil
}
func getJarEntry(absPath string) ([]Class, error) {
	file_list, err := ioutil.ReadDir(absPath)
	if err != nil {
		return nil, err
	}

	jarEntrys := make([]Class, 0)

	for _, file := range file_list {
		if file.IsDir() {
			continue
		}
		if !isJarFile(file.Name()) {
			continue
		}
		jarEntry, err := NewJarEntry(filepath.Join(absPath, file.Name()))
		if err != nil {
			continue
		}
		jarEntrys = append(jarEntrys, *jarEntry)
	}
	return jarEntrys, nil
}
func isJarFile(filename string) bool {
	isJar := strings.HasSuffix(filename, ".jar")
	isZip := strings.HasSuffix(filename, ".zip")

	if isJar || isZip {
		return true
	}

	return false
}

/******************************************************************
    根据类名(完全限定名,包之间用 . 分隔)读取Dir指定.class文件。如果文件结构中
没有对应的类，则在该路径下的.zip,.jar包内查找。
    注:用/分隔也没问题，只要不把.class后缀名加上就行
******************************************************************/
func (this DirEntry) ReadClass(classname string) ([]byte, error) {
	//根据类名获取文件地址
	folder := strings.Split(classname, ".")
	classpath := this.Path
	for _, value := range folder {
		classpath += "/" + value
	}
	classpath += ".class"

	//读取.class文件
	data, err := ioutil.ReadFile(classpath)

	if err == nil {
		//.class文件存在
		return data, nil
	} else {
		if this.OtherEntrys != nil {
			//.class文件不存在，查找.jar和.zip文件中是否存在该class
			for _, class := range this.OtherEntrys {
				data, err = class.ReadClass(classname)
				if err == nil {
					return data, nil
				}
			}
		}
		return nil, errors.New("class not found")
	}
}
