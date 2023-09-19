// CocoaSyncer - 心爱酱多节点智能解析平台 - 统一定时模块管理中心
//	@CreateTime		: 2023/8/30 15:10
//	@LastModified	: 2023/8/30 15:10
//	@Author			: Luckykeeper
//	@Email			: luckykeeper@luckykeeper.site
//	@Project		: CocoaSyncer

package subFunction

import (
	"log"

	"github.com/fatih/color"
	"github.com/robfig/cron"
)

// 在 Gin 跑定时任务的方法：https://www.cnblogs.com/alisleepy/p/17034621.html

// CronTaskCenter 定时器
var CronTaskCenter *cron.Cron

// 定时任务入口
// init 阶段初始化，main 包导入即可在 Init 阶段执行
// 手动调用也可
func CronTaskInit() {
	log.Println("Initializing CronCenter During Early Start!")
	CronTaskCenter = cron.New() // 定时任务
	CronTaskCenter.Start()
	// CronTaskCenter.AddFunc("@every 1m", func() { log.Println("1m") }) // 每隔1分钟执行一次的示例

	// 同步节点信息数据到前端
	CronTaskCenter.AddFunc("@every 10s", CronAddressData)

	// 更新看门狗状态
	CronTaskCenter.AddFunc("@every 30s", CronFeedDogAndWatchOtherCocoaDogStatus)

	// 定时推送数据到对端
	CronTaskCenter.AddFunc("@every 35s", CronUpdateClusterInfo)

	// 定时更新状态码
	CronTaskCenter.AddFunc("@every 40s", CronUpdateNodeStatus)

	// 定时更新 DNS 解析记录
	CronTaskCenter.AddFunc("@every 60s", CronUpdateDNS)

	color.Green("Initializing CronCenter Succeeded!")

}
