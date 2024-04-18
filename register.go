package main

import (
	"fmt"
	"github.com/gonutz/w32"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func main() {
	// 操作注册表
	createRegisterTable()
}

// 添加注册表
func createRegisterTable() {
	deleteRegisterTable()
	var currentPath string
	var exePath string
	currentPath, _ = os.Getwd()
	exePath = "\\\\Mac\\Home\\Downloads\\downloadobject\\apipost_win_x64.exe"
	if file, err := os.Create("./proxyPassword.reg"); err != nil {
		w32.MessageBox(w32.HWND(uintptr(0)), "创建注册表文件失败", "失败", w32.MB_ICONERROR|w32.MB_ICONINFORMATION)
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
			w32.MessageBox(w32.HWND(uintptr(0)), "写入注册表失败:"+err.Error(), "失败", w32.MB_ICONERROR|w32.MB_ICONINFORMATION)
		} else {
			w32.MessageBox(w32.HWND(uintptr(0)), "写入注册表成功!", "成功", w32.MB_OK|w32.MB_ICONINFORMATION)
		}
		if err := os.RemoveAll("./proxyPassword.reg"); err != nil {
			fmt.Println(err)
		}
	}
}

// 删除注册表
func deleteRegisterTable() {
	cmd := exec.Command("reg", "delete", "HKLM\\SOFTWARE\\Classes\\proxyPassword", "/f")
	// 隐藏命令行窗口（仅适用于Windows平台）
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err := cmd.Run(); err != nil {
		w32.MessageBox(w32.HWND(uintptr(0)), "移除旧的注册表失败:"+err.Error(), "失败", w32.MB_ICONERROR|w32.MB_ICONINFORMATION)
	} else {
		w32.MessageBox(w32.HWND(uintptr(0)), "移除旧的注册表成功!", "成功", w32.MB_OK|w32.MB_ICONINFORMATION)
	}
}