package router

import (
	"ginblog/control"
	"github.com/gin-gonic/gin"
)

// wx相关 通用访问
func wxRouter(r  gin.IRouter) {
	//	api.Any(`/echo`, control.WxApiPost)   // 微信下发验证
	router := r.Group("/wx")
	// test 接口
	router.GET("/test",  control.WxApiPost)

}
