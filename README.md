## BTC trading market broker for HaoBTC

# 本地搭建 #

1、安装golang开发运行环境
	
	http://golang.org/doc/install

2、下载安装依赖库并编译 broker

	./install
	
3、导入数据库表结构

	导入etc目录下的*.sql文件到PostgresDB

4、运行 broker

	./bin/broker

一切顺利的话，broker应该就启动了。
