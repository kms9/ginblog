package router

import (
	"ginblog/conf"
	"ginblog/internal/jwt"
	"ginblog/model"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
	"github.com/gin-contrib/multitemplate"

	//"github.com/labstack/echo/v4"
	//"github.com/labstack/echo/v4/middleware"
	"ginblog/utils"
)

var funcMap template.FuncMap

func init() {
	funcMap = template.FuncMap{"str2html": Str2html, "str2js": Str2js, "date": Date, "md5": Md5}
}
func midRecover() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}
				stack := make([]byte, 1<<10)
				length := runtime.Stack(stack, false)
				// stdlog.Println(string(stack[:length]))
				os.Stdout.Write(stack[:length])
				//ctx.Error(err)
				_ = ctx.Error(err.(error)) // nolint: errcheck
				ctx.Abort()
				return
			}
		}()
		ctx.Next()
	}
}



//// 跨越配置
//var crosConfig = middleware.CORSConfig{
//	AllowOrigins: []string{"*"},
//	AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
//}

// TplRender is a custom html/template renderer for Echo framework
type TplRender struct {
	templates *template.Template
}

// Render renders a template document
func (t *TplRender) Render(w io.Writer, name string, data interface{}, ctx *gin.Context) error {
	// 获取数据配置项
	if mp, is := data.(map[string]interface{}); is {
		mp["appjs"] = AppJsUrl
		mp["appcss"] = AppCssUrl
		mp["title"] = model.MapOpts.MustGet("title")
		mp["favicon"] = model.MapOpts.MustGet("favicon")
		mp["comment"] = model.MapOpts.MustGet("comment")
		mp["analytic"] = model.MapOpts.MustGet("analytic")
		mp["site_url"] = model.MapOpts.MustGet("site_url")
		mp["logo_url"] = model.MapOpts.MustGet("logo_url")
		mp["keywords"] = model.MapOpts.MustGet("keywords")
		mp["miitbeian"] = model.MapOpts.MustGet("miitbeian")
		mp["weibo_url"] = model.MapOpts.MustGet("weibo_url")
		mp["custom_js"] = model.MapOpts.MustGet("custom_js")
		mp["github_url"] = model.MapOpts.MustGet("github_url")
		mp["description"] = model.MapOpts.MustGet("description")
	}
	//开发模式
	//每次强制读取模板
	//每次强制加载函数
	//if conf.App.IsDev() {
	//	t.templates = ctx. LoadTmpl("./views", funcMap)
	//}
	return t.templates.ExecuteTemplate(w, name, data)
}

// Str2html Convert string to template.HTML type.
func Str2html(raw string) template.HTML {
	return template.HTML(raw)
}

// Str2js Convert string to template.JS type.
func Str2js(raw string) template.JS {
	return template.JS(raw)
}

// Date Date
func Date(t time.Time, format string) string {
	return t.Format(format) //"2006-01-02 15:04:05"
}

// Md5 Md5
func Md5(str string) string {
	ctx := md5.New()
	ctx.Write([]byte(str))
	return hex.EncodeToString(ctx.Sum(nil))
}

// 初始化模板和函数
//func initRender() *TplRender {
//	tpl := utils.LoadTmpl("./views", funcMap)
//	return &TplRender{
//		templates: tpl,
//	}
//}

// midJwt 中间件-jwt验证
func midJwt() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenRaw := ctx.DefaultPostForm("token", "") // query/form 查找 token
		if tokenRaw == "" {
			tokenRaw = ctx.Request.Header.Get( "x-token") // header 查找token
			if tokenRaw == "" {
				ctx.JSON(utils.ErrJwt(`请重新登陆`, `未发现jwt`))
				return
			}
			tokenRaw = tokenRaw[7:] // Bearer token len("Bearer ")==7
		}
		jwtAuth, err := jwt.Verify(tokenRaw, conf.App.Jwtkey)
		if err == nil {
			ctx.Set("auth", jwtAuth)
			ctx.Set("uid", jwtAuth.Id)
		} else {
			ctx.JSON(utils.ErrJwt(`请重新登陆","jwt验证失败`))
			return
		}
		// 自定义头
		ctx.Next()
	}
}

const html404 = `<!DOCTYPE html><html lang="zh-CN"><head><meta charset="UTF-8"><title>404 Not Found zxysilent</title><style>* { margin: 0; padding: 0; }body { background-color: #f8f8f8; -webkit-font-smoothing: antialiased; }.error { position: absolute; left: 50%; top: 25rem; width: 483px; margin: -300px 0 0 -242px; padding-top: 199px; font-size: 18px; color: #666; text-align: center; background: #f8f8f8 url(/static/imgs/404.jpg) 0 0 no-repeat; }.error .remind { margin: 30px 0; }.error .button { display: inline-block; padding: 0 20px; line-height: 40px; font-size: 14px; color: #fff; background-color: #f8912d; text-decoration: none; }.error .button:hover { opacity: .9; }</style></head><body><div class="error">    <p class="remind">您访问的页面不存在，请返回主页！</p>    <p><a class="button" href="/">返回主页</a></p></div></body></html>`

func LoadTemplates(templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	// 非模板嵌套
	htmls, err := filepath.Glob(templatesDir + "/*.html")
	if err != nil {
		panic(err.Error())
	}
	for _, html := range htmls {
		r.AddFromGlob(filepath.Base(html), html)
	}


	// 嵌套的内容模板
	includes, err := filepath.Glob(templatesDir + "/*.html")
	if err != nil {
		panic(err.Error())
	}

	// template自定义函数
	funcMap := template.FuncMap{
		"StringToLower": func(str string) string {
			return strings.ToLower(str)
		},
	}

	// 将主模板，include页面，layout子模板组合成一个完整的html页面
	for _, include := range includes {
		files := []string{}
		files = append(files, templatesDir+"/frame.html", include)
		r.AddFromFilesFuncs(filepath.Base(include), funcMap, files...)
	}

	return r
}

func createMyRender() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	r.AddFromFiles("index", "templates/base.html", "templates/index.html")
	r.AddFromFiles("article", "templates/base.html", "templates/index.html", "templates/article.html")
	return r
}