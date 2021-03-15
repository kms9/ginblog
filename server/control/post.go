package control

import (
	"ginblog/model"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"time"
	"ginblog/utils"
)

// PostGet 一个
// id int
func PostGet(ctx *gin.Context)  {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		 ctx.JSON(utils.ErrIpt(`数据输入错误,请重试`, err.Error()))
		return
	}
	mod, has := model.PostGet(id)
	if !has {
		 ctx.JSON(utils.ErrOpt(`未查询信息`))
		return
	}
	 ctx.JSON(utils.Succ(`信息`, mod))
	return
}

// PostPageAll 页面列表
func PostPageAll(ctx *gin.Context)  {
	mods, err := model.PostPageAll()
	if err != nil {
		 ctx.JSON(utils.ErrOpt(`未查询到页面信息`, err.Error()))
		return
	}
	if len(mods) < 1 {
		 ctx.JSON(utils.ErrOpt(`未查询到页面信息`, "len"))
		return
	}
	 ctx.JSON(utils.Succ(`页面信息`, mods))
	return
}

// PostTagGet 通过文章id 获取 标签ids
func PostTagGet(ctx *gin.Context)  {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		 ctx.JSON(utils.ErrIpt(`数据输入错误,请重试`, err.Error()))
		return
	}
	mods := model.PostTagGet(id)
	if mods == nil {
		 ctx.JSON(utils.ErrOpt(`未查询到标签信息`))
		return
	}
	 ctx.JSON(utils.Succ(`标签ids`, mods))
	return
}

// PostOpts 文章操作
func PostOpts(ctx *gin.Context)  {
	ipt := &struct {
		Post model.Post `json:"post" form:"post"` // 文章信息
		Type int        `json:"type" form:"type"` // 0 文章 1 页面
		Tags []int      `json:"tags" form:"tags"` // 标签
		Edit bool       `json:"edit" form:"edit"` // 是否编辑
	}{}
	err := ctx.Bind(ipt)
	if err != nil {
		 ctx.JSON(utils.ErrIpt(`数据输入错误,请重试`, err.Error()))
		return
	}
	if !ipt.Edit && model.PostExist(ipt.Post.Path) {
		 ctx.JSON(utils.ErrIpt(`当前访问路径已经存在,请重新输入`))
		return
	}
	// 同步类型
	ipt.Post.Type = ipt.Type
	if strings.Contains(ipt.Post.Content, "<!--more-->") {
		ipt.Post.Summary = strings.Split(ipt.Post.Content, "<!--more-->")[0]
	}
	// 生成目录
	if ipt.Type == 0 {
		ipt.Post.Content = getTocHTML(ipt.Post.Content)
	}
	// 编辑 文章/页面
	if ipt.Edit {
		// 修改日期在发布日期之前
		if ipt.Post.UpdateTime.Before(ipt.Post.CreateTime) {
			// 修改时间再发布时间后1分钟
			ipt.Post.UpdateTime = ipt.Post.CreateTime.Add(time.Minute * 2)
		}
		if model.PostEdit(&ipt.Post) {
			if ipt.Type == 0 {
				// 处理变动标签
				old := model.PostTagGet(ipt.Post.Id)
				new := ipt.Tags
				add := make([]int, 0)
				del := make([]int, 0)
				for _, itm := range old {
					if !inOf(itm, new) {
						del = append(del, itm)
					}
				}
				for _, itm := range new {
					if !inOf(itm, old) {
						add = append(add, itm)
					}
				}
				tagAdds := make([]model.PostTag, 0, len(add))
				for _, itm := range add {
					tagAdds = append(tagAdds, model.PostTag{
						TagId:  itm,
						PostId: ipt.Post.Id,
					})
				}
				// 删除标签
				model.PostTagDrops(ipt.Post.Id, del)
				// 添加标签
				model.TagPostAdds(&tagAdds)
				 ctx.JSON(utils.Succ(`文章修改成功`))
				return
			}
			 ctx.JSON(utils.Succ(`页面修改成功`))
			return
		}
		if ipt.Type == 0 {
			 ctx.JSON(utils.Fail(`文章修改失败,请重试`))
			return
		}
		 ctx.JSON(utils.Fail(`页面修改失败,请重试`))
		return
	}
	// 添加 文章/页面
	ipt.Post.UpdateTime = time.Now()
	if model.PostAdd(&ipt.Post) {
		// 添加标签
		// 文章
		if ipt.Type == 0 {
			//添加标签
			tagPosts := make([]model.PostTag, 0, len(ipt.Tags))
			for _, itm := range ipt.Tags {
				tagPosts = append(tagPosts, model.PostTag{
					TagId:  itm,
					PostId: ipt.Post.Id,
				})
			}
			model.TagPostAdds(&tagPosts)
			 ctx.JSON(utils.Succ(`文章添加成功`))
			return
		}
		 ctx.JSON(utils.Succ(`页面添加成功`))
		return
	}
	if ipt.Type == 0 {
		 ctx.JSON(utils.Fail(`文章添加失败,请重试`))
		return
	}
	 ctx.JSON(utils.Fail(`页面添加失败,请重试`))
	return
}
func similar(a, b string) int {
	if a[:4] == b[:4] {
		return 0
	}
	if a[:4] < b[:4] {
		return 1
	}
	return -1
}

// PostDrop  删除
func PostDrop(ctx *gin.Context)  {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		 ctx.JSON(utils.ErrIpt(`数据输入错误,请重试`, err.Error()))
		return
	}
	if !model.PostDrop(id) {
		 ctx.JSON(utils.Fail(`删除失败,请重试`))
		return
	}
	// 删除 文章对应的标签信息
	model.PostTagDrop(id)
	 ctx.JSON(utils.Succ(`删除成功`))
	return
}
func inOf(goal int, arr []int) bool {
	for idx := range arr {
		if goal == arr[idx] {
			return true
		}
	}
	return false
}
