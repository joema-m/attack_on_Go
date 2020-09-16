# Go爬取B站UP所有专栏的图片

上次分析了怎么获取B站UP的所有专栏id，在这里：[link](../base/json_parse.md)。这次来进一步获取所有的图片，并存到`MongoDB`或文件中

`CentOS`安装`MongoDB`可以看这里：[link](https://github.com/joema-m/attack_on_Python/blob/master/database/CentOS%E5%AE%89%E8%A3%85MongoDB.md)

代码在这里：[biliTdb.go](./biliTdb.go)

### 知识点：

* 分析URL
* 解析JSON
* 操作MongoDB
* 文件读写

### 需要改进的：

* 数据库连接池
* 健壮性
* 读取数据库连接配置
* 添加保存路径
* 并发获取专栏的图片，涉及到共享变量，还需要多看看协程
* ......

### Next:

既然都获取了文件名，下次就可以并发下载了。

