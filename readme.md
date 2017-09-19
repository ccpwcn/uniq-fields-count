# 1. 使用示例：
## Windows环境
```dos
.\uniq-fields-count.exe -f .\alive-one-day-log.example1.log -r "\""userId\"":\""[a-zA-Z0-9\-]+\"""
```
> 因为Windows命令行上使用双引号执行多重组合，所以要进行转义

## Linux环境
```bash
./uniq-fields-count -f ./alive-one-day-log.example3.log -r '"userId":"[a-zA-Z0-9\-]+"'
```

生产环境服务器高优先级启动：
```bash
nice --10 ./uniq-fields-count -f /opt/web_app/tomcat_logApi_8087/logs/logsApi/info.log.2017-09-15 -r '"userId":"[a-zA-Z0-9\-]+"'
```

# 2. 部署执行
### 基准测试--服务器小文件
```bash
nice --10 ./uniq-fields-count -f /opt/web_app/tomcat_baseApi_8088/logs/baseApi/info.log -r '"userId":"[a-zA-Z0-9\-]+"'
nice: cannot set niceness: Permission denied
2017/09/16 00:56:19 [提示]同时启用40个CPU执行任务
2017/09/16 00:56:19 [提示]文件名：/opt/web_app/tomcat_baseApi_8088/logs/baseApi/info.log
2017/09/16 00:56:19 [提示]当前文件大小：1.1661053GB
2017/09/16 00:56:19 [提示]匹配模式："userId":"[a-zA-Z0-9\-]+"
2017/09/16 00:56:26 [提示]搜索完成，文件共2562393行，找到有效数据199477个，任务耗时：7.424983828s
```
> 结论预测：在进程优先级能保证的情况下，100GB的文件时间在5分钟左右。

### 基准测试--服务器大文件
```bash
-sh-4.1$ ./uniq-fields-count -f /opt/web_app/tomcat_logApi_8087/logs/logsApi/info.log.2017-09-15 -r '"userId":"[a-zA-Z0-9\-]+"'   
2017/09/16 00:52:04 [提示]同时启用40个CPU执行任务
2017/09/16 00:52:04 [提示]文件名：/opt/web_app/tomcat_logApi_8087/logs/logsApi/info.log.2017-09-15
2017/09/16 00:52:04 [提示]当前文件大小：93.88157GB
2017/09/16 00:52:04 [提示]匹配模式："userId":"[a-zA-Z0-9\-]+"
2017/09/16 01:01:48 [提示]搜索完成，文件共181766328行，找到有效数据2006745个，任务耗时：9m44.194685769s
```
> 结论预测：在进程优先级能保证的情况下，100GB的文件时间在5分钟左右。

### 基准测试--本地小文件
```dos
PS E:\code\go\uniq-fields-count\bin> .\uniq-fields-count.exe -f E:\apps\apache-tomcat-7.0.57\logs\msgpushService\info.log -r "\""userId\"":\""[a-zA-Z0-9\-]+\"""
2017/09/16 01:07:31 [提示]同时启用8个CPU执行任务
2017/09/16 01:07:31 [提示]文件名：E:\apps\apache-tomcat-7.0.57\logs\msgpushService\info.log
2017/09/16 01:07:31 [提示]当前文件大小：1018.096MB
2017/09/16 01:07:31 [提示]匹配模式："userId":"[a-zA-Z0-9\-]+"
2017/09/16 01:07:31 [提示]搜索完成，文件共1728075行，找到有效数据0个，任务耗时：640.6774ms
```
> 结论预测：100GB的文件时间在3分钟左右。
