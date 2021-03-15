package control

import (
	"ginblog/model"
	"github.com/gin-gonic/gin"
	"github.com/zxysilent/utils"
)

// OptsGet 获取某个配置项
func OptsGet(ctx *gin.Context)  {
	key := ctx.Param("key")
	if key == "" {
		 ctx.JSON(utils.ErrIpt(`请填写key值`))
		return
	}
	if val, ok := model.OptsGet(key); ok {
		 ctx.JSON(utils.Succ(``, val))
		return
	}
	 ctx.JSON(utils.ErrIpt(`错误的key值`))
	return
}

// OptsEdit 编辑某个配置项
func OptsEdit(ctx *gin.Context)  {
	ipt := &model.Opts{}
	err := ctx.Bind(ipt)
	if err != nil {
		 ctx.JSON(utils.ErrIpt(`数据输入错误,请重试`, err.Error()))
		return
	}
	if !model.OptsEdit(ipt) {
		 ctx.JSON(utils.Fail(`配置项修改失败`))
		return
	}
	 ctx.JSON(utils.Succ(`配置项修改成功`))
	return
}

// OptsBase 基本配置项目
func OptsBase(ctx *gin.Context)  {
	// ipt := &model.Opts{}
	// err := ctx.Bind(ipt)
	// if err != nil {
	// 	return ctx.JSON(utils.ErrIpt(`数据输入错误,请重试`, err.Error()))
	// }
	// if !model.OptsEdit(ipt) {
	// 	return ctx.JSON(utils.Fail(`配置项修改失败`))
	// }
	mp := model.MapOpts
	// delete(mp, "analytic")
	// delete(mp, "comment")
	 ctx.JSON(utils.Succ(`基本配置项目`, mp))
	return
}
