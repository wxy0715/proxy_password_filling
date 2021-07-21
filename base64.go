package main
import (
	"encoding/base64"
	"fmt"
	"strings"
)
func baseStdEncode(srcBtye []byte) string {
	encoding := base64.StdEncoding.EncodeToString(srcBtye)
	return encoding
}
func baseDeEncode(src string) string {
	reader := strings.NewReader(src)
	decoder := base64.NewDecoder(base64.StdEncoding, reader)
	// 以流式解码
	buf := make([]byte, 2)
	// 保存解码后的数据
	dst := ""
	for {
		n, err := decoder.Read(buf)
		if n == 0 || err != nil {
			break
		}
		dst += string(buf[:n])
	}
	fmt.Println("解码后的数据为:", dst)
	return dst
}
func main(){
	code:="admin"
	encodeCode:=baseStdEncode([]byte(code))
	println(code+" base64加密后的数据为: "+encodeCode)
	println(encodeCode+" base64解密后的数据为: "+baseDeEncode("eyJIb3N0IjoibG9jYWxob3N0IiwiUG9ydCI6IjU1MjEiLCJVc2VybmFtZSI6MjgxLCJQYXNzd29yZCI6MjgxLCJTb2Z0IjoiZmlsZXppbGxhIiwiUHJvdG9jb2wiOiJmdHAifQ=="))

}