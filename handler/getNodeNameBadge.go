// CocoaSyncer - 心爱酱多节点智能解析平台 - 获取节点名称（用于反馈是访客访问的是哪台节点）
//	@CreateTime		: 2023/8/30 15:10
//	@LastModified	: 2023/9/19 17:06
//	@Author			: Luckykeeper
//	@Email			: luckykeeper@luckykeeper.site
//	@Project		: CocoaSyncer

package handler

import (
	"cocoaSyncer/model"
	subFunction "cocoaSyncer/subfunction"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// show - V1 - Badge
// GetNodeNameBadge godoc
//
//	@Summary		展示 CocoaSyncer 节点卡片，请使用浏览器打开链接，Swagger 无法调用
//	@Description	【浏览器打开，非程序调用】向用户展示 CocoaSyncer 节点卡片（Badge）
//	@Tags			监控
//	@Success		200
//	@Router			/badge [get]
func GetNodeNameBadge(context *gin.Context) {
	thisCocoa := &model.CocoaBasic{ConfigImported: true}
	// thisCocoa.ConfigImported = true
	subFunction.CocoaDataEngine.Where("configImported=true").Get(thisCocoa)

	// 屏蔽掉不应该在此的输出
	thisCocoa.DataBaseType = "*****"
	thisCocoa.Dsn = "*****"
	thisCocoa.CocoaSecret = "*****"
	thisCocoa.CloudPlatformInfo = nil
	thisCocoa.OtherCocoaSyncer = nil
	thisCocoa.CocoaManagedService = nil

	// 转 map[string]interface{} 给前端
	thisCocoaJson, _ := json.Marshal(&thisCocoa)
	var thisCocoaMap map[string]interface{}
	_ = json.Unmarshal(thisCocoaJson, &thisCocoaMap)

	// tmplate 用法：https://www.liwenzhou.com/posts/Go/template/
	context.HTML(http.StatusOK, "nodeName.tmpl", thisCocoaMap)
	// context.HTML(http.StatusOK, "status.tmpl", gin.H{
	// 	"title": thisCocoa.NodeName,
	// })
	// context.JSON(http.StatusOK, thisCocoa)
}
