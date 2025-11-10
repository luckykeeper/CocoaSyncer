// CocoaSyncer - 心爱酱多节点智能解析平台 - 路由
//	@CreateTime		: 2023/8/30 15:10
//	@LastModified	: 2025/11/10 19:58
//	@Author			: Luckykeeper
//	@Email			: luckykeeper@luckykeeper.site
//	@Project		: CocoaSyncer

package cocoaSyncerRouter

import (
	"cocoaSyncer/handler"
	"cocoaSyncer/model"
	"cocoaSyncer/static"
	"cocoaSyncer/templates"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
)

var CocoaSyncerRouterDemo *gin.Engine

// 打包静态文件
// https://github.com/gin-gonic/examples/blob/master/assets-in-binary/example02/main.go
// 用注释告诉 embed 要打包哪些文件
// 嵌入父目录：https://www.ixiqin.com/2022/10/02/in-the-embedded-in-the-parent-directory-of-the-file/
// 嵌入父目录：https://github.com/golang/go/issues/46056
// 见 static 和 Templates 文件夹的 go 文件

func GinRouterDemo(router *gin.Engine, cocoaBasic model.CocoaBasic) {
	template := template.Must(template.New("").ParseFS(templates.TemplateFs, "*.tmpl"))
	router.SetHTMLTemplate(template)
	router.StaticFS("/static", http.FS(static.StaticFs))

	// 根路径
	router.Any("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{})
	})

	// logo 302
	router.Any("/favicon.ico", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/static/favicon.ico")
	})

	// ------------------------------------------
	// show - V1 - 路由组
	v1 := router.Group("/v1")

	// show - V1 - 根路径
	v1.Any("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "v1.tmpl", gin.H{})
	})

	// show - V1 - 节点名称卡片
	v1.GET("/badge", func(c *gin.Context) {
		handler.GetNodeNameBadgeDemo(c, cocoaBasic)
	})
}
