// CocoaSyncer - 心爱酱多节点智能解析平台 - 检查 Token
//	@CreateTime		: 2023/8/30 15:10
//	@LastModified	: 2023/8/30 15:10
//	@Author			: Luckykeeper
//	@Email			: luckykeeper@luckykeeper.site
//	@Project		: CocoaSyncer

package subFunction

import (
	"cocoaSyncer/model"
	"log"
)

// 检查 Token
func CheckToken(token string) bool {
	log.Println("Checking Token:", token)
	cocoaBasic := new(model.CocoaBasic)
	CocoaDataEngine.Cols("cocoaSecret").Get(cocoaBasic)
	if token == cocoaBasic.CocoaSecret {
		return true
	} else {
		log.Println("Checking Token Failed")
		return false
	}
}
