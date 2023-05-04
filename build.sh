#!/usr/bin/env sh

local_src_path="/Users/thesky/go/src/github.com/thesky341/my_stock_market"
docker_src_path="/go/src/github.com/thesky341/my_stock_market"

docker run --rm \
  -v ${local_src_path}:${docker_src_path} \
  -v ${local_src_path}/.gopkg:/go/pkg \
  --link stock_market_data \
  -t -i golang:latest /bin/bash ${docker_src_path}/run.sh