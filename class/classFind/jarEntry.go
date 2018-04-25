package classFind

import (
	"archive/zip"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type JarEntry struct {
	Path string
}

func NewJarEntry(jarPath string) (*JarEntry, error) {
	absPath, err := filepath.Abs(jarPath)
	if err != nil {
		return nil, err
	}
	f, err := os.Stat(absPath)

	if err != nil {
		return nil, err
	}
	if f.IsDir() {
		return nil, errors.New("不是文件！")
	}
	return &JarEntry{absPath}, nil
}

/******************************************************************
    根据类名(完全限定名,包之间用 . 分隔)从.jar或.zip中查找.class文件。
******************************************************************/
func (this JarEntry) ReadClass(classname string) ([]byte, error) {
	classpath := strings.Replace(classname, ".", "/", -1)
	classpath += ".class"
	jar, err := zip.OpenReader(this.Path)

	if err != nil {
		return nil, err
	}
	defer jar.Close()
	for _, file := range jar.File {

		if file.Name == classpath {

			rc, err := file.Open()
			if err != nil {
				return nil, err
			}

			defer rc.Close()
			data, err := ioutil.ReadAll(rc)
			if err != nil {
				return nil, err
			}
			return data, nil
		}
	}

	return nil, errors.New("Class未找到")
}
