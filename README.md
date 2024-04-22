## 程序使用方法

### 第一步:

下载proxyPassword.exe程序,然后用管理员身份运行程序

![image-20240422132139355](https://wxy-md.oss-cn-shanghai.aliyuncs.com/image-20240422132139355.png)

### 第二步:

程序会生成exe程序同目录下生成password_proxy_path.ini配置文件,存放运行程序exe路径
该程序会自动查询注册表里面的程序路径,默认为空需要手动配置

![image-20240422132320960](https://wxy-md.oss-cn-shanghai.aliyuncs.com/image-20240422132320960.png)

### 最后一步:

**xshell浏览器访问示例:**

```
proxyPassword://Soft=xshell&Protocol=ssh&Username=root&Password=123&Port=21&Host=127.0.0.1
```

**自动打开本地软件并且自动连接成功**

![image-20240422133530412](https://wxy-md.oss-cn-shanghai.aliyuncs.com/image-20240422133530412.png)

## 前端对接方式

```
主要通过构建a标签打开新页面 <a href="proxyPassword://Soft=xshell&Protocol=ssh&Username=root&Password=123&Port=21&Host=127.0.0.1" />
```

## 基础命令
删除注册表:
```
reg delete "HKLM\SOFTWARE\Classes\proxyPassword" /f
```
添加注册表:
```
reg.exe import install.reg
```
基础编译:
```
go build -o proxyPassword.exe proxyPassword.go
```
其他Windows架构（如ARM64）或在非Windows系统上编译Windows EXE文件:
```
GOOS=windows GOARCH=amd64 go build -o proxyPassword.exe proxyPassword.go
```
