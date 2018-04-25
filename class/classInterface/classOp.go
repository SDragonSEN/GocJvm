package classInterface

import (
	. "basic/memCtrl"
	. "basic/symbol"
	"fmt"
)

/******************************************************************
    功能:获取超类
	入参:无
    返回值:1、超类

	注:不存在超类则返回nil
******************************************************************/
func (self *CLASS_INFO) GetSuperClass() *CLASS_INFO {
	if self.SuperClassAddr == INVALID_MEM {
		return nil
	}
	return (*CLASS_INFO)(GetPointer(self.SuperClassAddr, CLASS_INFO_SIZE))
}

/******************************************************************
    功能:查找函数
	入参:无
    返回值:1、方法结构体指针（不存在方法则返回nil）
          2、code地址（不存在方法则返回INVALID_MEM）
******************************************************************/
func (self *CLASS_INFO) FindMethod(methodName, descriptor uint32) (*METHOD, uint32) {
	methods := *(*[]METHOD)(GetArrayPointer(self.MethodDev+self.LocalAdr, self.MethodNum*METHOD_SIZE, METHOD_SIZE))
	for _, method := range methods {
		if method.MethodName == methodName &&
			method.Descriptor == descriptor {
			if method.CodeAddr != INVALID_MEM {
				return &method, self.LocalAdr + method.CodeAddr
			}
			return &method, INVALID_MEM
		}
	}
	return nil, INVALID_MEM
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
		methods := *(*[]METHOD)(GetArrayPointer(classInfo.MethodDev+classInfo.LocalAdr, classInfo.MethodNum*METHOD_SIZE, METHOD_SIZE))
		for _, method := range methods {
			if method.MethodName == methodName &&
				method.Descriptor == descriptor {
				if method.CodeAddr != INVALID_MEM {
					return &method, classInfo, self.LocalAdr + method.CodeAddr
				}
				return &method, classInfo, INVALID_MEM
			}
		}
		if self.SuperClassAddr != INVALID_MEM {
			classInfo = (*CLASS_INFO)(GetPointer(classInfo.SuperClassAddr, CLASS_INFO_SIZE))
		} else {
			break
		}
	}
	return nil, nil, INVALID_MEM
}

/******************************************************************
    功能:获取static字段值
	入参:无
    返回值:1、值
******************************************************************/
func (self *CLASS_INFO) GetStaticData32(filedName, descriptor uint32) uint32 {
	//fmt.Println(string(GetSymbol(filedName)))
	fileds := *(*[]FILED_ITEM)(GetArrayPointer(self.StaticParaDev+self.LocalAdr, self.StaticParaNum*FILED_ITEM_SIZE, FILED_ITEM_SIZE))
	for _, filed := range fileds {
		if filed.FiledName == filedName {
			staticAdr := self.StaticMem
			return *(*uint32)(GetPointer(staticAdr+filed.Index*4, 4))
		}
	}
	panic("GetStaticData32")
	return 0
}

/******************************************************************
    功能:保存static字段值
	入参:无
    返回值:1、值
******************************************************************/
func (self *CLASS_INFO) PutStaticData32(filedName, descriptor, v uint32) {
	//fmt.Println(string(GetSymbol(filedName)))
	fileds := *(*[]FILED_ITEM)(GetArrayPointer(self.StaticParaDev+self.LocalAdr, self.StaticParaNum*FILED_ITEM_SIZE, FILED_ITEM_SIZE))
	for _, filed := range fileds {
		if filed.FiledName == filedName {
			staticAdr := self.StaticMem
			p := (*uint32)(GetPointer(staticAdr+filed.Index*4, 4))
			*p = v
			return
		}
	}
	panic("PutStaticData32")
}

/******************************************************************
    功能:获取static字段值(long,double)
	入参:无
    返回值:1、值
******************************************************************/
func (self *CLASS_INFO) GetStaticData64(filedName, descriptor uint32) [2]uint32 {
	var v [2]uint32
	fileds := *(*[]FILED_ITEM)(GetArrayPointer(self.StaticParaDev+self.LocalAdr, self.UnstaticParaDev-self.StaticParaDev, FILED_ITEM_SIZE))
	for _, filed := range fileds {
		if filed.FiledName == filedName {
			staticAdr := self.StaticMem
			v[0] = *(*uint32)(GetPointer(staticAdr+filed.Index*4, 4))
			v[1] = *(*uint32)(GetPointer(staticAdr+filed.Index*4+4, 4))
			return v
		}
	}
	panic("GetStaticData64")
	return v
}

/******************************************************************
    功能:保存static字段值(long,double)
	入参:无
    返回值:1、值
******************************************************************/
func (self *CLASS_INFO) PutStaticData64(filedName, descriptor, v0, v1 uint32) {
	fileds := *(*[]FILED_ITEM)(GetArrayPointer(self.StaticParaDev+self.LocalAdr, self.UnstaticParaDev-self.StaticParaDev, FILED_ITEM_SIZE))
	for _, filed := range fileds {
		if filed.FiledName == filedName {
			staticAdr := self.StaticMem
			p0 := (*uint32)(GetPointer(staticAdr+filed.Index*4, 4))
			p1 := (*uint32)(GetPointer(staticAdr+filed.Index*4+4, 4))
			*p0 = v0
			*p1 = v1
			return
		}
	}
	panic("PutStaticData64")
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
		fileds := *(*[]FILED_ITEM)(GetArrayPointer(classInfo.UnstaticParaDev+classInfo.LocalAdr, classInfo.UnstaticParaNum*FILED_ITEM_SIZE, FILED_ITEM_SIZE))
		for _, filed := range fileds {
			if filed.FiledName == filedName {
				return index + filed.Index
			}
		}
		index += classInfo.UnstaticParaSize
		if classInfo.SuperClassAddr == INVALID_MEM {
			fmt.Println("GetUnstaticDataIndex:", self.LocalAdr, string(GetSymbol(filedName)), string(GetSymbol(descriptor)))
			panic("GetUnstaticDataIndex")
		}
		classInfo = (*CLASS_INFO)(GetPointer(classInfo.SuperClassAddr, CLASS_INFO_SIZE))
	}
	panic("GetUnstaticDataIndex")
	return 0
}
