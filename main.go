package main

import (
	"context"
	"fmt"
	"my_stock_market/config"
)

func main() {
	fmt.Println("hello world!")

	ctx := context.Background()

	config.MustInitConf(ctx)

}
