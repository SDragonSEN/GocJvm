package classAnaly

import (
	"../memoryControl"
)

/******************************************************************
    功能:获取超类
	入参:无
    返回值:1、超类

	注:不存在超类则返回nil
******************************************************************/
func (self *CLASS_INFO) GetSuperClass() *CLASS_INFO {
	if self.SuperClassAddr == memCtrl.INVALID_MEM {
		return nil
	}
	return (*CLASS_INFO)(memCtrl.GetPointer(self.SuperClassAddr, CLASS_INFO_SIZE))
}

/******************************************************************
    功能:查找函数
	入参:无
    返回值:1、方法结构体指针（不存在方法则返回nil）
          2、code地址（不存在方法则返回INVALID_MEM）
******************************************************************/
func (self *CLASS_INFO) FindMethod(methodName, descriptor uint32) (*METHOD, uint32) {
	methods := *(*[]METHOD)(memCtrl.GetArrayPointer(self.MethodDev+self.LocalAdr, self.MethodNum*METHOD_SIZE, METHOD_SIZE))
	for _, method := range methods {
		if method.MethodName == methodName &&
			method.Descriptor == descriptor {
			if method.CodeAddr != memCtrl.INVALID_MEM {
				return &method, self.LocalAdr + method.CodeAddr
			}
			return &method, memCtrl.INVALID_MEM
		}
	}
	return nil, memCtrl.INVALID_MEM
}
