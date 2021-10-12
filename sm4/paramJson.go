package sm4

import (
	"encoding/json"
	"fmt"
	"log"
)

type Param struct {
	Soft string     //软件
	Protocol string //协议
	Username int //用户名
	Password int //密码
	Port int 		//端口
	Host string     //主机
}


func testMarshal() []byte {
	param := Param{
		Host:"localhost",
		Port:5521,
		Username:281,
		Password:281,
		Soft:"filezilla",
		Protocol:"ftp",
	}
	data, err := json.Marshal(param)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func TestUnmarshal(data []byte) {
	var param Param
	err := json.Unmarshal(data, &param)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(param)
}
func main41() {
/*	var data []byte
	data = testMarshal()
	fmt.Println(string(data))
	TestUnmarshal(data)*/
	var data = `{"Host":"localhost","Port":"5521","Username":281,"Password":281,"Soft":"filezilla","Protocol":"ftp"}`
	p := &Param{}
	json.Unmarshal([]byte(data), p)
	fmt.Println(*p)
}