// CocoaSyncer - 心爱酱多节点智能解析平台 - 获取节点名称（用于反馈是访客访问的是哪台节点）
//	@CreateTime		: 2023/8/30 15:10
//	@LastModified	: 2025/11/10 19:59
//	@Author			: Luckykeeper
//	@Email			: luckykeeper@luckykeeper.site
//	@Project		: CocoaSyncer

package handler

import (
	"cocoaSyncer/model"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetNodeNameBadgeDemo(context *gin.Context, basic model.CocoaBasic) {
	thisCocoa := &model.CocoaBasic{ConfigImported: true}

	// 屏蔽掉不应该在此的输出
	thisCocoa.DataBaseType = "*****"
	thisCocoa.Dsn = "*****"
	thisCocoa.CocoaSecret = "*****"
	thisCocoa.CloudPlatformInfo = nil
	// thisCocoa.OtherCocoaSyncer = nil
	thisCocoa.CocoaManagedService = nil

	thisCocoa.NodeName = basic.NodeName

	// 在线节点数量显示
	nodesTotalNum := 2
	nodeLiveNum := 2

	thisCocoa.NodeName = thisCocoa.NodeName + " | " + fmt.Sprint(nodeLiveNum) + "/" + fmt.Sprint(nodesTotalNum)

	// 转 map[string]interface{} 给前端
	thisCocoaJson, _ := json.Marshal(&thisCocoa)
	var thisCocoaMap map[string]interface{}
	_ = json.Unmarshal(thisCocoaJson, &thisCocoaMap)

	// tmplate 用法：https://www.liwenzhou.com/posts/Go/template/
	// 注意要想 img 能够显示，必须指定 Content-Type
	context.Header("Content-Type", "image/svg+xml;charset=utf-8")
	context.Header("Content-Disposition", "Badge of "+thisCocoa.NodeName)
	context.HTML(http.StatusOK, "nodeName.tmpl", thisCocoaMap)
	// context.HTML(http.StatusOK, "status.tmpl", gin.H{
	// 	"title": thisCocoa.NodeName,
	// })
	// context.JSON(http.StatusOK, thisCocoa)
}
