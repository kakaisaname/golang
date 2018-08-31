Go builtin包提供了go预先声明的函数、变量等的文档。这些函数变量等的实现其实并不是在builtin包里，只是为了方便文档组织。

这些内置的变量、函数、类型无需引入包即可使用。



默认提供的有：

1、常量：

true,false,iota

关于iota,请参考周一的文章“go iota用法”



2、变量

nil



3、函数

func append(slice []Type, elems ...Type) []Type

向slice中添加元素



func cap(v Type) int

返回指定类型的容量，根据不同类型，返回意义不同。

数组: 元素个数 (和len(v)一样).

数组指针: *v的元素个数 (和len(v)一样).

Slice: the maximum length the slice can reach when resliced;如果v==nil, cap(v) 值为0；

Channel: channel 缓存区的容量, 以其中的元素为单位;如果v==nil, cap(v) 值为0；

参考之前的文章“Go内置函数cap”



func close(c chan<- Type)

关闭一个channel



func complex(r, i FloatType) ComplexType

创建一个复数



func copy(dst, src []Type) int

用于slice间复制数据，参考之前的文章“Go内建函数copy”



func delete(m map[Type]Type1, key Type)

删除map中指定key的元素



func imag(c ComplexType) FloatType

返回复数的虚部



func len(v Type) int

返回变量的长度。参考之前的文章“Go内置函数len”



func make(Type, size IntegerType) Type



func new(Type) *Type



func panic(v interface{})

产生一个异常，参考"go异常处理"



func print(args ...Type)

打印输出，可用于调试



func println(args ...Type)

打印输出，可用于调试



func real(c ComplexType) FloatType

返回复数的实部



func recover() interface{}

参考"go异常处理"





4、数据类型

type ComplexType

type FloatType

type IntegerType

type Type

type Type1

type bool

type byte

type complex128

type complex64

type error

type float32

type float64

type int

type int16

type int32

type int64

type int8

type rune

type string

type uint

type uint16

type uint32

type uint64

type uint8

type uintptr