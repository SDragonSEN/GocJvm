package classAnaly

import (
	"memoryControl"
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

/******************************************************************
    功能:查找函数,如果没找到，则去找父类
	入参:无
    返回值:1、方法结构体指针（不存在方法则返回nil）
	      2、Class Info
          3、code地址（不存在方法则返回INVALID_MEM）
******************************************************************/
func (self *CLASS_INFO) FindMethodEx(methodName, descriptor uint32) (*METHOD, *CLASS_INFO, uint32) {
	classInfo := self
	for {
		methods := *(*[]METHOD)(memCtrl.GetArrayPointer(classInfo.MethodDev+classInfo.LocalAdr, classInfo.MethodNum*METHOD_SIZE, METHOD_SIZE))
		for _, method := range methods {
			if method.MethodName == methodName &&
				method.Descriptor == descriptor {
				if method.CodeAddr != memCtrl.INVALID_MEM {
					return &method, classInfo, self.LocalAdr + method.CodeAddr
				}
				return &method, classInfo, memCtrl.INVALID_MEM
			}
		}
		if self.SuperClassAddr != memCtrl.INVALID_MEM {
			classInfo = (*CLASS_INFO)(memCtrl.GetPointer(classInfo.SuperClassAddr, CLASS_INFO_SIZE))
		} else {
			break
		}
	}
	return nil, nil, memCtrl.INVALID_MEM
}

/******************************************************************
    功能:获取static字段值
	入参:无
    返回值:1、值
******************************************************************/
func (self *CLASS_INFO) GetStaticData32(filedName, descriptor uint32) uint32 {
	//fmt.Println(string(memCtrl.GetSymbol(filedName)))
	fileds := *(*[]FILED_ITEM)(memCtrl.GetArrayPointer(self.StaticParaDev+self.LocalAdr, self.StaticParaNum*FILED_ITEM_SIZE, FILED_ITEM_SIZE))
	for _, filed := range fileds {
		if filed.FiledName == filedName {
			staticAdr := self.StaticMem
			return *(*uint32)(memCtrl.GetPointer(staticAdr+filed.Index*4, 4))
		}
	}
	panic("GetStaticData32")
	return 0
}

/******************************************************************
    功能:获取static字段值(long,double)
	入参:无
    返回值:1、值
******************************************************************/
func (self *CLASS_INFO) GetStaticData64(filedName, descriptor uint32) [2]uint32 {
	var v [2]uint32
	fileds := *(*[]FILED_ITEM)(memCtrl.GetArrayPointer(self.StaticParaDev+self.LocalAdr, self.UnstaticParaDev-self.StaticParaDev, FILED_ITEM_SIZE))
	for _, filed := range fileds {
		if filed.FiledName == filedName {
			staticAdr := self.StaticMem
			v[0] = *(*uint32)(memCtrl.GetPointer(staticAdr+filed.Index*4, 4))
			v[1] = *(*uint32)(memCtrl.GetPointer(staticAdr+filed.Index*4+4, 4))
			return v
		}
	}
	panic("GetStaticData64")
	return v
}

/******************************************************************
    功能:获取unstatic字段在实例中的位置
	入参:无
    返回值:1、值
******************************************************************/
func (self *CLASS_INFO) GetUnstaticDataIndex(filedName, descriptor uint32) uint32 {
	var index uint32 = 0
	classInfo := self
	for {
		fileds := *(*[]FILED_ITEM)(memCtrl.GetArrayPointer(classInfo.UnstaticParaDev+classInfo.LocalAdr, classInfo.UnstaticParaNum*FILED_ITEM_SIZE, FILED_ITEM_SIZE))
		for _, filed := range fileds {
			if filed.FiledName == filedName {
				return index + filed.Index
			}
		}
		index += classInfo.UnstaticParaSize
		if classInfo.SuperClassAddr == memCtrl.INVALID_MEM {
			panic("GetUnstaticDataIndex")
		}
		classInfo = (*CLASS_INFO)(memCtrl.GetPointer(classInfo.SuperClassAddr, CLASS_INFO_SIZE))
	}
	panic("GetUnstaticDataIndex")
	return 0
}
