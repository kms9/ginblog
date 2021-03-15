package control

import (
	"ginblog/model"
	"github.com/gin-gonic/gin"
	"strconv"
	"ginblog/utils"
)

// TagAll 所有标签
func TagAll(ctx *gin.Context)  {
	mods, err := model.TagAll()
	if err != nil {
		 ctx.JSON(utils.ErrOpt(`未查询到标签信息`, err.Error()))
		return
	}
	if len(mods) < 1 {
		 ctx.JSON(utils.ErrOpt(`未查询到标签信息`, "len"))
		return
	}
	 ctx.JSON(utils.Succ(`分类信息`, mods))
	return
}

// TagAdd 添加标签
func TagAdd(ctx *gin.Context)  {
	ipt := &model.Tag{}
	err := ctx.Bind(ipt)
	if err != nil {
		 ctx.JSON(utils.ErrIpt(`数据输入错误,请重试`, err.Error()))
		return
	}
	if !model.TagAdd(ipt) {
		 ctx.JSON(utils.Fail(`添加标签失败,请重试`))
		return
	}
	 ctx.JSON(utils.Succ(`添加标签成功`))
	return
}

// TagEdit 修改标签
func TagEdit(ctx *gin.Context)  {
	ipt := &model.Tag{}
	err := ctx.Bind(ipt)
	if err != nil {
		 ctx.JSON(utils.ErrIpt(`数据输入错误,请重试`, err.Error()))
		return
	}
	if !model.TagEdit(ipt) {
		 ctx.JSON(utils.Fail(`标签修改失败`))
		return
	}
	 ctx.JSON(utils.Succ(`标签修改成功`))
	return
}

// TagDrop  删除标签
func TagDrop(ctx *gin.Context)  {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		 ctx.JSON(utils.ErrIpt(`数据输入错误,请重试`, err.Error()))
		return
	}
	if !model.TagDrop(id) {
		 ctx.JSON(utils.Fail(`标签删除失败,请重试`))
		return
	}
	// 删除标签相关联的数据
	model.TagPostDrop(id)
	 ctx.JSON(utils.Succ(`标签删除成功`))
	return
}
