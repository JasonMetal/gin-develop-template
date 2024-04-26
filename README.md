# 介绍
`gin-develop-template`

Gin-based web backend api framework for business development

# 1. 初始化项目
替换项目中import的 `develop-template` 为项目名称

# 2. go mod 初始化
go mod init 项目名称

在go.mod文件种加入以下replace信息
```shell
go 1.22

replace (
	github.com/JasonMetal/submodule-support-go.git v0.0.0 => ./submodule/support-go.git
	github.com/JasonMetal/submodule-services-proto.git v0.0.0 => ./submodule/services-proto.git
)
```

# 3. git 子模块初始化

#### 若根目录下无 `submodule` 则新增 `submodule` 文件夹

mkdir submodule

git submodule add git@github.com:JasonMetal/submodule-services-proto.git submodule/services-proto.git

git submodule add git@github.com:JasonMetal/submodule-support-go.git submodule/support-go.git


# 4. 同步子模块
##### 遇到 类似 fatal: A git directory for 'submodule/services-proto.git' is found locally with remote(s)

git submodule add --force git@github.com:JasonMetal/submodule-services-proto.git submodule/services-proto.git

git submodule add --force git@github.com:JasonMetal/submodule-support-go.git submodule/support-go.git

##### 如：遇到 The system cannot find the file specified

git submodule update --remote --init

# 5. 设置私有仓库
go env -w GOPRIVATE=*.github.com

# 6. 整理go mod
go mod tidy

#### 部署运行
#### 0. go build -o go-test8888 cli.go
#### 测试服
#### 1. ./go-test -e test savePageDataCron
#### linux下用 `supervisord` 进行监控相关任务

```shell

directory = /home/www/go-test
command = /home/www/go-develop-template/go-test -e test savePageDataCron
autostart = true
autorestart = true
loglevel = info
stdout_logfile = /var/log/supervisor/goTestCrawl.log
stderr_logfile = /var/log/supervisor/goTestCrawl_stderr.log
stdout_logfile_maxbytes = 30MB
stdout_logfile_backups = 3
stdout_events_enabled = false

```shell
#### 状态
supervisorctl status|grep goTestCrawl
supervisorctl restart goTestCrawl
#### conf /etc/supervisord.d/conf/goTestCrawl.conf
[program:goTestCrawl]
directory = /home/www/go-test
command = /home/www/go-develop-template/go-test8888 -e prod savePageDataCron
autostart = true
autorestart = true
loglevel = info
stdout_logfile = /var/log/supervisor/goTestCrawl.log
stderr_logfile = /var/log/supervisor/goTestCrawl_stderr.log
stdout_logfile_maxbytes = 30MB
stdout_logfile_backups = 3
stdout_events_enabled = false

[xxx@test go-websites]# cat /etc/supervisord.d/conf/goTestCrawl.conf
[program:goTestCrawl]
directory = /home/www/go-test
command = /home/www/go-develop-template/go-test8888 -e prod savePageDataCron
autostart = true
autorestart = true
loglevel = info
stdout_logfile = /var/log/supervisor/goTestCrawl.log
stderr_logfile = /var/log/supervisor/goTestCrawl_stderr.log
stdout_logfile_maxbytes = 30MB
stdout_logfile_backups = 3
stdout_events_enabled = false