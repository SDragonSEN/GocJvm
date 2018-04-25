package access

import (
	"errors"

	. "basic/memCtrl"
)

type ACCESS_INFO struct {
	TypeAddr uint32
	DataAddr uint32
	NextAddr uint32
}

const ACCESS_INFO_SIZE = 3 * 4

var AccHeaderAdr uint32 = INVALID_MEM

/******************************************************************
    功能:新建引用
	入参:无
    返回值:1、引用指针
	      2、地址
		  3、error
******************************************************************/
func NewAccessInfo() (*ACCESS_INFO, uint32, error) {
	if AccHeaderAdr == INVALID_MEM {
		//初始化类表头结点
		AccHeaderAdr, _ = Malloc(ACCESS_INFO_SIZE, ACCESS_NODE)
		accHeader := (*ACCESS_INFO)(GetPointer(AccHeaderAdr, ACCESS_INFO_SIZE))
		accHeader.DataAddr = INVALID_MEM
		accHeader.NextAddr = INVALID_MEM
		accHeader.TypeAddr = INVALID_MEM

		return accHeader, AccHeaderAdr, nil
	}
	curAddr := AccHeaderAdr
	var curAccess *ACCESS_INFO
	for curAddr != INVALID_MEM {
		curAccess = (*ACCESS_INFO)(GetPointer(curAddr, ACCESS_INFO_SIZE))

		curAddr = curAccess.NextAddr
	}
	newAccAdr, err := Malloc(ACCESS_INFO_SIZE, ACCESS_NODE)
	if err != nil {
		return nil, INVALID_MEM, errors.New("NewAccessInfo():内存不足")
	}
	curAccess.NextAddr = newAccAdr
	newAcc := (*ACCESS_INFO)(GetPointer(newAccAdr, ACCESS_INFO_SIZE))
	newAcc.NextAddr = INVALID_MEM
	return newAcc, newAccAdr, nil
}

/******************************************************************
    功能:获取Access的数据
	入参:access地址
    返回值:1、数据切片
******************************************************************/
func GetData(accAdr uint32) []byte {
	acc := (*ACCESS_INFO)(GetPointer(accAdr, ACCESS_INFO_SIZE))
	return Memory[acc.DataAddr:]
}

/******************************************************************
    功能:获取Access的类
	入参:access地址
    返回值:1、数据切片
******************************************************************/
func GetClassInfo(accAdr uint32) uint32 {
	acc := (*ACCESS_INFO)(GetPointer(accAdr, ACCESS_INFO_SIZE))
	return acc.TypeAddr
}

/******************************************************************
    功能:将引用类型的Type地址修改
	入参:1、源Type地址
	    2、目标Type地址
    返回值:无
******************************************************************/
func ModifyTypeAddr(src, dest uint32) {
	curAddr := AccHeaderAdr
	var curAccess *ACCESS_INFO
	for curAddr != INVALID_MEM {
		curAccess = (*ACCESS_INFO)(GetPointer(curAddr, ACCESS_INFO_SIZE))
		if curAccess.TypeAddr == src {
			curAccess.TypeAddr = dest
		}
		curAddr = curAccess.NextAddr
	}
}
