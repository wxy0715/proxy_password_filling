package main

import (
	"github.com/gonutz/w32"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"
)

var log = logrus.New()

// 定义软件
var soft = []string{"xshell", "xftp", "filezilla"}

func main() {
	// 获取用户目录
	userDir, err := user.Current()
	if err != nil {
		errorMessage("获取用户目录失败:" + err.Error())
		os.Exit(0)
	}
	// 定义配置文件路径
	configPath := userDir.HomeDir + "\\AppData\\Roaming\\password_proxy\\password_proxy_path.ini"
	logPath := userDir.HomeDir + "\\AppData\\Roaming\\password_proxy\\password_proxy.log"
	configPath = "./password_proxy_path.ini"
	logPath = "./password_proxy.log"
	// 判断入参是注册还是启动
	if os.Args == nil {
		// 初始化日志
		initLog(logPath)
		// 操作注册表
		createRegisterTable()
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
	} else {
		if len(os.Args) != 6 {
			errorMessage("参数信息不正确!")
		}
		start(configPath)
	}
}

// 启动软件
func start(configPath string) {
	// 根据入参启动程序
	Soft := os.Args[1]
	Protocol := os.Args[2]
	Username := os.Args[3]
	Password := os.Args[4]
	Host := os.Args[5]
	Port := os.Args[6]
	// 通过Soft获取路径
	// 读取INI文件
	cfg, err := ini.LoadSources(ini.LoadOptions{
		SkipUnrecognizableLines: true,
	}, configPath)
	if err != nil {
		errorMessage("无法加载INI文件:" + configPath)
		os.Exit(0)
	}
	kv, err := cfg.Section("path").GetKey(Soft)
	if err != nil {
		errorMessage("读取软件路径失败!" + err.Error())
	}
	if !fileExists(kv.Value()) {
		errorMessage("软件路径不存在:" + kv.Value())
	}
	execParam := Protocol + "://" + Username + ":" + Password + "@" + Host + ":" + Port
	cmd := exec.Command(kv.Value(), execParam)
	if err := cmd.Run(); err != nil {
		log.Println("启动失败:", err)
	} else {
		log.Println("启动成功!")
	}
}

// 删除注册表
func deleteRegisterTable() {
	cmd := exec.Command("reg", "delete", "HKLM\\SOFTWARE\\Classes\\proxyPassword", "/f")
	// 隐藏命令行窗口（仅适用于Windows平台）
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err := cmd.Run(); err != nil {
		errorMessage("移除旧的注册表失败:" + err.Error())
		os.Exit(0)
	} else {
		log.Info("移除旧的注册表成功!")
	}
}

// 添加注册表
func createRegisterTable() {
	deleteRegisterTable()
	var currentPath string
	var exePath string
	currentPath, _ = os.Getwd()
	exePath = "\\\\Mac\\Home\\Desktop\\proxyPassword.exe"
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
		//data += strings.Replace(strings.Replace(str, "\\", "\\\\", -1), "\\\\proxyPassword.exe", "", -1) + "\"\n\n"
		data += "[HKEY_LOCAL_MACHINE\\SOFTWARE\\Classes\\proxyPassword\\DefaultIcon]\n"
		data += "@=\"%SystemRoot%\\system32\\url.dll,0\"\n\n"
		data += "[HKEY_LOCAL_MACHINE\\SOFTWARE\\Classes\\proxyPassword\\Shell]\n"
		data += "[HKEY_LOCAL_MACHINE\\SOFTWARE\\Classes\\proxyPassword\\Shell\\open]\n"
		data += "[HKEY_LOCAL_MACHINE\\SOFTWARE\\Classes\\proxyPassword\\Shell\\open]\n"
		data += "[HKEY_LOCAL_MACHINE\\SOFTWARE\\Classes\\proxyPassword\\Shell\\open\\command]\n"
		data += "@=\"" + strings.Replace(exePath, "\\", "\\\\", -1) + " "
		data += "\\\"%1\\\" \\\"%2\\\" \\\"%3\\\" \\\"%4\\\" \\\"%5\\\"\""
		//写入byte的slice数据
		file.Write([]byte(data))
		file.Close()
		// 导入注册表
		cmd := exec.Command("reg.exe", "import", currentPath+"\\proxyPassword.reg")
		// 隐藏命令行窗口（仅适用于Windows平台）
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		if err := cmd.Run(); err != nil {
			errorMessage("写入注册表失败:" + err.Error())
			os.Exit(0)
		}
		if err := os.RemoveAll("./proxyPassword.reg"); err != nil {
			log.Error("删除文件失败:" + err.Error())
		}
		log.Info("添加新注册表成功!")
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
				log.Error("Error closing file: %v\n", err)
			}
		}(file)
	}
	log.Info("Directory structure and/or file already exist or created successfully.")
}

// 初始化日志
func initLog(logPath string) {
	log.Formatter = &logrus.JSONFormatter{}
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		errorMessage("Failed to open log file: " + err.Error())
		os.Exit(0)
	}
	log.Out = file
}

// 成功提示
func successMessage(message string) {
	w32.MessageBox(w32.HWND(uintptr(0)), message, "成功", w32.MB_OK|w32.MB_ICONINFORMATION)
}

// 失败提示
func errorMessage(message string) {
	w32.MessageBox(w32.HWND(uintptr(0)), message, "操作失败", w32.MB_ICONERROR|w32.MB_ICONINFORMATION)
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
