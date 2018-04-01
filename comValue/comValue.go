package comValue

const NOP = 0x00         //空操作
const ACONST_NULL = 0x01 //null值入栈
const ICONST_M1 = 0x02   //-1(int)入栈
const ICONST_0 = 0x03    //0(int)入栈，下同
const ICONST_1 = 0x04
const ICONST_2 = 0x05
const ICONST_3 = 0x06
const ICONST_4 = 0x07
const ICONST_5 = 0x08
const LCONST_0 = 0x09 //0(long)入栈
const LCONST_1 = 0x0a
const FCONST_0 = 0x0b //0(float)入栈
const FCONST_1 = 0x0c
const FCONST_2 = 0x0d
const DCONST_0 = 0x0e //0(double)入栈
const DCONST_1 = 0x0f
const BIPUSH = 0x10  //操作数byte,拓展成int型入栈
const SIPUSH = 0x11  //操作数int16,拓展成int型入栈
const LDC = 0x12     //操作数byte,将常量(int,float,string)入栈
const LDC_W = 0x13   //操作数uint16,将常量入栈
const LDC2_W = 0x14  //操作数uint16,将常量(long,double)入栈
const ILOAD = 0x15   //操作数byte,从局部变量装载int类型入栈,可拓展
const LLOAD = 0x16   //操作数byte,从局部变量装载long类型入栈,可拓展
const FLOAD = 0x17   //操作数byte,从局部变量装载float类型入栈,可拓展
const DLOAD = 0x18   //操作数byte,从局部变量装载double类型入栈，可拓展
const ALOAD = 0x19   //操作数byye,从局部变量装载引用类型入栈,可拓展
const ILOAD_0 = 0x1a //局部变量0(int)入栈
const ILOAD_1 = 0x1b
const ILOAD_2 = 0x1c
const ILOAD_3 = 0x1d
const LLOAD_0 = 0x1e //局部变量0(long)
const LLOAD_1 = 0x1f
const LLOAD_2 = 0x20
const LLOAD_3 = 0x21
const FLOAD_0 = 0x22 //局部变量0(float)入栈
const FLOAD_1 = 0x23
const FLOAD_2 = 0x24
const FLOAD_3 = 0x25
const DLOAD_0 = 0x26 //局部变量0(double)入栈
const DLOAD_1 = 0x27
const DLOAD_2 = 0x28
const DLOAD_3 = 0x29
const ALOAD_0 = 0x2a //局部变量0(引用)入栈
const ALOAD_1 = 0x2b
const ALOAD_2 = 0x2c
const ALOAD_3 = 0x2d
const IALOAD = 0x2e //装载int数组的指定项
const LALOAD = 0x2f
const FALOAD = 0x30
const DALOAD = 0x31
const AALOAD = 0x32
const BALOAD = 0x33   //装载boolean或byte数组的指定项(拓展成int型)
const CALOAD = 0x34   //装载char数组的指定项(拓展成int型)
const SALOAD = 0x35   //装载short数组的指定项(拓展成int型)
const ISTORE = 0x36   //操作数byte,栈顶元素保存到局部变量(int)
const LSTORE = 0x37   //操作数byte,栈顶元素保存到局部变量(long)
const FSTORE = 0x38   //操作数byte,栈顶元素保存到局部变量(float)
const DSTORE = 0x39   //操作数byte,栈顶元素保存到局部变量(double)
const ASTORE = 0x3a   //操作数byte,栈顶元素保存到局部变量(引用)
const ISTORE_0 = 0x3b //栈顶元素保存到局部变量0(int)
const ISTORE_1 = 0x3c
const ISTORE_2 = 0x3d
const ISTORE_3 = 0x3e
const LSTORE_0 = 0x3f
const LSTORE_1 = 0x40
const LSTORE_2 = 0x41
const LSTORE_3 = 0x42
const FSTORE_0 = 0x43
const FSTORE_1 = 0x44
const FSTORE_2 = 0x45
const FSTORE_3 = 0x46
const DSTORE_0 = 0x47
const DSTORE_1 = 0x48
const DSTORE_2 = 0x49
const DSTORE_3 = 0x4a
const ASTORE_0 = 0x4b
const ASTORE_1 = 0x4c
const ASTORE_2 = 0x4d
const ASTORE_3 = 0x4e
const IASTORE = 0x4f //栈顶int元素保存到数组中
const LASTORE = 0x50
const FASTORE = 0x51
const DASTORE = 0x52
const AASTORE = 0x53
const BASTORE = 0x54
const CASTORE = 0x55
const SASTORE = 0x56
const POP = 0x57     //从栈顶弹出一个字长的元素
const POP2 = 0x58    //从栈顶弹出两个字长的元素
const DUP = 0x59     //复制栈顶一个元素并压栈
const DUP_X1 = 0x5a  //复制栈顶一个字长的数据，弹出栈顶两个字长数据，先将复制后的数据压栈，再将弹出的两个字长数据压栈
const DUP_X2 = 0x5b  //复制栈顶一个字长的数据，弹出栈顶三个字长的数据，将复制后的数据压栈，再将弹出的三个字长的数据压栈
const DUP2 = 0x5c    //复制栈顶两个字长的数据，将复制后的两个字长的数据压栈
const DUP2_X1 = 0x5d //复制栈顶两个字长的数据，弹出栈顶三个字长的数据，将复制后的两个字长的数据压栈，再将弹出的三个字长的数据压栈
const DUP2_X2 = 0x5e //复制栈顶两个字长的数据，弹出栈顶四个字长的数据，将复制后的两个字长的数据压栈，再将弹出的四个字长的数据压栈
const SWAP = 0x5f    //交换栈顶两个字长的数据的位置。Java指令中没有提供以两个字长为单位的交换指令
const IADD = 0x60    //int+int,结果入栈
const LADD = 0x61    //long+long,结果入栈
const FADD = 0x62    //float+float,结果入栈
const DADD = 0x63    //double+double,结果入栈
const ISUB = 0x64    //int-int,结果入栈
const LSUB = 0x65    //long-long.结果入栈
const FSUB = 0x66    //float-float,结果入栈
const DSUB = 0x67    //double-double,结果入栈
const IMUL = 0x68    //int*int,结果入栈
const LMUL = 0x69    //long*long,结果入栈
const FMUL = 0x6a    //float*float,结果入栈
const DMUL = 0x6b    //double*double,结果入栈
const IDIV = 0x6c    //int/int,结果入栈
const LDIV = 0x6d    //long/long,结果入栈
const FDIV = 0x6e    //float/float,结果入栈
const DDIV = 0x6f    //double/double,结果入栈
const IREM = 0x70    //int%int,结果入栈
const LREM = 0x71    //long%long,结果入栈
const FREM = 0x72    //float%float,结果入栈
const DREM = 0x73    //double%double,结果入栈
const INEG = 0x74    //int取负,结果入栈
const LENG = 0x75    //long取负,结果入栈
const FNEG = 0x76    //float取负,结果入栈
const DNEG = 0x77    //double取负,结果入栈
const ISHL = 0x78    //左移int类型
const LSHL = 0x79    //左移long类型
const ISHR = 0x7a    //算数右移int类型
const LSHR = 0x7b    //算数右移long类型
const IUSHR = 0x7c   //逻辑右移int类型
const LUSHR = 0x7d   //逻辑右移long类型
const IAND = 0x7e    //int按位与
const LAND = 0x7f    //long按位与
const IOR = 0x80     //int按位或
const LOR = 0x81     //long按位或
const IXOR = 0x82    //int按位异或
const LXOR = 0x83    //long按位异或
const IINC = 0x84    //操作数,indexByte,constbyte,将整数值constbyte加到indexbyte指定的int类型的局部变量中
const I2L = 0x85     //将栈顶的int转成long
const I2F = 0x86
const I2D = 0x87
const L2I = 0x88
const L2F = 0x89
const L2D = 0x8a
const F2I = 0x8b
const F2L = 0x8c
const F2D = 0x8d
const D2I = 0x8e
const D2L = 0x8f
const D2F = 0x90
const I2B = 0x91 //int->byte
const I2C = 0x92
const I2S = 0x93
const LCMP = 0x94            //比较栈顶两long类型值，前者大，1入栈；相等，0入栈；后者大，-1入栈
const FCMPL = 0x95           //比较栈顶两float类型值，前者大，1入栈；相等，0入栈；后者大，-1入栈；有NaN存在，-1入栈
const FCMPG = 0x96           //比较栈顶两float类型值，前者大，1入栈；相等，0入栈；后者大，-1入栈；有NaN存在，-1入栈
const DCMPL = 0x97           //比较栈顶两double类型值，前者大，1入栈；相等，0入栈；后者大，-1入栈；有NaN存在，-1入栈
const DCMPG = 0x98           //比较栈顶两double类型值，前者大，1入栈；相等，0入栈；后者大，-1入栈；有NaN存在，-1入栈
const IFEQ = 0x99            //操作数,uint16,栈顶等于0则跳转
const IFNE = 0x9a            //操作数,uint16,栈顶不等于0则跳转
const IFLT = 0x9b            //操作数,uint16,栈顶小于0则跳转
const IFGE = 0x9c            //操作数,uint16,栈顶大于等于0则跳转
const IFLE = 0x9e            //操作数,uint16,栈顶小于等于0则跳转
const IF_ICMPEQ = 0x9f       //操作数,uint16,栈顶两元素相等则跳转
const IF_ICMPNE = 0xa0       //操作数,uint16,栈顶两元素不相等则跳转
const IF_ICMPLT = 0xa1       //操作数,uint16,若栈顶两int类型值前小于后则跳转
const IF_ICMPGE = 0xa2       //操作数,uint16,若栈顶两int类型值前大于等于后则跳转
const IF_ICMPGT = 0xa3       //操作数,uint16,若栈顶两int类型值前大于后则跳转
const IF_ICMPLE = 0xa4       //操作数,uint16,若栈顶两int类型值前小于等于后则跳转
const IF_ACMPEQ = 0xa5       //操作数,uint16,若栈顶两引用类型值相等则跳转
const IF_ACMPNE = 0xa6       //操作数,uint16,若栈顶两引用类型值不相等则跳转
const GOTO = 0xa7            //操作数,uint16,无条件跳转
const JSR = 0xa8             //操作数,uint16,跳转到子例程序,JDK1.4之后不再编译出该指令，java7之后禁用该指令，故不实现
const RET = 0xa9             //同上，废弃
const TABLESWITCH = 0xaa     //跳转表指令(索引),操作数略多
const LOOKUPSWITCH = 0xab    //跳转表指令(索引),操作数略多
const IRETURN = 0xac         //return int类型
const LRETURN = 0xad         //return long类型
const FRETURN = 0xae         //return float类型
const DRETURN = 0xaf         //return double类型
const ARETURN = 0xb0         //return 引用类型
const RETURN = 0xb1          //return void
const GETSTATIC = 0xb2       //操作数,uint16,获取静态字段的值
const PUTSTATIC = 0xb3       //操作数,uint16,给静态字段赋值
const GETFIELD = 0xb4        //操作数,uint16,获取对象字段的值
const PUTFIELD = 0xb5        //操作数,uint16,给对象字段赋值
const INVOKEVIRTUAL = 0xb6   //操作数,uint16,运行时调用方法，从实例类开始一路往父类找实现方法
const INVOKESPECIAL = 0xb7   //操作数,uint16,调用指定方法(构造函数、父类的方法、私有方法)
const INVOKESTATIC = 0xb8    //操作数,uint16,调用静态方法
const INVOKEINTERFACE = 0xb9 //操作数,uint16,调用接口方法
const NEW = 0xbb             //操作数,uint16,new一个对象
const NEWARRAY = 0xbc        //操作数,byte(指示数据类型),新建一个基本类型数组
const ANEWARRAY = 0xbd       //操作数,uint16,新建引用类型数组
const ARRAYLENGTH = 0xbe     //获取数组长度
const ATHROW = 0xbf          //抛异常
const CHECKCAST = 0xc0       //操作数,uint16,检查类型是否可以强转，不可以则抛异常
const INSTANCEOF = 0xc1      //操作数,uint16,instance关键字的实现，结果压栈
const MONITORENTER = 0xc2    //进入并获得对象监视器
const MONITOREXIT = 0xc3     //释放并退出对象监视器
const WIDE = 0xc4            //使用附加字节扩展局部变量索引
const MULTIANEWARRAY = 0xc5  //操作数,uint16,byte,创建多维数组
const IFNULL = 0xc6          //操作数,uint16,如果栈顶元素为null，则跳转
const IFNONNULL = 0xc7       //操作数,uint16,如果栈顶元素不为null，则跳转
const GOTO_W = 0xc8          //操作数,byte,byte,byte,byte,无条件跳转
const JSR_W = 0xc9           //不实现
