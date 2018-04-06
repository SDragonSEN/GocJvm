package classAnaly

import (
	"../memoryControl"
)

func (self *CLASS_INFO) GetSuperClass() *CLASS_INFO {
	if self.SuperClassAddr == memCtrl.INVALID_MEM {
		return nil
	}
	return (*CLASS_INFO)(memCtrl.GetPointer(self.SuperClassAddr, CLASS_INFO_SIZE))
}
