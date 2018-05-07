package classInterface

type CONSTANT_TYPE_16 struct {
	Param1 uint16
	Param2 uint16
}

type CONSTANT_TYPE_32 struct {
	Param uint32
}

type CLASS_INFO struct {
	ClassName             uint32 //类名
	SuperClassAddr        uint32 //父类地址,为0代表是Object类
	AccessFlag            uint16 //可访问属性
	IsCInit               bool
	rsv                   uint8
	ConstNum              uint32 //常量数量
	ClassConstDev         uint32 //类常量偏移
	FiledInfoDev          uint32 //参数信息偏移
	UnstaticParaDev       uint32 //非static参数地址
	UnstaticParaSize      uint32 //非static参数大小
	UnstaticParaNum       uint32 //非static参数个数
	UnstaticParaTotalSize uint32 //非static参数内存总大小(即，分配实例的大小)
	StaticParaDev         uint32 //static参数地址
	StaticParaSize        uint32 //static参数大小
	StaticParaNum         uint32 //static参数个数
	StaticMem             uint32 //类实例地址
	InterfaceDev          uint32 //接口定义偏移
	InterfaceNum          uint32 //接口数量
	MethodDev             uint32 //方法定义偏移
	MethodNum             uint32 //方法数量
	LocalAdr              uint32 //该类的地址
}

const CLASS_INFO_SIZE = 19 * 4

type FILED_ITEM struct {
	FiledName    uint32 //字段名(符号表索引)
	Index        uint32 //实例(包括类实例)中的索引值,从0开始，遇到long和double则跳1
	FiledInfoDev uint32 //字段描述偏移
}

const FILED_ITEM_SIZE = 3 * 4

type FILED_INFO struct {
	AccessFlag uint16 //可访问性
	rsv        [2]uint8
	Descriptor uint32 //描述符(符号表索引)
	AttriCount uint32 //属性数量
}

const FILED_INFO_SIZE = 3 * 4

type ATTRI_INFO struct {
	AttriName uint32 //属性名(符号表中的地址)
	Length    uint32 //长度
}

const ATTRI_INFO_SIZE = 8

type CONST_PAIR struct {
	StaticFiledIndex uint32 //static字段索引
	ConstIndex       uint32 //常量索引
	IsLongOrDouble   bool
}

type METHOD struct {
	AccessFlag uint16 //可访问属性
	rsv        [2]uint8
	MethodName uint32 //方法名
	Descriptor uint32 //描述符
	CodeAddr   uint32 //code地址,code属性里没有属性和长度，直接就是Code结构体开始
	Attribute  uint32 //属性地址
	AttriNum   uint32 //属性数量
}

const METHOD_SIZE = 6 * 4

type CODE_ATTRI struct {
	MaxStack       uint32 //方法栈
	MaxLocal       uint32 //局部变量大小
	CodeLength     uint32
	ExceptionCount uint32
	AttriNum       uint32
}

const CODE_ATTRI_SIZE = 20

var MagicNum = []byte{0xCA, 0xFE, 0xBA, 0xBE}

const FILED_ACC_PUBLIC = 0x0001
const FILED_ACC_PRIVATE = 0x0002
const FILED_ACC_PROTECTED = 0x0004
const FILED_ACC_STATIC = 0x0008
const FILED_ACC_FINAL = 0x0010
const FILED_ACC_VOILATIE = 0x0040
const FILED_ACC_TRANSIENT = 0x0080
const FILED_ACC_SYNTHETIC = 0x1000
const FILED_ACC_ENUM = 0x4000

const CLASS_ACC_PUBLIC = 0x0001
const CLASS_ACC_FINAL = 0x0010
const CLASS_ACC_SUPER = 0x0020 //必选
const CLASS_ACC_INTERFACE = 0x0200
const CLASS_ACC_ABSTRACT = 0x0400
const CLASS_ACC_SYNTHETIC = 0x1000
const CLASS_ACC_ANNOTATION = 0x2000
const CLASS_ACC_ENUM = 0x4000

const (
	METHOD_ACC_PUBLIC       = 0x0001
	METHOD_ACC_PRIVATE      = 0x0002
	METHOD_ACC_PROTECTED    = 0x0004
	METHOD_ACC_STATIC       = 0x0008
	METHOD_ACC_FINAL        = 0x0010
	METHOD_ACC_SYNCHRONIZED = 0x0020
	METHOD_ACC_BRIDGE       = 0x0040
	METHOD_ACC_VARARGS      = 0x0080
	METHOD_ACC_NATIVE       = 0x0100
	METHOD_ACC_ABSTRACT     = 0x0400
	METHOD_ACC_STRICT       = 0x0800
	METHOD_ACC_SYNTHETIC    = 0x1000
)
