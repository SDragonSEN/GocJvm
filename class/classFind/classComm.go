package classFind

import (
	"errors"

	"os"
	"strings"
)

type Class interface {
	ReadClass(classname string) ([]byte, error)
}

var BootstrapClassPath []Class //启动类路径
var ExtensionClassPath []Class //拓展类路径
var UserClassPath []Class      //用户类路径

func InitClassPath(userPath string) {
	initUserClassPath(userPath)
	dir, err := NewDirEntry("./", false)
	if err == nil {
		UserClassPath = append(UserClassPath, dir)
	} else {
		panic("InitClassPath()")
	}
	initBootstrapClassPath()
}

/******************************************************************
    根据类名(完全限定名,包之间用 . 分隔)读取.class文件。
	顺序:Bootstrap Extension User
******************************************************************/
func ReadClass(classname string) ([]byte, error) {
	for _, boot := range BootstrapClassPath {
		context, err := boot.ReadClass(classname)
		if nil == err {
			return context, err
		}
	}
	for _, extension := range ExtensionClassPath {
		context, err := extension.ReadClass(classname)
		if nil == err {
			return context, err
		}
	}
	for _, user := range UserClassPath {
		context, err := user.ReadClass(classname)
		if nil == err {
			return context, err
		}
	}
	return nil, errors.New("class not found")
}

//指定用户Class路径,即-cp后的路径
func initUserClassPath(userPath string) {
	UserClassPath = make([]Class, 0)
	paths := strings.Split(userPath, ";")
	for _, path := range paths {
		if strings.HasSuffix(path, ".jar") || strings.HasSuffix(path, ".zip") {
			jar, err := NewJarEntry(path)
			if err == nil {
				UserClassPath = append(UserClassPath, jar)
			}
		} else if strings.HasSuffix(path, "/*") {
			dir, err := NewDirEntry(path, true)
			if err == nil {
				UserClassPath = append(UserClassPath, dir)

			}
		} else {
			dir, err := NewDirEntry(path, false)
			if err == nil {
				UserClassPath = append(UserClassPath, dir)
			}
		}
	}
}

//初始化启动类和拓展类路径
func initBootstrapClassPath() {
	BootstrapClassPath = make([]Class, 0)
	ExtensionClassPath = make([]Class, 0)

	classpath := os.Getenv("classpath")

	if classpath != "" {
		dir, err := NewDirEntry(classpath+"/lib", true)
		if err == nil {
			BootstrapClassPath = append(BootstrapClassPath, dir)
		}
		dir, err = NewDirEntry(classpath+"/lib/ext", true)
		if err == nil {
			ExtensionClassPath = append(ExtensionClassPath, dir)
		}
	}
}
