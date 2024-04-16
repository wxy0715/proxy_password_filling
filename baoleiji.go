package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"io"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"syscall"
)
type Param struct {
	Soft string     //软件
	Protocol string //协议
	Username string //用户名
	Password string //密码
	Port string 	//端口
	Host string     //主机
}

type MyMainWindow struct {
	*walk.MainWindow
	edit *walk.TextEdit
}

var path = ""
var isTrue = false
var paramNumber = ""
var mw = &MyMainWindow{}
func main1() {
	createRegisterTable()
}

func run(){
	// 1.先判断配置文件是否存在
	userDir, err := user.Current()
	if err != nil {
		fmt.Println(err)
	}
	configPath := userDir.HomeDir+"\\AppData\\Roaming\\password_proxy\\password_proxy.config"
	pathTemp := userDir.HomeDir+"\\AppData\\Roaming\\password_proxy"
	pathTemp = strings.Replace(pathTemp,"\\","\\\\",-1)
	err = os.Mkdir(pathTemp, 0666)
	if err != nil {
		fmt.Println("文件存在")
	}
	if exist ,err := fileExists(configPath);err!=nil{
		fmt.Println(err)
	}else if !exist {
		if file,err := os.Create(configPath); err!=nil{
			fmt.Println(err)
		} else{
			file.Close()
		}
	}
	// 2.判断文件里面配置的路径是否存在
	var execPath = ""
	f, err := os.Open(configPath)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	isTrue = false
	for true {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if  strings.Contains(line,"filezilla=") {
			execPath = strings.Replace(line,"filezilla=","",-1)
			if exist ,err := fileExists(execPath);err!=nil{
				fmt.Println(err)
			}else if !exist {
				MainWindow{
					AssignTo: &mw.MainWindow,  //窗口重定向至mw，重定向后可由重定向变量控制控件
					//Icon:     "./Ico.ico", //窗体图标
					Title:    "文件选择对话框",       //标题
					MinSize:  Size{Width: 500, Height: 130},
					Size:     Size{Width: 500, Height: 130},
					Layout:   VBox{}, //样式，纵向
					Children: []Widget{ //控件组
						Composite{
							Layout: Grid{Columns: 2},
							Children: []Widget{
								TextEdit{
									AssignTo: &mw.edit,
								},
								Composite{
									Layout: Grid{Rows: 2},
									Children: []Widget{
										PushButton{
											Text:      "浏览",
											OnClicked: mw.selectFile,
										},
										PushButton{
											Text:      "确定",
											OnClicked: mw.ensurePath,
										},
									},
								},
							},
						},
					},
				}.Run()
				//mw.MainWindow.Run()
			}
			isTrue = true
			break
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
		}
	}
	fmt.Println(os.Getwd())
	if !isTrue {
		MainWindow{
			AssignTo: &mw.MainWindow,  //窗口重定向至mw，重定向后可由重定向变量控制控件
			//Icon:     "./Ico.ico", //窗体图标
			Title:    "文件选择对话框",       //标题
			MinSize:  Size{Width: 500, Height: 130},
			Size:     Size{Width: 500, Height: 130},
			Layout:   VBox{}, //样式，纵向
			Children: []Widget{ //控件组
				Composite{
					Layout: Grid{Columns: 2},
					Children: []Widget{
						TextEdit{
							AssignTo: &mw.edit,
						},
						Composite{
							Layout: Grid{Rows: 2},
							Children: []Widget{
								PushButton{
									Text:      "浏览",
									OnClicked: mw.selectFile,
								},
								PushButton{
									Text:      "确定",
									OnClicked: mw.ensurePath,
								},
							},
						},
					},
				},
			},
		}.Run()
		//mw.MainWindow.Run()
	}else{
		var data = ""
		if os.Args==nil {
			data = baseDeEncode(strings.Replace(os.Args[1],"password_proxy://","",-1))
		}else{
			paramNumber = baseDeEncode(strings.Replace(os.Args[1],"password_proxy://","",-1))
			data = baseDeEncode(strings.Replace(os.Args[1],"password_proxy://","",-1))
		}
		p := &Param{}
		json.Unmarshal([]byte(data), p)
		// exec
		execParam := p.Protocol+"://"+p.Username+":"+p.Password+"@"+p.Host+":"+p.Port
		cmd := exec.Command(execPath,execParam)
		if err := cmd.Run(); err != nil {
			log.Println("启动失败:", err)
		} else {
			log.Println("启动成功!")
		}
	}
}
func main() {
	run()
}

/*func main2(){
	a := []byte("admin")
	key := []byte("A55BC477B1C0E79E")
	decrypto := SM4Encrypt(a,key)
	fmt.Println("sm4加密后：",hex.EncodeToString(decrypto))
	i := SM4Decrypt(decrypto, key)
	fmt.Println("sm4解密后：",string(i))
}*/

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
	return dst
}

func (mw *MyMainWindow) selectFile() {
	dlg := new(walk.FileDialog)
	dlg.Title = "选择文件"
	dlg.Filter = "可执行文件 (*.exe)|*.exe|所有文件 (*.*)|*.*"
	mw.edit.SetText("") //通过重定向变量设置TextEdit的Text
	if ok, err := dlg.ShowOpen(mw); err != nil {
		mw.edit.AppendText("Error : File Open\r\n")
		return
	} else if !ok {
		mw.edit.AppendText("Cancel\r\n")
		return
	}
	s := fmt.Sprintf("Select : %s\r\n", dlg.FilePath)
	path = dlg.FilePath
	mw.edit.AppendText(s)
}

func (mw *MyMainWindow) ensurePath() {
	createConfigFile(path,"filezilla")
}

// 创建配置文件
func createConfigFile(path string,name string){
	if file,err := os.Create(getUserPath()); err!=nil{
		fmt.Println(err)
	} else{
		//写入数据
		var data string
		data += name+"="+path
		//写入byte的slice数据
		file.Write([]byte(data))
		file.Close()
		mw.MainWindow.Close()
		run()
	}
}

func getUserPath() string {
	userDir, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	pathTemp := userDir.HomeDir+"\\AppData\\Roaming\\password_proxy"
	pathTemp = strings.Replace(pathTemp,"\\","\\\\",-1)
	err = os.Mkdir(pathTemp, 0666)
	if err != nil {
		fmt.Println("文件存在")
	}
	configPath := userDir.HomeDir+"\\AppData\\Roaming\\password_proxy\\password_proxy.config"
	return strings.Replace(configPath,"\\","\\\\",-1)
}

func fileExists(file string) (bool, error) {
	_, err := os.Stat(file)
	if err == nil {
		return true, nil //文件或者文件夹存在
	}
	if os.IsNotExist(err) {
		return false, nil //不存在
	}
	return false, err //不存在，这里的err可以查到具体的错误信息
}

// 添加注册表
func createRegisterTable() {
	var str string
	//var result string
	str, _ = os.Getwd()
	str = str+"\\\\baoleiji.exe"
	if file,err := os.Create("./install.reg"); err!=nil{
		fmt.Println(err)
	} else{
		//写入数据
		var data string
		data += "Windows Registry Editor Version 5.00\n\n"
		data += "[HKEY_LOCAL_MACHINE\\SOFTWARE\\Classes\\password_proxy]\n"
		data += "@=\"password_proxy\"\n"
		data += "\"URL Protocol\"=\""
		data += strings.Replace(strings.Replace(str, "\\", "\\\\", -1), "\\\\baoleiji.exe", "", -1)+"\"\n\n"
		data += "[HKEY_LOCAL_MACHINE\\SOFTWARE\\Classes\\password_proxy\\DefaultIcon]\n"
		data += "@=\"%SystemRoot%\\system32\\url.dll,0\"\n\n"
		data += "[HKEY_LOCAL_MACHINE\\SOFTWARE\\Classes\\password_proxy\\Shell]\n"
		data += "[HKEY_LOCAL_MACHINE\\SOFTWARE\\Classes\\password_proxy\\Shell\\open]\n"
		data += "[HKEY_LOCAL_MACHINE\\SOFTWARE\\Classes\\password_proxy\\Shell\\open]\n"
		data += "[HKEY_LOCAL_MACHINE\\SOFTWARE\\Classes\\password_proxy\\Shell\\open\\command]\n"
		data += "@=\"\\\""
		data += strings.Replace(strings.Replace(str, "\\", "\\\\", -1), "\\\\baoleiji.exe", "baoleiji.exe", -1)+"\\\""+" "
		data += "\\\"%1\\\" \\\"%2\\\" \\\"%3\\\" \\\"%4\\\" \\\"%5\\\"\""
		//写入byte的slice数据
		file.Write([]byte(data))
		file.Close()
		var result string
		// 导入注册表
		result, _ = os.Getwd()
		cmd := exec.Command("cmd")
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		in := bytes.NewBuffer(nil)
		cmd.Stdin = in
		var out bytes.Buffer
		cmd.Stdout = &out
		go func() {
			in.WriteString("reg.exe import \""+result+"\\install.reg\"\n")
		}()
		if err := cmd.Run(); err != nil {
			log.Println("导入失败:", err)
		} else {
			log.Println("导入成功!")
		}
	}
}