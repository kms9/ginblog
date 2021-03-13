package router

import (
	"blog/control"

	"github.com/labstack/echo/v4"
)

// wx相关 通用访问
func wxRouter(api *echo.Group) {
	api.Any(`/echo`, control.WxApiPost)   // 微信下发验证
}
