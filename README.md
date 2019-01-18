# go-selpg
服务计算作业

# 任务要求

根据 [开发 Linux 命令行实用程序](https://www.ibm.com/developerworks/cn/linux/shell/clutil/index.html) 文中的要求用go语言实现selpg

# flag包的使用

作业要求使用标准库—命令行参数解析FLAG，导入包`flag`，这是整个作业最重要的知识点，下面讲一下作业中`flag` 的使用；

## 安装

首先，由于作业要求使用 pflag 替代 goflag 以满足 Unix 命令行规范，因此在在导入包之前需要提前安装pflag，在终端执行代码：

```
go get github.com/spf13/pflag
```

安装成功后，我们就可以使用下面的代码导入pflag包：

```go
import (
	flag "github.com/spf13/pflag"
)
```

## 将flag绑定到变量上

```go
var ip = flag.Int("flagname", 1234, "help message for flagname")
```

上面的代码会将返回一个`int`类型的指针，在flag解析后，`flagname`的值（默认是1234，后面会介绍解析的过程）会自动赋到指针 ip 上，另一种写法是：

```go
var ip int

func init() {
	flag.IntVar(&ip, "flagname", 1234, "help message for flagname")
}
```

## flag 解析

函数`flag.Parse()` 会调用函数解析，你在终端输入的内容会被 `flag.Parse()` 所解析，

以下面的代码为例：

```go
import (
	flag "github.com/spf13/pflag"
)

// 将flag绑定到变量上
var start_page = flag.Int("s", -1, "Input The Page Start")
var end_page = flag.Int("e", -1, "Input The Page End")
var fileName = flag.String("file", "", "Input The FileName")
var page_len = flag.Int("l", 72 , "Input the length of page")
var page_type = flag.Bool("f", false , "Form-feed-delimited")
var destination = flag.String("d","","The destination of the output file")

func main() {
	// flag 解析
	flag.Parse()
}
```

如果你在终端输入：

````
selpg --s 1 --e 1 --d output
````

在flag解析后，指针start_page指向的值会变为1，指针end_page指向的值会变为1，指针destination 指向的值会变为output，没有提供的值会选择它们的默认值，比如说，指针page_len 指向的值是72；

需要注意的是，对于指向布尔值的指针，如page_type，在需要将其置为true的时候不是：

```
selpg --f true
```

实际使用的时候，只需要提供`--f` 即可以，即

```
selpg --f
```

# 项目


根据任务要求，程序selpg需要提供以下功能：

1. -s 指定打印的起始页
2. -e 指定打印的终止页
3. -l 指定打印的每页的行数，默认为72
4. -f 指定换页的方法，以`\f` 为界
5. -d 指定输出的地址，可以是打印机的地址
6. `<` 使用重定向输入
7. `>` 使用重定向输出
8. `2>` 重定向错误输出

具体的参数设计过程可以看上面提供的代码，下面是对输入参数的一些情况进行分析：

1. 起始页小于1
2. 终止页小于1
3. 起始页大于终止页
4. 输入参数出错

```go
if(*start_page < 1){
		fmt.Println("You must input the start page and greater than 0, e.g. --s 1")
		return
	}else if (*end_page < 1){
		fmt.Println("You must input the end page and greater than 0, e.g. --e 1")
		return
	}else if(*start_page > *end_page){
		fmt.Println("The number of start_page must not be bigger than end_page")
		return
	}else if(flag.NArg() > 1){
		fmt.Println("Please check your input, input selpg -h for help")
		return
	}
```

因为`flag.Parse()`可以自行对我们的输入参数进行转化，因此可以直接使用`flao.Parse()` 而不是对输入的命令字符串进行分析，大大简化了编程的困难；

下面就是对打印操作进行处理：

对于不同的输入流进行处理：

```go
// 控制台输入或重定向的内容
scanner := bufio.NewScanner(os.Stdin)
// ....
// 文件输入
file, err := os.Open(*fileName) // 打开文件
if(err != nil){
    panic(err)
}
reader := bufio.NewReader(file)
//...
```

在文件流输入的时候，我们需要对输入进行判断，以`EOF` 作为文件输入的结尾：

```go
// 对页面的结尾进行判断
if(err == io.EOF){
	break
}else if(err != nil){
	panic(err)
}
```

然后就是对不同的页面设置：

第一种设置就是以`\f` 为换页符，此时需要逐字符输入进行处理：

```go
for scanner.Scan() {
    // 逐字符输入处理
    line := scanner.Text()
    line += "\n"
    for i:= 0;i < len(line);i++{
        // 遇到换行符，页数+1
        if(line[i] == '\f'){
            curPage++
        }
    }
    if(curPage >= (*start_page)-1 && curPage < (*end_page)) {
        outputContent += line
    }
}
```

另一种就是每一页的行数固定，因此我们需要使用一个变量计算当前的行数：

```go
for{
    line, _, err := reader.ReadLine()
    str := string(line)
    str += "\n"
    if(err == io.EOF){
        break
    }else if(err != nil){
        panic(err)
    }
    // curPage表示现在的页数
    curPage := counter/(*page_len)
    if(curPage >= (*start_page)-1 && curPage < (*end_page)) {
        outputContent += str
    }
    // 使用counter计行
    counter++
}
```
