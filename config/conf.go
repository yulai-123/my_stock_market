package config

import (
	"context"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"sync"
)

type Conf struct {
	MySQL   *MySQLConf   `yaml:"mySQL"`
	TuShare *TuShareConf `yaml:"tuShare"`
}

type TuShareConf struct {
	Host  string `yaml:"host"`
	Token string `yaml:"token"`
}

type MySQLConf struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
}

var confOnce sync.Once
var myConf Conf

func MustInitConf(ctx context.Context) error {
	var outErr error
	confOnce.Do(func() {
		file, err := ioutil.ReadFile("config/myconf.yml")
		if err != nil {
			outErr = err
			logrus.Errorf("[MustInitConf] error: %v", err)
			return
		}
		err = yaml.Unmarshal(file, &myConf)
		if err != nil {
			outErr = err
			logrus.Errorf("[MustInitConf] error: %v", err)
			return
		}
	})

	if outErr != nil {
		return outErr
	}

	logrus.Info("[MustInitConf] init conf success")
	return nil
}

func GetMySQLConf(ctx context.Context) *MySQLConf {
	return myConf.MySQL
}

func GetTuShareConf(ctx context.Context) *TuShareConf {
	return myConf.TuShare
}
