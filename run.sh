go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct

cd /go/src/github.com/thesky341/my_stock_market

go run main.go