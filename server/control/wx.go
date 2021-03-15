package control

import (
	"fmt"
	"ginblog/internal/wx"
	"ginblog/libs"
	"ginblog/process"
	"github.com/gin-gonic/gin"
	"time"
)

// OptsGet 获取某个配置项
func WxApiPost(ctx *gin.Context)  {
	//wx.Wechat.VerifyURL(ctx.Response().Writer, ctx.Request())
	//fmt.Println(ctx.Request().Form)

	wxCtx:= wx.Wechat.VerifyURL(ctx.Writer, ctx.Request)
	replyUser:=wxCtx.Msg.FromUserName
	//fmt.Println(replyUser)
	adminOpenId:=libs.Viper.GetString("admin.open_id")

	replyMsg := time.Now().String()
	if replyUser == adminOpenId {
		replyMsg = "管理员账号:adminOpenId:"+replyMsg
		go process.SendWxTimes()
	}
	err :=wxCtx.NewText(replyMsg).Reply()

	if err!=nil {
		fmt.Println(err)
	}
	
	return
}
