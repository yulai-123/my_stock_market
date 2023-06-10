# my_stock_market
1. 项目介绍 \
这是一个通过拉取tu_share数据，进行数据分析的项目 \
使用 Golang 开发，会将tu_share数据保存到本地，使用docker运行
2. 项目结构 
```shell
├── .gopkg
├── cmd # 运行命令
├── common
│   ├── util # 通用工具
├── config
│   ├── config.go
│   ├── myconf.yml # 配置文件
├── ddl # 数据库ddl
├── excel # 命令生成的excel
├── infra
│   ├── mysql # 数据库初始化
│   ├── infra.go
├── logs # 日志
├── model # 模型层
├── repo # 数据库dal
├── service # 服务层，目前只有拉取tushare数据的服务
├── test # 测试
├── .gitgnore
├── build.sh # 启动脚本
├── go.mod
├── main.go
├── README.md
├── run.sh
```
