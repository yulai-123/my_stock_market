#!/usr/bin/env sh

docker run -v /Users/thesky/go/src/github.com/thesky341/my_stock_market:/go/src/github.com/thesky341/my_stock_market \
  -t -i golang:latest  /bin/bash