package main

import (
	"io"
	flag "github.com/spf13/pflag"
	"fmt"
	"bufio"
	"os"
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
	
	if(*start_page < 1){    //起始页小于1
		fmt.Println("You must input the start page and greater than 0, e.g. --s 1")
		return
	}else if (*end_page < 1){ //终止页小于1
		fmt.Println("You must input the end page and greater than 0, e.g. --e 1")
		return
	}else if(*start_page > *end_page){   //起始页大于终止页
		fmt.Println("The number of start_page must not be bigger than end_page")
		return
	}else if(flag.NArg() > 1){  //输入参数出错
		fmt.Println("Please check your input, input selpg -h for help")
		return
	}
	getPrint()
}

func getPrint(){
	var counter int
	var outputContent string

	if(*page_type == true){
		// 以\f为界，逐字符处理
		if(flag.NArg() == 0){
			// 控制台输入
			scanner := bufio.NewScanner(os.Stdin)
			counter = 0
			outputContent = ""
			curPage := 0
 			// 逐字符输入处理
			for scanner.Scan() {
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
			fmt.Println(outputContent)
		}else if(flag.NArg() == 1){
			// 文件输入
			counter = 0
			outputContent = ""
			curPage := 0

			*fileName = flag.Arg(0)
			// 文件输入
			file, err := os.Open(*fileName)
			if(err != nil){
				panic(err)
			}
			reader := bufio.NewReader(file)
			
			for{
				// 逐字符输入处理
				line, _, err := reader.ReadLine()
				str := string(line)
				str += "\n"
				if(err == io.EOF){
					break
				}else if(err != nil){
					panic(err)
				}
				for i:= 0;i < len(str);i++{
					// 遇到换行符，页数+1
					if(line[i] == '\f'){
						curPage++
					}
				}
				if(curPage >= (*start_page)-1 && curPage < (*end_page)) {
					outputContent += str
				}
			}
			fmt.Println(outputContent)
		}
	}else{
		// 每页指定行数，对每一页的输入计行
		if(flag.NArg() == 0){
			scanner := bufio.NewScanner(os.Stdin)
			counter = 0
			outputContent = ""
		
			for scanner.Scan() {
				line := scanner.Text()
				line += "\n"
				curPage := counter/(*page_len)
				if(curPage >= (*start_page)-1 && curPage < (*end_page)) {
					outputContent += line
				}
				counter++
			}
			fmt.Println(outputContent)
		}else if(flag.NArg() == 1){
			counter = 0
			outputContent = ""

			*fileName = flag.Arg(0)
			file, err := os.Open(*fileName)
			if(err != nil){
				panic(err)
			}
			reader := bufio.NewReader(file)
			
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
			fmt.Println(outputContent)
		}
	}
}