## ProcZygote
一个轻量化的分布式子程序调用程序

### 使用说明
支持两种运行方式：

- 前台运行，直接使用`./ProcZygote start`
- 后台运行，使用`./ProcZygote -daemon -out 日志文件路径`

**特别注意：**

1. 如需更改端口，请查看`constants/default.go`
2. 首次运行，需手动创建`/run/ProcZygote/logs`目录
3. 如需使用userNS，则需使用root权限运行本程序