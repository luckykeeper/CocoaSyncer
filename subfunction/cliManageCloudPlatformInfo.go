// CocoaSyncer - 心爱酱多节点智能解析平台 - CLI 管理云服务商信息
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

// 管理 DNS API 信息
func CLIManageCloudPlatformInfo() {
	if !CheckExecRole() {
		color.Red("不允许在 follower 节点进行此操作，请前往 leader 节点执行！")
		os.Exit(0)
	}
	color.Red("更新云平台信息 —— 请选择你要进行的操作（输入对应操作的数字即可）")
	color.Green("1.查询当前云平台信息")
	color.Green("2.新增云平台信息")
	color.Green("3.修改当前云平台信息")
	color.Green("4.删除指定云平台信息")
	color.Cyan("请输入：")
	var userChoice string
	fmt.Scan(&userChoice)
	switch userChoice {
	case "1":
		color.Cyan("你选择了查询当前云平台信息")
		thisCocoa := &model.CocoaBasic{ConfigImported: true}
		CocoaDataEngine.Cols("CloudPlatformInfo").Where("configImported=true").Get(thisCocoa)

		for seq, info := range thisCocoa.CloudPlatformInfo {
			color.Green("————————————————————————————————")
			fmt.Println("顺序：", seq+1)
			fmt.Println("云平台：", info.PlatformProvider)
			fmt.Println("AccessKey：", info.AccessKey)
			fmt.Println("AccessSecret：", info.AccessSecret)
			color.Green("————————————————————————————————")
		}
		color.Green("查询成功!")

	case "2":
		color.Cyan("你选择了新增云平台信息")
		var accessKey, accessSecret, platformProvider string
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

		color.Cyan("请输入云平台的 AccessKey :")
		fmt.Scan(&accessKey)
		fmt.Println("AccessKey:", accessKey)

		color.Cyan("请输入云平台的 AccessSecret :")
		fmt.Scan(&accessSecret)
		fmt.Println("AccessSecret:", accessSecret)

		// var thisCocoa = &model.CocoaBasic{ConfigImported: true}
		// thisCocoa.CloudPlatformInfo.AccessKey = accessKey
		// thisCocoa.CloudPlatformInfo.AccessSecret = accessSecret
		// thisCocoa.CloudPlatformInfo.PlatformProvider = platformProvider

		var thisCloudPlatformInfo = &model.CloudPlatformInfo{PlatformProvider: platformProvider,
			AccessKey:    accessKey,
			AccessSecret: accessSecret}

		thisCocoa := &model.CocoaBasic{ConfigImported: true}
		thisCocoa.CloudPlatformInfo = append(thisCocoa.CloudPlatformInfo, *thisCloudPlatformInfo)

		CocoaDataEngine.Cols("CloudPlatformInfo").Update(thisCocoa)
		color.Green("新增成功!")

	case "3":
		color.Cyan("你选择了修改当前云平台信息")
		thisCocoa := &model.CocoaBasic{ConfigImported: true}
		CocoaDataEngine.Cols("CloudPlatformInfo").Where("configImported=true").Get(thisCocoa)

		for seq, info := range thisCocoa.CloudPlatformInfo {
			color.Green("————————————————————————————————")
			fmt.Println("顺序：", seq+1)
			fmt.Println("云平台：", info.PlatformProvider)
			fmt.Println("AccessKey：", info.AccessKey)
			fmt.Println("AccessSecret：", info.AccessSecret)
			color.Green("————————————————————————————————")
		}
		color.Red("请输入要修改的平台的名称:")
		var accessKey, accessSecret, platformProvider string
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

		color.Cyan("请输入云平台的 AccessKey :")
		fmt.Scan(&accessKey)
		fmt.Println("AccessKey:", accessKey)

		color.Cyan("请输入云平台的 AccessSecret :")
		fmt.Scan(&accessSecret)
		fmt.Println("AccessSecret:", accessSecret)

		var thisCloudPlatformInfo = &model.CloudPlatformInfo{PlatformProvider: platformProvider,
			AccessKey:    accessKey,
			AccessSecret: accessSecret}

		thisCocoaModified := &model.CocoaBasic{ConfigImported: true}
		thisCocoaModified.CloudPlatformInfo = append(thisCocoaModified.CloudPlatformInfo, *thisCloudPlatformInfo)

		CocoaDataEngine.Cols("CloudPlatformInfo").Update(thisCocoaModified)
		color.Green("修改成功!")

	case "4":
		color.Cyan("你选择了删除当前云平台信息")
		thisCocoa := &model.CocoaBasic{ConfigImported: true}
		CocoaDataEngine.Cols("CloudPlatformInfo").Where("configImported=true").Get(thisCocoa)

		for seq, info := range thisCocoa.CloudPlatformInfo {
			color.Green("————————————————————————————————")
			fmt.Println("顺序：", seq+1)
			fmt.Println("云平台：", info.PlatformProvider)
			fmt.Println("AccessKey：", info.AccessKey)
			fmt.Println("AccessSecret：", info.AccessSecret)
			color.Green("————————————————————————————————")
		}
		color.Red("请输入要删除的平台的名称:")
		var platformProvider string
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

		// 嵌套json需要自己处理，orm 框架不能处理
		var cloudPlatformInfo []model.CloudPlatformInfo

		for _, info := range thisCocoa.CloudPlatformInfo {
			if info.PlatformProvider != platformProvider {
				cloudPlatformInfo = append(cloudPlatformInfo, info)
			}
		}

		thisCocoaModified := &model.CocoaBasic{ConfigImported: true}
		thisCocoaModified.CloudPlatformInfo = cloudPlatformInfo

		CocoaDataEngine.Cols("CloudPlatformInfo").Update(thisCocoaModified)
		color.Green("删除成功!")

	default:
		color.Red("未输入或输入错误，你的输入是：", userChoice)
		os.Exit(1)
	}
}
