/*
  trader API Engine
*/

package main

import (
	"config"
	"logger"
	"math/rand"
	"runtime"
	"strategy"
	"time"
	"trade_server"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// 设置随机数种子
	rand.Seed(time.Now().Unix())
}

func main() {
	config.LoadConfig()

	printBanner()

	strategy.TradeCenter()

	thriftserver := config.Config["thriftserver"]
	server := trade_server.NewTraderServer(thriftserver)
	server.Run()
}

func printBanner() {
	version := "V1.1 postgres"
	logger.Infoln("[ ---------------------------------------------------------->>> ")
	logger.Infoln(" trading broker Engine", version)
	logger.Infoln(" <<<----------------------------------------------------------] ")
}
