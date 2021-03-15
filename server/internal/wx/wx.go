package wx

import (
	"ginblog/model"
	"fmt"
	"github.com/esap/wechat"
) // 微信SDK包

var Wechat  *wechat.Server

func Init()  {

	fmt.Println("微信信息初始化")

	wechat.Debug = true

	cfg := &wechat.WxConfig{
		AppName:		"kms技术分享",
		Token:         	model.MapOpts.MustGet("wechat_token"),
		AppId:          model.MapOpts.MustGet("wechat_appid"),
		Secret:         model.MapOpts.MustGet("wechat_secret"),
		EncodingAESKey:	model.MapOpts.MustGet("wechat_aeskey"),
	}

	fmt.Println(cfg.Token, cfg.AppId,	cfg.Secret, 	cfg.EncodingAESKey)

	Wechat = wechat.New(cfg)
}
