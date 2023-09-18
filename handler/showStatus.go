// CocoaSyncer - 心爱酱多节点智能解析平台 - 本节点状态展示
// @CreateTime : 2023/8/30 15:10
// @LastModified : 2023/8/30 15:10
// @Author : Luckykeeper
// @Email : luckykeeper@luckykeeper.site
// @Project : CocoaSyncer

package handler

import (
	"cocoaSyncer/model"
	subFunction "cocoaSyncer/subfunction"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// show - V1 - 展示集群状态
func ShowStatus(context *gin.Context) {
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
	context.HTML(http.StatusOK, "status.tmpl", thisCocoaMap)
	// context.HTML(http.StatusOK, "status.tmpl", gin.H{
	// 	"title": thisCocoa.NodeName,
	// })
	// context.JSON(http.StatusOK, thisCocoa)
}
