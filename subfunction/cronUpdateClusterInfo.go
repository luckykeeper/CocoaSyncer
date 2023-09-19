// CocoaSyncer - 心爱酱多节点智能解析平台 - 定时更新节点状态
//	@CreateTime		: 2023/8/30 15:10
//	@LastModified	: 2023/8/30 15:10
//	@Author			: Luckykeeper
//	@Email			: luckykeeper@luckykeeper.site
//	@Project		: CocoaSyncer

package subFunction

import (
	"cocoaSyncer/model"
	"log"
	"time"

	"github.com/fatih/color"
)

// 更新看门狗状态
// 1、喂本端看门狗
// 2、拉对端看门狗数据
func CronFeedDogAndWatchOtherCocoaDogStatus() {
	color.Cyan("___________________________")
	log.Println("CronTaskCenter - CronFeedDogAndWatchOtherCocoaDogStatus - Start")
	log.Println("定时任务中心 - 更新看门狗状态 - 开始")

	// 喂本端狗
	thisCocoa := &model.CocoaBasic{ConfigImported: true}
	CocoaDataEngine.Where("configImported=true").Get(thisCocoa)

	thisCocoaModified := &model.CocoaBasic{ConfigImported: true}

	for _, info := range thisCocoa.OtherCocoaSyncer {
		// 每个节点仅喂自己的狗
		if thisCocoa.NodeCode == info.NodeCode {
			info.WatchDog = time.Now()
		}
		thisCocoaModified.OtherCocoaSyncer = append(thisCocoaModified.OtherCocoaSyncer, info)
	}

	CocoaDataEngine.Cols("otherCocoaSyncer").Update(thisCocoaModified)
	color.Green("CronTaskCenter - CronFeedDogAndWatchOtherCocoaDogStatus - 1st/2")
	color.Green("定时任务中心 - 更新看门狗状态 - 完成 1st/2 - 本端喂狗")

	// 拉对端狗数据
	for _, info := range thisCocoaModified.OtherCocoaSyncer {
		// 这个操作不对本端生效
		if info.NodeCode != thisCocoa.NodeCode {
			log.Println("CronTaskCenter - CronFeedDogAndWatchOtherCocoaDogStatus - 拉对端狗数据 - ", info.NodeCode)
			data, success := pullClusterInfo(info)
			if success {
				// 仅更新对应节点的数据
				thisCocoaStage2 := &model.CocoaBasic{ConfigImported: true}
				// 这里应该遍历 thisCocoaModified 而不是 thisCocoa ，否则旧的时间就会把新写上去的时间覆盖
				for _, info := range thisCocoaModified.OtherCocoaSyncer {
					if info.NodeCode == data.NodeCode {
						// 抓取的节点名称与对端一致，拉取对端狗的数据（仅拉喂狗数据，不拉其它数据）
						info.WatchDog = data.WatchDog
						thisCocoaStage2.OtherCocoaSyncer = append(thisCocoaStage2.OtherCocoaSyncer, info)
					} else {
						// 不是拉取节点的数据，不更新
						thisCocoaStage2.OtherCocoaSyncer = append(thisCocoaStage2.OtherCocoaSyncer, info)
					}
				}
				CocoaDataEngine.Cols("otherCocoaSyncer").Update(thisCocoaStage2)

				color.Green("CronTaskCenter - CronFeedDogAndWatchOtherCocoaDogStatus - Done 2/2")
				color.Green("定时任务中心 - 更新看门狗状态 - 完成 2/2")
				color.Cyan("___________________________")
			} else {
				color.Red("CronTaskCenter - CronFeedDogAndWatchOtherCocoaDogStatus - Failed 2nd/2")
				color.Red("定时任务中心 - 更新看门狗状态 - 失败 2nd/2")
				color.Cyan("___________________________")
			}
		}
	}

}

// 更新节点状态码
// 根据狗的情况判断节点是否在线
func CronUpdateNodeStatus() {
	color.Cyan("___________________________")
	log.Println("CronTaskCenter - CronUpdateNodeStatus - Start")
	log.Println("定时任务中心 - 更新节点状态码 - 开始")

	thisCocoa := &model.CocoaBasic{ConfigImported: true}
	CocoaDataEngine.Where("configImported=true").Get(thisCocoa)

	thisCocoaModified := &model.CocoaBasic{ConfigImported: true}

	for _, info := range thisCocoa.OtherCocoaSyncer {
		// 超时时间，按照 config 中的 TimeOut 执行
		if time.Since(info.WatchDog) > time.Duration(thisCocoa.Timeout)*60*time.Second {
			info.StatusCode = 504
		} else {
			info.StatusCode = 200
		}
		thisCocoaModified.OtherCocoaSyncer = append(thisCocoaModified.OtherCocoaSyncer, info)
	}
	CocoaDataEngine.Cols("otherCocoaSyncer").Update(thisCocoaModified)

	color.Green("CronTaskCenter - CronUpdateNodeStatus - Done")
	color.Green("定时任务中心 - 更新节点状态码 - 完成")
	color.Cyan("___________________________")
}

// 同步数据
// 仅 leader 节点执行此方法
// 调用 /getClusterInfo 接口更新节点的数据
// 不会覆盖看门狗和状态码
func CronUpdateClusterInfo() {
	if CheckExecRole() {
		color.Cyan("___________________________")
		log.Println("CronTaskCenter - CronUpdateClusterInfo - Start")
		log.Println("定时任务中心 - 同步节点数据 - 开始")

		thisCocoa := &model.CocoaBasic{ConfigImported: true}
		CocoaDataEngine.Where("configImported=true").Get(thisCocoa)

		for _, info := range thisCocoa.OtherCocoaSyncer {
			// 这个操作不对本端生效
			if thisCocoa.NodeCode != info.NodeCode {
				if pushClusterInfo(info, *thisCocoa) {
					color.Green("CronTaskCenter - CronUpdateClusterInfo - Done")
					color.Green("定时任务中心 - 同步节点数据 - 完成")
					color.Cyan("___________________________")
				} else {
					color.Red("CronTaskCenter - CronUpdateClusterInfo - Failed")
					color.Red("定时任务中心 - 同步节点数据 - 失败")
					color.Cyan("___________________________")
				}
			}
		}
	}
}
