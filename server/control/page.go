package control

import (
	"ginblog/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AboutView 关于
func AboutView(ctx *gin.Context)  {
	mod, has := model.PostSingle("about")
	if !has {
		ctx.Redirect(302, "/")
		return
	}
	mod.Content = reg.ReplaceAllString(mod.Content, `<img class="lazy-load" src="data:image/gif;base64,R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw==" data-src="$1" alt="$2">`)
	ctx.HTML(http.StatusOK, "page.html", map[string]interface{}{
		"Page": mod,
		"Show": mod.IsPublic && mod.Status == 3,
	})
}

// LinksView 友链
func LinksView(ctx *gin.Context)  {
	mod, has := model.PostSingle("links")
	if !has {
		ctx.Redirect(302, "/")
		return
	}
	mod.Content = reg.ReplaceAllString(mod.Content, `<img class="lazy-load" src="data:image/gif;base64,R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw==" data-src="$1" alt="$2">`)
	ctx.HTML(http.StatusOK, "page.html", map[string]interface{}{
		"Page": mod,
		"Show": mod.IsPublic && mod.Status == 3,
	})
}

// PageView 页面
func PageView(ctx *gin.Context)  {
	mod, has := model.PostSingle(ctx.Param("*"))
	if !has {
		ctx.Redirect(302, "/")
		return
	}
	mod.Content = reg.ReplaceAllString(mod.Content, `<img class="lazy-load" src="data:image/gif;base64,R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw==" data-src="$1" alt="$2">`)
	ctx.HTML(http.StatusOK, "page.html", map[string]interface{}{
		"Page": mod,
		"Show": mod.IsPublic && mod.Status == 3,
	})
}
