package main

import (
	"fmt"
	"github.com/gonutz/w32"
	"gopkg.in/ini.v1"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
)

// 定义软件类型
var soft = []string{"xshell", "xftp", "filezilla"}

func main() {
	// 定义配置文件路径
	exePath, err := os.Executable()
	if err != nil {
		errorMessage("获取当前程序路径失败")
		os.Exit(0)
	}
	configPath := strings.Replace(exePath, "proxyPassword.exe", "password_proxy_path.ini", -1)
	if len(os.Args) == 1 {
		initServer(configPath, exePath)
		return
	}
	// 启动程序
	if len(os.Args) == 7 {
		start(configPath)
		return
	}
	errorMessage("参数信息不正确!")
	os.Exit(0)
}

// 初始化服务
func initServer(configPath string, exePath string) {
	// 操作注册表
	createRegisterTable(exePath)
	// 创建文件
	if !fileExists(configPath) {
		createFileIfNotExist(configPath)
		// 读取INI文件
		cfg, err := ini.LoadSources(ini.LoadOptions{
			SkipUnrecognizableLines: true,
		}, configPath)
		if err != nil {
			errorMessage("无法加载INI文件:" + configPath)
			os.Exit(0)
		}
		// 指定exe程序
		for _, exe := range soft {
			cfg.Section("path").Key(exe).SetValue("")
			if err := cfg.SaveTo(configPath); err != nil {
				errorMessage("无法保存INI文件: %v" + err.Error())
				os.Exit(0)
			}
		}
	}
	successMessage("配置完成,请在配置文件password_proxy_path.ini里填写exe程序地址")
}

type Param struct {
	Soft     string //软件
	Protocol string //协议
	Username string //用户名
	Password string //密码
	Port     string //端口
	Host     string //主机
}

// 启动软件 proxyPassword://Soft=xshell&Protocol=ssh&Username=root&Password=123&Port=21&Host=127.0.0.1
// todo 密码使用非对称加密方式
func start(configPath string) {
	// 获取全路径
	var data = strings.Replace(os.Args[1], "proxypassword://", "", -1)
	// 解析查询字符串
	values, err := url.ParseQuery(data)
	if err != nil {
		errorMessage("读取程序参数失败:" + err.Error())
		os.Exit(0)
	}
	// 创建Param结构体对象并映射查询参数
	p := Param{
		Soft:     values.Get("Soft"),
		Protocol: values.Get("Protocol"),
		Username: values.Get("Username"),
		Password: values.Get("Password"),
		Port:     values.Get("Port"),
		Host:     values.Get("Host"),
	}
	// 读取INI文件
	cfg, err := ini.LoadSources(ini.LoadOptions{
		SkipUnrecognizableLines: true,
	}, configPath)
	if err != nil {
		errorMessage("无法加载INI文件:" + err.Error())
		os.Exit(0)
	}
	kv, err := cfg.Section("path").GetKey(p.Soft)
	if err != nil {
		errorMessage("读取软件路径失败!" + err.Error())
	}
	if !fileExists(kv.Value()) {
		errorMessage("软件路径不存在:" + kv.Value())
	}
	execParam := p.Protocol + "://" + p.Username + ":" + p.Password + "@" + p.Host + ":" + p.Port
	cmd := exec.Command(kv.Value(), execParam)
	// 隐藏Windows平台上命令行窗口
	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	}
	if err := cmd.Run(); err != nil {
		fmt.Println("启动失败:", err)
	} else {
		fmt.Println("启动成功!")
	}
}

// 删除注册表
func deleteRegisterTable() {
	cmd := exec.Command("reg", "delete", "HKLM\\SOFTWARE\\Classes\\proxyPassword", "/f")
	// 隐藏Windows平台上命令行窗口
	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	}
	if err := cmd.Run(); err != nil {
		fmt.Println("移除旧的注册表失败:" + err.Error())
	} else {
		fmt.Println("移除旧的注册表成功!")
	}
}

// 添加注册表
func createRegisterTable(exePath string) {
	deleteRegisterTable()
	var currentPath string
	currentPath, _ = os.Getwd()
	if file, err := os.Create("./proxyPassword.reg"); err != nil {
		errorMessage("创建注册表文件失败:" + err.Error())
		os.Exit(0)
	} else {
		//写入数据
		var data string
		data += "Windows Registry Editor Version 5.00\n\n"
		data += "[HKEY_LOCAL_MACHINE\\SOFTWARE\\Classes\\proxyPassword]\n"
		data += "@=\"proxyPassword\"\n"
		data += "\"URL Protocol\"=\"\"" + "\n\n"
		data += "[HKEY_LOCAL_MACHINE\\SOFTWARE\\Classes\\proxyPassword\\DefaultIcon]\n"
		data += "@=\"%SystemRoot%\\system32\\url.dll,0\"\n\n"
		data += "[HKEY_LOCAL_MACHINE\\SOFTWARE\\Classes\\proxyPassword\\Shell]\n"
		data += "[HKEY_LOCAL_MACHINE\\SOFTWARE\\Classes\\proxyPassword\\Shell\\open]\n"
		data += "[HKEY_LOCAL_MACHINE\\SOFTWARE\\Classes\\proxyPassword\\Shell\\open]\n"
		data += "[HKEY_LOCAL_MACHINE\\SOFTWARE\\Classes\\proxyPassword\\Shell\\open\\command]\n"
		data += "@=\"" + strings.Replace(exePath, "\\", "\\\\", -1) + " "
		data += "\\\"%1\\\" \\\"%2\\\" \\\"%3\\\" \\\"%4\\\" \\\"%5\\\" \\\"%6\\\"\""
		//写入byte的slice数据
		file.Write([]byte(data))
		file.Close()
		// 导入注册表
		cmd := exec.Command("reg.exe", "import", currentPath+"\\proxyPassword.reg")
		// 隐藏Windows平台上命令行窗口
		if runtime.GOOS == "windows" {
			cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		}
		if err := cmd.Run(); err != nil {
			errorMessage("写入注册表失败:" + err.Error())
			os.Exit(0)
		}
		if err := os.RemoveAll("./proxyPassword.reg"); err != nil {
			fmt.Println("删除文件失败:" + err.Error())
		}
		fmt.Println("添加新注册表成功!")
	}
}

// 创建配置文件
func createFileIfNotExist(configPath string) {
	// 分离路径中的目录部分（不含文件名）
	dirName := filepath.Dir(configPath)
	// 检查目录是否存在，不存在则创建
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		err := os.MkdirAll(dirName, 0755)
		if err != nil {
			errorMessage("Error creating directory:" + dirName)
			os.Exit(0)
		}
	}
	// 检查文件是否存在，不存在则创建
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		file, err := os.Create(configPath)
		if err != nil {
			errorMessage("Error creating file:" + configPath)
			os.Exit(0)
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
			}
		}(file)
	}
}

// 文件是否存在
func fileExists(file string) bool {
	_, err := os.Stat(file)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// 成功提示
func successMessage(message string) {
	w32.MessageBox(w32.HWND(uintptr(0)), message, "成功", w32.MB_OK|w32.MB_ICONINFORMATION)
}

// 失败提示
func errorMessage(message string) {
	w32.MessageBox(w32.HWND(uintptr(0)), message, "操作失败", w32.MB_ICONERROR|w32.MB_ICONINFORMATION)
}
