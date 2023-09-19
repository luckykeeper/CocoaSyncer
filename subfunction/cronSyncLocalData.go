// CocoaSyncer - 心爱酱多节点智能解析平台 - 定时同步 NodeInfo 到本机（仅同步访问地址，供 web 展示）
//	@CreateTime		: 2023/8/30 15:10
//	@LastModified	: 2023/8/30 15:10
//	@Author			: Luckykeeper
//	@Email			: luckykeeper@luckykeeper.site
//	@Project		: CocoaSyncer

package subFunction

import (
	"cocoaSyncer/model"
	"log"

	"github.com/fatih/color"
)

// 定时同步 OtherCocoaSyncer 的地址数据到前端显示
func CronAddressData() {
	color.Cyan("___________________________")
	log.Println("CronTaskCenter - CronAddressData - Start")
	log.Println("定时任务中心 - 同步地址数据到前端 - 开始")

	thisCocoa := &model.CocoaBasic{ConfigImported: true}
	CocoaDataEngine.Where("configImported=true").Get(thisCocoa)
	for _, info := range thisCocoa.OtherCocoaSyncer {
		if thisCocoa.NodeCode == info.NodeCode {
			thisCocoa.NodeAddress = info.NodeAddress
		}
	}
	CocoaDataEngine.Cols("nodeAddress").Update(thisCocoa)

	color.Green("CronTaskCenter - CronAddressData - Done")
	color.Green("定时任务中心 - 同步地址数据到前端 - 完成")
	color.Cyan("___________________________")
}
