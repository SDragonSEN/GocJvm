package class

import (
	"strings"
)

type Class interface {
	ReadClass(classname string) ([]byte, error)
}

var BootstrapClassPath []Class //启动类路径
var ExtensionClassPath []Class //拓展类路径
var UserClassPath []Class      //用户类路径

func InitClassPath(userPath string) {
	paths := strings.Split(userPath, ";")
	for _, path := range paths {
		if strings.HasSuffix(path, ".jar") || strings.HasSuffix(path, ".zip") {
			jar, err := NewJarEntry(path)
			if err != nil {
				UserClassPath = append(UserClassPath, jar)
			}
		} else if strings.HasSuffix(path, "/*") {
			dir, err := NewDirEntry(path, true)
			if err != nil {
				UserClassPath = append(UserClassPath, dir)

			}
		} else {
			dir, err := NewDirEntry(path, false)
			if err != nil {
				UserClassPath = append(UserClassPath, dir)
			}
		}
	}
}
