package control

import (
	"blog/internal/wx"
	"blog/libs"
	"blog/process"
	"github.com/labstack/echo/v4"
	"time"
)

// OptsGet 获取某个配置项
func WxApiPost(ctx echo.Context) error {
	//wx.Wechat.VerifyURL(ctx.Response().Writer, ctx.Request())
	//fmt.Println(ctx.Request().Form)

	wxCtx:= wx.Wechat.VerifyURL(ctx.Response().Writer, ctx.Request())
	replyUser:=wxCtx.Msg.FromUserName
	//fmt.Println(replyUser)
	adminOpenId:=libs.Viper.GetString("admin.open_id")

	replyMsg := time.Now().String()
	if replyUser == adminOpenId {
		replyMsg = "管理员账号:adminOpenId:"+replyMsg
		go process.SendWxTimes()
	}
	reply :=wxCtx.NewText(replyMsg).Reply()

	return reply
}
