## BTC trading market maker for HaoBTC

# 本地搭建 #

1、安装golang开发运行环境
	
	http://golang.org/doc/install

2、下载安装依赖库并编译 trader

	./install
	
这样便编译好了 trader

3、导入数据库表结构

	导入etc目录下的*.sql文件到PostgresDB

4、运行 trader

	// linux/mac或者Git Bash 下执行
	./bin/trader

一切顺利的话，trader应该就启动了。
