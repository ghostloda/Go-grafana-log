Go语言学习~~
知识点：并发中Channel的使用
读取模块：
1、打开文件
2、从文件末尾开始逐行读取至Read Channel
解析模块：
1、从Read Channel中读取每行日志数据
2、利用正则匹配来提取所需要的监控数据（path、status、method..）
3、写入Write Channel
写入模块：
1、初始化数据库
2、将Write Channel写入数据库
展示：
1、grafana连接influxdb做图形展示
----
各个组件都是容器化的


