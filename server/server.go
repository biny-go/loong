package server

import "fmt"

func Start() {
	// wire注入如何处理？
	fmt.Println("reg wire server...")
	// 自动注册路由
	fmt.Println("reg router server...")
	// 自动向注册中心注册
	fmt.Println("reg center server...")
	// 启动 http
	// 启动 prc
	fmt.Println("Starting http rpc server...")
}
