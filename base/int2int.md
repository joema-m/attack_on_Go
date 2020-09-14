# Go语言的int强制类型转换

## 0 基础

先看看各数值范围

```go
// 输出各数值范围
fmt.Println("int8 range:", math.MinInt8, math.MaxInt8)
fmt.Println("uint8 range:", 0, math.MaxUint8)
fmt.Println("int16 range:", math.MinInt16, math.MaxInt16)
fmt.Println("uint16 range:", 0, math.MaxInt16)
fmt.Println("int32 range:", math.MinInt32, math.MaxInt32)
fmt.Println("uint32 range:", 0, math.MaxUint32)
fmt.Println("int64 range:", math.MinInt64, math.MaxInt64)
// 下面这句报错，constant 18446744073709551615 overflows int
//fmt.Println("uint64 range:", 0, math.MaxUint64)

/*
int8 range: -128 127
uint8 range: 0 255
int16 range: -32768 32767
uint16 range: 0 32767
int32 range: -2147483648 2147483647
uint32 range: 0 4294967295
int64 range: -9223372036854775808 9223372036854775807
*/
```



## 1. 类型转换

### 1.1 低转高

```go
var a int8 = 0x80 - 1
var b int16 = int16(a)
fmt.Println(a, b)
// 127 127
var c int16 = 0x8000 - 1
var d int32 = int32(c)
fmt.Println(c, d)
// 32767 32767
```

数值不变

### 1.2 高转低

```go
var e int16 = 0x1234
var f int8 = int8(e)
fmt.Println(e, f)
// 4660 52
// ox34 = 52

var g int32 = 0x12345678
var h int16 = int16(g)
var i int16 = 0x5678
fmt.Println(g, h, i)
// 305419896 22136 22136
```

可见由高转低是截取低位的部分

### 1.3 重点来了

再来看看下面的

```go
var j int16 = 0x0081
var k int8 = int8(j)
fmt.Println(j, k)
// 129 -127
```

预期结果应该是`129 129`，但是第二个数却是`-127`，这是怎么回事？

想一想补码的概念就应该猜到了

```
0x81，把它看作带符号的二进制
二进制：1000 0001
反码：  1111 1110
补码：  1111 1111
十进制：-127
```

看看32转16

```go
var l int32 = 0x00008001
var m int16 = int16(l)
fmt.Println(l, m)
// 32769 -32767
```

