package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	//创建文件
	if file,err := os.Create("./a.txt"); err!=nil{
		fmt.Println(err)
	}else{
		//写入数据
		data:="我是王星宇\n" +
			"hello"
		//写入byte的slice数据
		file.Write([]byte(data))
		//写入字符串
		file.WriteString(data)
		file.Close()
	}

	// 读取文件
	if file, err := os.Open("./a.txt"); err!=nil {
		fmt.Println(err)
	}else{
		buf2 := make([]byte, 1024)
		ix := 0
		for {
			//ReadAt从指定的偏移量开始读取，不会改变文件偏移量
			len, _ := file.ReadAt(buf2, int64(ix))
			ix = ix + len
			if len == 0 {
				break
			}
			fmt.Println(string(buf2))
		}
		file.Close()
	}
}

func readLine(fileName string) ([]string,error){
	f, err := os.Open(fileName)
	if err != nil {
		return nil,err
	}
	buf := bufio.NewReader(f)
	var result []string
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			if err == io.EOF {   //读取结束，会报EOF
				return result,nil
			}
			return nil,err
		}
		result = append(result,line)
	}
	return result,nil
}

//判断文件或者文件夹是否存在，一般判断第一个参数即可，第二个参数可以忽略，或者严谨一些，把err日志记录起来
func FileExists(file string) (bool, error) {
	_, err := os.Stat(file)
	if err == nil {
		return true, nil //文件或者文件夹存在
	}
	if os.IsNotExist(err) {
		return false, nil //不存在
	}
	return false, err //不存在，这里的err可以查到具体的错误信息
}

//判断目录是否存在
func isDir(dir string) bool {
	info, err := os.Stat(dir)
	if err == nil {
		return false
	}
	return info.IsDir()
}

//判断文件是否存在
func isFile(file string) bool {
	info, err := os.Stat(file)
	if err != nil {
		return false
	}
	return !info.IsDir()
}
