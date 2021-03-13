package process

import (
	"blog/libs"
	"fmt"
	"strconv"
	"time"
	"blog/internal/wx"
)

//func SendWxTimes()  {
//
//}

func SendWxTimes() {
	timePer:=libs.Viper.GetInt("admin.msg_push_per")
	heartbeat := time.NewTicker(time.Duration(timePer) * time.Second)
	defer heartbeat.Stop()
	count 	 := 0
	maxCount := 3
	for {

		select {

		case <- heartbeat.C:
			adminOpenId := libs.Viper.GetString("admin.open_id")
			replyMsg := time.Now().String()
			replyMsg = "管理员账号:adminOpenId:"+replyMsg

			count++

			replyMsg = replyMsg+  "循环次数:"+strconv.Itoa(count)
			err:= wx.Wechat.SendText(adminOpenId, replyMsg)
			if err!=nil{
				fmt.Println(err)
			}

			if count>=maxCount{
				heartbeat.Stop()
				return
			}
			//
			////fmt.Println(replyUser)
			//adminOpenId:=libs.Viper.GetString("admin.open_id")

		}
	}
}