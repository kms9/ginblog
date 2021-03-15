package main

import (
	"fmt"
	"ginblog/cache"
	"ginblog/conf"
	"ginblog/model"
	"ginblog/router"
	"ginblog/services"
	"ginblog/setup"
	"github.com/kms9/publicyc"
	"github.com/kms9/publicyc/pkg/server/ogin"
)

// Engine ..
type Engine struct {
	yc.Application
}

// NewEngine 初始化相关方法
func NewEngine() *Engine {
	eng := &Engine{}
	if err := eng.Start(
		setup.StartLogger,
		setup.StartRedis,
		//setup.StartDB,
		setup.StartMysqlDB,
		//setup.StartRequest,

		cache.Init,
		services.ServiceInit,
		eng.serveHTTP,
	); err != nil {
		//onion_log.Panicf("Engine: %s", err)
		fmt.Errorf("Engine: %s", err)
	}
	return eng
}

// serveHTTP 启动http
func (eng *Engine) serveHTTP() error {
	server := ogin.UseConfig("http").Build()
	router.StartHttp(server)
	return eng.Serve(server)
}

func main() {

	conf.Init()
	model.Init()

	eng := NewEngine()
	if err := eng.Run(); err != nil {
		fmt.Println(err.Error())
	}
}
