package class

import (
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
}

//指定用户Class路径,即-cp后的路径
func initUserClassPath(userPath string) {
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
	classpath := os.Getenv("classpath")
	if classpath != nil {
		dir, err := NewDirEntry(classpath, false)
		if err == nil {
			BootstrapClassPath = append(BootstrapClassPath, dir)
		}
		dir, err = NewDirEntry(classpath+"/lib/ext", false)

		ExtensionClassPath = append(ExtensionClassPath, dir)
	}
}
