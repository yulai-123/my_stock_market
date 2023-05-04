package test

import (
	"context"
	"github.com/sirupsen/logrus"
	"my_stock_market/config"
	"my_stock_market/infra"
	"my_stock_market/repo/impl/daily"
	"os"
	"testing"
)

func initContainer(ctx context.Context) {
	logFile, err := os.OpenFile("./log.txt", os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	logrus.SetOutput(logFile)

	err = config.MustInitConf(ctx)
	if err != nil {
		panic(err)
	}

	err = infra.MustInitInfra(ctx)
	if err != nil {
		panic(err)
	}
}

func TestAny(t *testing.T) {
	ctx := context.Background()

	initContainer(ctx)

	dailyDAL := daily.NewStockDailyDAL(ctx)

}
