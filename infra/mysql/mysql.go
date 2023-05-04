package mysql

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"my_stock_market/config"
)

type DBProvider struct {
	DSN string
}

func GetDBProvider(ctx context.Context) *DBProvider {
	return dbProvider
}

func (d *DBProvider) WithContext(ctx context.Context) *gorm.DB {
	return defaultDB.WithContext(ctx)
}

var dbProvider *DBProvider
var defaultDB *gorm.DB

func MustInitMySQL(ctx context.Context) error {
	mysqlConf := config.GetMySQLConf(ctx)

	dbProvider = &DBProvider{}

	dbProvider.DSN = fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		mysqlConf.Username, mysqlConf.Password, mysqlConf.Host, mysqlConf.Port, mysqlConf.Database)

	db, err := gorm.Open(mysql.Open(dbProvider.DSN), &gorm.Config{})
	if err != nil {
		logrus.Errorf("[MustInitMySQL] get db error: %v", err)
		return err
	}
	defaultDB = db

	err = testDBConnection(ctx)
	if err != nil {
		return err
	}

	logrus.Info("[MustInitMySQL] init mysql success")
	return nil
}

func testDBConnection(ctx context.Context) error {
	err := GetDBProvider(ctx).WithContext(ctx).Table("test").Select("name").Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.Errorf("[TestDBConnection] error: %v", err)
		return err
	}
	logrus.Info("[testDBConnection] connect mysql success")
	return nil
}
