package server

import (
	"fmt"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

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

// server/server.go
func NewHTTPServer(opts HTTPServerOptions) *http.Server {
	var lOpts = []http.ServerOption{
		http.Address(opts.Address),
		http.Middleware(
			recovery.Recovery(),
		),
		http.Middleware(opts.Middlewares...),
	}
	srv := http.NewServer(lOpts...)

	for _, r := range opts.Registrars {
		r.Register(srv)
	}

	return srv
}

type HTTPServerOptions struct {
	Address     string
	Middlewares []middleware.Middleware
	Registrars  []HTTPServiceRegistrar
}

type HTTPServiceRegistrar interface {
	Register(*http.Server)
}
