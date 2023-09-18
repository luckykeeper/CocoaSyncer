// CocoaSyncer - 心爱酱多节点智能解析平台 - 集群状态信息
// @CreateTime : 2023/8/30 15:10
// @LastModified : 2023/8/30 15:10
// @Author : Luckykeeper
// @Email : luckykeeper@luckykeeper.site
// @Project : CocoaSyncer

package handler

import (
	"cocoaSyncer/model"
	subFunction "cocoaSyncer/subfunction"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// API - V1 - 获取集群信息
func GetClusterInfo(context *gin.Context) {
	var apiRequest model.APIRequest
	context.ShouldBind(&apiRequest)
	log.Println("API request:", apiRequest)
	var serverReturn model.ServerReturn
	// 验证 Token
	if !subFunction.CheckToken(apiRequest.CocoaSecret) {
		serverReturn.StatusCode = 401
		serverReturn.StatusString = "Unauthorized! You are not Cocoa!"
		context.JSON(http.StatusOK, serverReturn)
		return
	} else {
		// 拿节点数据
		getCocoaBasic := &model.CocoaBasic{ConfigImported: true}
		subFunction.CocoaDataEngine.Get(getCocoaBasic)
		serverReturn.CocoaBasic = *getCocoaBasic
		// 节点类型为 Leader ，业务码返回 200
		if serverReturn.CocoaBasic.WorkMode == "leader" {
			serverReturn.StatusCode = 200
			serverReturn.StatusString = "MafuMafu~"
			context.JSON(http.StatusOK, serverReturn)
			return
		} else if serverReturn.CocoaBasic.WorkMode == "follower" {
			// 节点类型为 follower ，且已经加入集群，业务码返回 200
			if len(serverReturn.CocoaBasic.OtherCocoaSyncer) > 1 {
				serverReturn.StatusCode = 200
				serverReturn.StatusString = "MafuMafu~"
				context.JSON(http.StatusOK, serverReturn)
				return
			} else if len(serverReturn.CocoaBasic.OtherCocoaSyncer) == 1 {
				// 节点类型为 follower ，且没有加入集群，业务码返回 201
				serverReturn.StatusCode = 201
				serverReturn.StatusString = "Mafu?"
				context.JSON(http.StatusOK, serverReturn)
				return
			}
		}
	}
}
