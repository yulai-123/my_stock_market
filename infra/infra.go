package infra

import (
	"context"
	"my_stock_market/infra/mysql"
)

func MustInitInfra(ctx context.Context) error {
	err := mysql.MustInitMySQL(ctx)
	if err != nil {
		return err
	}
	return nil
}
