// CocoaSyncer - 心爱酱多节点智能解析平台 - CLI 管理服务信息
// @CreateTime : 2023/8/30 15:10
// @LastModified : 2023/8/30 15:10
// @Author : Luckykeeper
// @Email : luckykeeper@luckykeeper.site
// @Project : CocoaSyncer

package subFunction

import (
	"cocoaSyncer/model"
	"fmt"
	"os"

	"github.com/fatih/color"
)

// 管理服务（域名）信息
func CLIManageServices() {
	if !CheckExecRole() {
		color.Red("不允许在 follower 节点进行此操作，请前往 leader 节点执行！")
		os.Exit(0)
	}
	color.Red("更新服务信息 —— 请选择你要进行的操作（输入对应操作的数字即可）")
	color.Green("1.查询当前受管服务信息")
	color.Green("2.新增受管服务信息")
	color.Green("3.删除指定受管服务信息")
	color.Cyan("请输入：")
	var userChoice string
	fmt.Scan(&userChoice)
	switch userChoice {
	case "1":
		color.Cyan("你选择了查询当前受管服务信息")
		thisCocoa := &model.CocoaBasic{ConfigImported: true}
		CocoaDataEngine.Cols("cocoaManagedService").Where("configImported=true").Get(thisCocoa)
		for seq, info := range thisCocoa.CocoaManagedService {
			color.Green("————————————————————————————————")
			fmt.Println("顺序：", seq+1)
			fmt.Println("受管服务名称：", info.ServiceName)
			fmt.Println("服务域名:", info.Domain)
			fmt.Println("主机记录:", info.Record)
			fmt.Println("服务商:", info.PlatformProvider)
			color.Green("————————————————————————————————")
		}

		color.Green("查询成功!")

	case "2":
		color.Cyan("你选择了新增受管服务信息")
		var serviceName, domain, record, platformProvider string
		color.Cyan("请输入云平台服务商名称（仅输入下面对应服务商的字母）")
		color.Green("——————————当前支持服务商——————————")
		color.Green("Ali:阿里云")
		color.Green("————————————————————————————————")
		color.Cyan("请输入：")

		for {
			fmt.Scan(&platformProvider)
			if platformProvider == "Ali" {
				fmt.Println("服务商：", platformProvider)
				break
			} else {
				color.Red("输入错误请检查！")
				continue
			}
		}

		color.Cyan("请输入受管服务名称（助记，可随意输入）:")
		fmt.Scan(&serviceName)
		fmt.Println("受管服务名称:", serviceName)

		color.Cyan("请输入你的域名，如 “luckykeeper.site” :")
		fmt.Scan(&domain)
		fmt.Println("域名:", domain)

		color.Cyan("请输入主机记录，如 “www” :")
		fmt.Scan(&record)
		fmt.Println("主机记录:", record)

		var thisManageServiceInfo = &model.CocoaManagedService{ServiceName: serviceName,
			Domain:           domain,
			Record:           record,
			PlatformProvider: platformProvider}

		thisCocoa := &model.CocoaBasic{ConfigImported: true}
		CocoaDataEngine.Cols("cocoaManagedService", "otherCocoaSyncer").Where("configImported=true").Get(thisCocoa)
		thisCocoaModified := &model.CocoaBasic{ConfigImported: true}
		thisCocoaModified.CocoaManagedService = append(thisCocoa.CocoaManagedService, *thisManageServiceInfo)
		// 同步更新节点信息，注意添加顺序：节点/DNS服务商->服务（域名）/线路
		for _, info := range thisCocoa.OtherCocoaSyncer {
			info.CocoaManagedService = append(info.CocoaManagedService, *thisManageServiceInfo)
			thisCocoaModified.OtherCocoaSyncer = append(thisCocoaModified.OtherCocoaSyncer, info)
		}

		CocoaDataEngine.Cols("cocoaManagedService", "otherCocoaSyncer").Update(thisCocoaModified)

		color.Green("新增成功!")

	case "3":
		color.Cyan("你选择了删除当前受管服务信息")
		thisCocoa := &model.CocoaBasic{ConfigImported: true}
		CocoaDataEngine.Cols("cocoaManagedService", "otherCocoaSyncer").Where("configImported=true").Get(thisCocoa)
		for seq, info := range thisCocoa.CocoaManagedService {
			color.Green("————————————————————————————————")
			fmt.Println("顺序：", seq+1)
			fmt.Println("受管服务名称：", info.ServiceName)
			fmt.Println("服务域名:", info.Domain)
			fmt.Println("主机记录:", info.Record)
			fmt.Println("服务商:", info.PlatformProvider)
			color.Green("————————————————————————————————")
		}

		var domain, record string

		color.Cyan("请输入你的域名，如 “luckykeeper.site” :")
		fmt.Scan(&domain)
		fmt.Println("域名:", domain)

		color.Cyan("请输入主机记录，如 “www” :")
		fmt.Scan(&record)
		fmt.Println("主机记录:", record)

		thisCocoaModified := &model.CocoaBasic{ConfigImported: true}

		// 同步更新节点信息，注意添加顺序：节点/DNS服务商->服务（域名）/线路
		// 节点信息，通过遍历双循环删除
		for _, info := range thisCocoa.OtherCocoaSyncer {
			var serviceGroup []model.CocoaManagedService
			for _, services := range info.CocoaManagedService {
				if !(services.Domain == domain && services.Record == record) {
					serviceGroup = append(serviceGroup, services)

					// 受管服务信息，在此循环内一并修改
					thisCocoaModified.CocoaManagedService = append(thisCocoaModified.CocoaManagedService, services)
				}
			}
			info.CocoaManagedService = serviceGroup
			// 修改后的节点信息
			thisCocoaModified.OtherCocoaSyncer = append(thisCocoaModified.OtherCocoaSyncer, info)
		}

		CocoaDataEngine.Cols("cocoaManagedService", "otherCocoaSyncer").Update(thisCocoaModified)

		color.Green("删除成功!")

	default:
		color.Red("未输入或输入错误，你的输入是：", userChoice)
		os.Exit(1)
	}
}
