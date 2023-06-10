#!/usr/bin/env sh

local_src_path=$(cd $(dirname $0);pwd)
docker_src_path="/go/src/github.com/thesky341/my_stock_market"

# 1. 运行后删除
# 2. 挂载本地目录到容器中，输出文件到本地
# 3. 挂载本地目录到容器中，保存.gopkg，以便下次编译时不用重新下载
# 4. 连接容器中的mysql，使用本地数据库，可以不需要
# 5. 启动run.sh
docker run --rm \
  -v ${local_src_path}:${docker_src_path} \
  -v ${local_src_path}/.gopkg:/go/pkg \
  --link stock_market_data \
  -t -i golang:latest /bin/bash ${docker_src_path}/run.sh