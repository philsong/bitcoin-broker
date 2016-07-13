## BTC trading market broker

    作为经纪商角色，提供thrift标准服务接口。

    提供合并交易所深度的标准价格模型，
    大并发订单内部自动撮合冲销，
    智能自动路由订单并拆大单为小单到不同交易所，
    失败订单自动重试处理。

    支持haobtc,okcoin,huobi,chbtc等交易所。

# 本地搭建 #

1、安装golang开发运行环境
	
	http://golang.org/doc/install

2、下载安装依赖库并编译 broker

	./install
	
3、导入数据库表结构

	导入etc目录下的*.sql文件到PostgresDB

4、修改配置

    修改conf/config_sample.json 为 conf/config.json

5、运行 broker

	./bin/broker

一切顺利的话，broker应该就启动了。
