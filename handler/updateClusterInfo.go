// CocoaSyncer - 心爱酱多节点智能解析平台 - 集群状态信息
//	@CreateTime		: 2023/8/30 15:10
//	@LastModified	: 2023/9/19 17:06
//	@Author			: Luckykeeper
//	@Email			: luckykeeper@luckykeeper.site
//	@Project		: CocoaSyncer

package handler

import (
	"cocoaSyncer/model"
	subFunction "cocoaSyncer/subfunction"
	"net/http"

	"github.com/gin-gonic/gin"
)

// api - V1 - 加入集群/更新信息
// UpdateClusterInfo godoc
//
//	@Summary		加入集群/更新信息
//	@Description	【程序调用】用于 leader 节点将 follower 节点加入集群，由 leader 节点向 follower 节点推送数据
//	@Tags			业务
//	@Accept			json
//	@Produce		json
//	@Param			apiRequest	body		model.APIRequest	false	"API 请求数据，鉴权参数是 apiRequest.CocoaSecret ，节点仅接收 `otherCocoaSyncer`, `cocoaManagedService`, `cloudPlatformInfo`的数据，其中状态码和喂狗数据不接收"
//	@Success		200			{object}	model.ServerReturn
//	@Router			/updateClusterInfo [post]
func UpdateClusterInfo(context *gin.Context) {
	// 这个方法仅限 follower 类型的节点可被调用
	getCocoaBasic := &model.CocoaBasic{ConfigImported: true}
	var serverReturn model.ServerReturn
	subFunction.CocoaDataEngine.Get(getCocoaBasic)

	if subFunction.CheckExecRole() {
		serverReturn.StatusCode = 405
		serverReturn.StatusString = "Not Allowed! Cocoa is soooo kawaii, Do you still sure you really want to eat Cocoa???"
		context.JSON(http.StatusOK, serverReturn)
	} else {
		// 验证 Token
		var apiRequest model.APIRequest
		context.ShouldBind(&apiRequest)
		if !subFunction.CheckToken(apiRequest.CocoaSecret) {
			serverReturn.StatusCode = 401
			serverReturn.StatusString = "Unauthorized! You are not Cocoa!"
			context.JSON(http.StatusOK, serverReturn)
		} else {
			// 加入集群，写入/更新节点信息
			// 不更改任何一端的看门狗和状态码
			// 状态码由本地判断，喂狗数据由各端处理
			for _, info := range apiRequest.CocoaBasic.OtherCocoaSyncer {
				for _, infoLocal := range getCocoaBasic.OtherCocoaSyncer {
					if infoLocal.NodeCode == info.NodeCode {
						info.StatusCode = infoLocal.StatusCode
						info.WatchDog = infoLocal.WatchDog
					}
				}

				// // 本端
				// if info.NodeCode == getCocoaBasic.NodeCode {
				// 	// 不更改本端的看门狗和状态码
				// 	for _, infoLocal := range getCocoaBasic.OtherCocoaSyncer {
				// 		// 本端
				// 		if info.NodeCode == getCocoaBasic.NodeCode {
				// 			// 修改对端发送数据中看门狗和状态码为本端数据再写入
				// 			info.WatchDog = infoLocal.WatchDog
				// 			info.StatusCode = infoLocal.StatusCode
				// 		} else {
				// 			// 对端数据，不修改状态码
				// 			if info.NodeCode == infoLocal.NodeCode {
				// 				info.StatusCode = infoLocal.StatusCode
				// 			}
				// 		}
				// 	}
				// }

			}

			subFunction.CocoaDataEngine.Cols("otherCocoaSyncer", "cocoaManagedService", "cloudPlatformInfo").Update(apiRequest.CocoaBasic)
			serverReturn.StatusCode = 200
			serverReturn.StatusString = "MafuMafu~"
			context.JSON(http.StatusOK, serverReturn)
		}
	}
}
