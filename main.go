package main

import (
	"flag"
	"fmt"
	"httpPressureTest/model"
	"httpPressureTest/server"
	"runtime"
)

type array []string

func (a *array) String() string {
	return fmt.Sprint(*a)
}

func (a *array) Set(s string) error {
	*a = append(*a, s)

	return nil
}

func main() {

	runtime.GOMAXPROCS(1)

	var (
		concurrency uint64 // 并发数
		totalNumber uint64 // 请求数(单个并发/协程)
		requestUrl  string // 压测的url 目前支持，http/https ws/wss
		headers     array  // 自定义头信息传递给服务器
		body        string // HTTP POST方式传送数据
	)

	flag.Uint64Var(&concurrency, "c", 1, "并发数")
	flag.Uint64Var(&totalNumber, "n", 1, "请求数(单个并发/协程)")
	flag.StringVar(&requestUrl, "u", "", "压测地址")
	flag.Var(&headers, "H", "自定义头信息传递给服务器 示例:-H 'Content-Type: application/json'")
	flag.StringVar(&body, "data", "", "HTTP POST方式传送数据")

	// 解析参数
	flag.Parse()
	if concurrency == 0 || totalNumber == 0 || (requestUrl == "") {
		fmt.Printf("示例: go run main.go -c 1 -n 1 -u https://www.baidu.com/ \n")
		fmt.Printf("当前请求参数: -c %d -n %d -u %s \n", concurrency, totalNumber, requestUrl)

		flag.Usage()

		return
	}

	request, err := model.NewRequest(requestUrl, headers, body)
	if err != nil {
		fmt.Printf("参数不合法 %v \n", err)

		return
	}

	fmt.Printf("\n 开始启动  并发数:%d 请求数:%d 请求参数: \n", concurrency, totalNumber)
	request.Print()

	// 开始处理
	server.Dispose(concurrency, totalNumber, request)

	return
}
