package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"my_stock_market/cmd/fund"
	"my_stock_market/config"
	"my_stock_market/infra"
	"os"
	"time"
)

func main() {
	ctx := context.Background()

	var err error

	var cstSh, _ = time.LoadLocation("Asia/Shanghai")
	currentTime := time.Now().In(cstSh).Format("2006-01-02-15:04:05")

	logFile, err := os.OpenFile("logs/"+currentTime+".txt", os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	logrus.SetOutput(logFile)
	logrus.Info("服务开始启动")

	err = config.MustInitConf(ctx)
	if err != nil {
		panic(err)
	}

	err = infra.MustInitInfra(ctx)
	if err != nil {
		panic(err)
	}

	//s := tushare.NewStock(ctx)
	t := fund.NewFund(ctx)

	err = t.Check(ctx)
	if err != nil {
		panic(err)
	}
}
