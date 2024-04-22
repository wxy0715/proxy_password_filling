## 基础命令
删除注册表:
```
reg delete "HKLM\SOFTWARE\Classes\db2db" /f
```
添加注册表:
```
reg.exe import install.reg
```
基础编译:
```
go build -o output.exe main.go
```
其他Windows架构（如ARM64）或在非Windows系统上编译Windows EXE文件:
```
GOOS=windows GOARCH=amd64 go build -o proxyPassword.exe proxyPassword.go
```

生成rsrc.syso
```
go get github.com/akavel/rsrc
rsrc -manifest proxy_password_filling.manifest -o rsrc.syso
```


## 程序使用方法
### xshell浏览器访问示例:
```
proxyPassword://Soft=xshell&Protocol=ssh&Username=root&Password=123&Port=21&Host=127.0.0.1
```