package router

import (
	"ginblog/control"
	"ginblog/middleware"

	//"github.com/gin-gonic/gin"
	"github.com/kms9/publicyc/pkg/conf"
	"github.com/kms9/publicyc/pkg/server/ogin"
	"html/template"
)

func StartHttp(server  *ogin.Server) {

	//engine := echo.New()
	//engine.Renderer = initRender()                    // 初始渲染引擎
	//server.HTMLRender = initRender()
	server.Use(middleware.GinRecovery())                 // 恢复 日志记录

	server.Use(middleware.Cors()) // 跨域设置
	//server.HideBanner = true                          // 不显示横幅
	//server.HideBanner = true                          // 不显示横幅
	//server.HTTPErrorHandler = HTTPErrorHandler        // 自定义错误处理
	//RegDocs(engine)                                   // 注册文档
	server.Static(`/dist`, "dist")                    // 静态目录 - 后端专用
	server.Static(`/static`, "static")                // 静态目录


	server.SetFuncMap(template.FuncMap{"str2html": Str2html, "str2js": Str2js, "date": Date, "md5": Md5})
	//server.HTMLRender = LoadTemplates("./views")
	server.LoadHTMLGlob("./views/*")

	server.StaticFile(`/favicon.ico`, "favicon.ico")        // ico
	server.StaticFile("/dashboard", "dist/index.html")     // 前后端分离页面

	//--- 页面 -- start
	server.GET(`/`, control.IndexView)                 // 首页
	server.GET(`/archives`, control.ArchivesView)      // 归档
	server.GET(`/archives.json`, control.ArchivesJson) // 归档 json
	server.GET(`/tags`, control.TagsView)              // 标签
	server.GET(`/tags.json`, control.TagsJson)         // 标签 json
	server.GET(`/tag/:tag`, control.TagPostView)       // 具体某个标签
	server.GET(`/cate/:cate`, control.CatePostView)    // 分类
	server.GET(`/about`, control.AboutView)            // 关于
	server.GET(`/links`, control.LinksView)            // 友链
	//server.GET(`/post/*`, control.PostView)            // 具体某个文章
	//server.GET(`/page/*`, control.PageView)            // 具体某个页面
	//--- 页面 -- end

	api := server.Group("/api")         // api/
	apiRouter(api)                      // 注册分组路由

	adm := server.Group("/adm") // adm/ 需要登陆才能访问
	adm.Use(midJwt())
	admRouter(adm)                      // 注册分组路由
	//wx := engine.Group("/wx")
	//wxRouter(wx)
	//err := engine.Start(conf.App.Addr)
	//if err != nil {
	//	log.Fatalln("run error :", err)
	//}

	bizRouter := server.Group("/"+conf.Detail().GetString("yc.server.http.name"))
	wxRouter(bizRouter)
}
