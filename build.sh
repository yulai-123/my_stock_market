#!/usr/bin/env sh

docker run --rm -v /Users/thesky/go/src/github.com/thesky341/my_stock_market:/go/src/github.com/thesky341/my_stock_market \
  -t -i golang:latest  /bin/bash