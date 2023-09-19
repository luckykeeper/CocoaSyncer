// CocoaSyncer - 心爱酱多节点智能解析平台 - 节点身份认证
//	@CreateTime		: 2023/8/30 15:10
//	@LastModified	: 2023/8/30 15:10
//	@Author			: Luckykeeper
//	@Email			: luckykeeper@luckykeeper.site
//	@Project		: CocoaSyncer

package subFunction

import "cocoaSyncer/model"

// 判断节点身份是否为 leader
func CheckExecRole() bool {
	thisCocoa := &model.CocoaBasic{ConfigImported: true}
	CocoaDataEngine.Where("configImported=true").Get(thisCocoa)
	if thisCocoa.WorkMode == "leader" {
		return true
	} else {
		return false
	}
}
