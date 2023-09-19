// CocoaSyncer - 心爱酱多节点智能解析平台 - 路由
//	@CreateTime		: 2023/8/30 15:10
//	@LastModified	: 2023/9/19 17:06
//	@Author			: Luckykeeper
//	@Email			: luckykeeper@luckykeeper.site
//	@Project		: CocoaSyncer

package cocoaSyncerRouter

import (
	"cocoaSyncer/handler"
	"cocoaSyncer/model"
	"cocoaSyncer/static"
	subFunction "cocoaSyncer/subfunction"
	"cocoaSyncer/templates"
	"html/template"
	"net/http"
	"strings"

	"cocoaSyncer/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var CocoaSyncerRouter *gin.Engine

// 打包静态文件
// https://github.com/gin-gonic/examples/blob/master/assets-in-binary/example02/main.go
// 用注释告诉 embed 要打包哪些文件
// 嵌入父目录：https://www.ixiqin.com/2022/10/02/in-the-embedded-in-the-parent-directory-of-the-file/
// 嵌入父目录：https://github.com/golang/go/issues/46056
// 见 static 和 Templates 文件夹的 go 文件

func GinRouter(router *gin.Engine, devMode bool) {
	// LOGO
	// router.StaticFile("/favicon.ico", "./static/favicon.ico")
	// router.StaticFile("/static/cocoa", "./static/cocoa")
	// router.StaticFile("/static/ba", "./static/ba")
	// router.StaticFile("/static/sakura.png", "./static/sakura.png")
	// router.StaticFile("/static/t9kX2ZsbR71DFh0IscwBnjtRgJVuakEywhN",
	// 	"./static/t9kX2ZsbR71DFh0IscwBnjtRgJVuakEywhN")
	// router.LoadHTMLGlob("templates/*")
	// 改造常见问题：cannot use template (variable of type *"text/template".Template) as *"html/template".Template value in argument to router.SetHTMLTemplate
	// 见：https://github.com/gin-gonic/gin/issues/3216，移除 "text/template" 依赖

	thisCocoa := model.CocoaBasic{ConfigImported: true}
	subFunction.CocoaDataEngine.Get(&thisCocoa)

	// programatically set swagger info
	docs.SwaggerInfo.Title = "CocoaSyncer Swagger Doc Center"
	docs.SwaggerInfo.Description = "CocoaSyncer 在线文档 - V1"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = strings.Split(thisCocoa.NodeAddress, "://")[1]
	docs.SwaggerInfo.BasePath = "/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

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

	// Prometheus Exporter
	handler.PrometheusExporter(router)

	// ------------------------------------------
	// show - V1 - 路由组
	v1 := router.Group("/v1")

	// show - V1 - 根路径
	v1.Any("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "v1.tmpl", gin.H{})
	})

	if devMode {
		// Swagger - V1
		v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// show - V1 - 展示集群状态
	v1.GET("/status", handler.ShowStatus)

	// show - V1 - 节点名称卡片
	v1.GET("/badge", handler.GetNodeNameBadge)

	// api - V1 - 获取集群信息
	v1.POST("/getClusterInfo", handler.GetClusterInfo)

	// api - V1 - 加入集群&更新集群信息
	v1.POST("/updateClusterInfo", handler.UpdateClusterInfo)

}
