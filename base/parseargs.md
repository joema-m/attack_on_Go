# Go解析命令行参数

## os.Args

os.Args的类型是 **[]string** ，也就是字符串切片,可以在for循环的range中遍历，还可以用 **len(os.Args)** 来获取其数量。

```go
a := os.Args
fmt.Println(a[1:])
fmt.Println(len(a))

/*
$ go run parseArgs.go hello world my name is joe
[hello world my name is joe]
7
*/

```



## flag

```go
var bo = flag.Bool("b", false, "bool类型参数") 
// 这里默认为false,如果使用-b，则为true,但如果 -b true，则错误
var str = flag.String("s", "", "string类型参数")

flag.Parse()
// 注意这里使用的是指针
fmt.Println("-b:", *bo)
fmt.Println("-s:", *str)
fmt.Println("其他参数：", flag.Args())

/*
$ go run parseArgs.go -s hello
-b: false
-s: hello
其他参数： []
*/
```

