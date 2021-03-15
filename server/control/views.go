package control

import (
	"ginblog/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"ginblog/utils"
)

// IndexView 主页面
func IndexView(ctx *gin.Context) {
	//return ctx.HTML(200, `<html><head><meta charset="UTF-8"><title>文档</title></head><body><a href="/swagger/index.html">doc</a></body></html>`)
	pi, _ := strconv.Atoi(ctx.DefaultPostForm("page", "1"))
	if pi == 0 {
		pi = 1
	}
	ps, _ := atoi(model.MapOpts.MustGet("page_size"), 6)
	mods, _ := model.PostPage(pi, ps)
	total := model.PostCount()
	naver := model.Naver{}
	if pi > 1 {
		naver.Prev = "/?page=" + strconv.Itoa(pi-1)
	}
	if total > (pi * ps) {
		naver.Next = "/?page=" + strconv.Itoa(pi+1)
	}
	//ctx.HTML(http.StatusOK, render.Render())
	ctx.HTML(http.StatusOK, "index.html", map[string]interface{}{
		"Posts": mods,
		"Naver": naver,
	})
}

// ArchivesView 归档
func ArchivesView(ctx *gin.Context) {
	mods, err := model.PostArchive()
	if err != nil {
		ctx.Redirect(302, "/")
	}
	ctx.HTML(http.StatusOK, "archive.html", map[string]interface{}{
		"Posts": mods,
	})
	return
}
func ArchivesJson(ctx *gin.Context){
	mods, err := model.PostArchive()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "未查询到数据")
		return
	}
	ctx.JSON(utils.Succ("归档", mods))

}

// CatePostView 分类文章列表
func CatePostView(ctx *gin.Context)  {
	cate := ctx.Param("cate")
	if cate == "" {
		ctx.Redirect(302, "/")
		return
	}
	mod, has := model.CateName(cate)
	if !has {
		ctx.Redirect(302, "/")
	}
	pi, _ := strconv.Atoi(ctx.DefaultPostForm("page",  "1"))
	if pi == 0 {
		pi = 1
	}
	ps, _ := atoi(model.MapOpts.MustGet("page_size"), 6)
	mods, err := model.CatePostList(mod.Id, pi, ps, true)
	if err != nil || len(mods) < 1 {
		ctx.Redirect(302, "/")
		return
	}
	total := model.CatePostCount(mod.Id, true)
	naver := model.Naver{}
	if pi > 1 {
		naver.Prev = "/cate/" + mod.Name + "?page=1"
	}
	if total > (pi * ps) {
		naver.Next = "/cate/" + mod.Name + "?page=" + strconv.Itoa(pi+1)
	}
	ctx.HTML(http.StatusOK, "cate-post.html", map[string]interface{}{
		"Cate":      mod,
		"CatePosts": mods,
		"Naver":     naver,
	})
}

// TagsView 标签列表
func TagsView(ctx *gin.Context) {
	mods, err := model.TagStateAll()
	if err != nil {
		ctx.Redirect(302, "/")
		return
	}
	ctx.HTML(http.StatusOK, "tags.html", map[string]interface{}{
		"Tags": mods,
	})
}
func TagsJson(ctx *gin.Context) {
	mods, err := model.TagStateAll()
	if err != nil {
		ctx.JSON(utils.Fail("未查询到标签", err))
		return
	}
	ctx.JSON(utils.Succ("标签", mods))
	return
}

// TagPostView 标签文章列表
func TagPostView(ctx *gin.Context) {
	tag := ctx.Param("tag")
	if tag == "" {
		ctx.Redirect(302, "/tags")
		return
	}
	mod, has := model.TagName(tag)
	if !has {
		ctx.Redirect(302, "/tags")
		return
	}
	pi, _ := strconv.Atoi(ctx.DefaultPostForm("page", "1"))
	if pi == 0 {
		pi = 1
	}
	ps, _ := atoi(model.MapOpts.MustGet("page_size"), 6)
	mods, err := model.TagPostList(mod.Id, pi, ps)
	if err != nil || len(mods) < 1 {
		ctx.Redirect(302, "/tags")
		return
	}
	total := model.TagPostCount(mod.Id)
	naver := model.Naver{}
	if pi > 1 {
		naver.Prev = "/tag/" + mod.Name + "?page=1"
	}
	if total > (pi * ps) {
		naver.Next = "/tag/" + mod.Name + "?page=" + strconv.Itoa(pi+1)
	}
	ctx.HTML(http.StatusOK, "tag-post.html", map[string]interface{}{
		"Tag":      mod,
		"TagPosts": mods,
		"Naver":    naver,
	})
}

// PostView 文章页面
func PostView(ctx *gin.Context) {
	//return ctx.HTML(200, `<html><head><meta charset="UTF-8"><title>文档</title></head><body><a href="/swagger/index.html">doc</a></body></html>`)
	paths := strings.Split(ctx.Param("*"), ".")
	if len(paths) == 2 {
		mod, naver, has := model.PostPath(paths[0])
		if !has {
			ctx.Redirect(302, "/")
			return
		}
		if paths[1] == "html" {
			mod.Content = reg.ReplaceAllString(mod.Content, `<img class="lazy-load" src="data:image/gif;base64,R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw==" data-src="$1" alt="$2">`)
			tags, _ := model.PostTags(mod.Id)
			ctx.HTML(http.StatusOK, "post.html", map[string]interface{}{
				"Post":    mod,
				"Naver":   naver,
				"Tags":    tags,
				"HasTag":  len(tags) > 0,
				"HasCate": mod.Cate != nil,
			})
		}
		ctx.JSON(utils.Succ("", mod))
		return
	}
	ctx.Redirect(302, "/404")
	return
}

var reg = regexp.MustCompile(`<img src="([^" ]+)" alt="([^" ]*)"\s?\/?>`)

// 生成目录并替换内容
func getTocHTML(html string) string {
	html = strings.Replace(html, `id="`, `id="toc_`, -1)
	regToc := regexp.MustCompile("<h[1-6]>.*?</h[1-6]>")
	regH := regexp.MustCompile(`<h[1-6]><a id="(.*?)"></a>(.*?)</h[1-6]>`)
	hs := regToc.FindAllString(html, -1)
	if len(hs) > 1 {
		sb := strings.Builder{}
		sb.WriteString(`<div class="toc"><ul>`)
		level := 0
		for i := 0; i < len(hs)-1; i++ {
			fg := similar(hs[i], hs[i+1])
			if fg == 0 {
				sb.WriteString(regH.ReplaceAllString(hs[i], `<li><a href="#$1">$2</a></li>`))
			} else if fg == 1 {
				level++
				sb.WriteString(regH.ReplaceAllString(hs[i], `<li><a href="#$1">$2</a><ul>`))
			} else {
				level--
				sb.WriteString(regH.ReplaceAllString(hs[i], `<li><a href="#$1">$2</a></li></ul></li>`))
			}
		}
		fg := similar(hs[len(hs)-2], hs[len(hs)-1])
		if fg == 0 {
			sb.WriteString(regH.ReplaceAllString(hs[len(hs)-1], `<li><a href="#$1">$2</a></li>`))
		} else if fg == 1 {
			level++
			sb.WriteString(regH.ReplaceAllString(hs[len(hs)-1], `<li><a href="#$1">$2</a><ul>`))
		} else {
			level--
			sb.WriteString(regH.ReplaceAllString(hs[len(hs)-1], `<li><a href="#$1">$2</a></li></ul></li>`))
		}
		for level > 0 {
			sb.WriteString(`</ul></li>`)
			level--
		}
		sb.WriteString(`</ul></div>`)
		return sb.String() + html
	}
	if len(hs) == 1 {
		sb := strings.Builder{}
		sb.WriteString(`<div class="toc"><ul>`)
		sb.WriteString(regH.ReplaceAllString(hs[0], `<li><a href="#$1">$2</a></li>`))
		sb.WriteString(`</ul></div>`)
		return sb.String() + html
	}
	return html
}
func atoi(raw string, def int) (int, error) {
	out, err := strconv.Atoi(raw)
	if err != nil {
		return def, err
	}
	return out, nil
}
