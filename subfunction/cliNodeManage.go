// CocoaSyncer - 心爱酱多节点智能解析平台 - CLI 节点管理
//	@CreateTime		: 2023/8/30 15:10
//	@LastModified	: 2023/8/30 15:10
//	@Author			: Luckykeeper
//	@Email			: luckykeeper@luckykeeper.site
//	@Project		: CocoaSyncer

package subFunction

import (
	"bytes"
	"cocoaSyncer/model"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/fatih/color"
)

// CLI 节点管理
func CLINodeManage() {
	if !CheckExecRole() {
		color.Red("不允许在 follower 节点进行此操作，请前往 leader 节点执行！")
		os.Exit(0)
	}
	color.Red("节点管理 —— 请选择你要进行的操作（输入对应操作的数字即可）")
	color.Green("1.查询当前节点信息")
	color.Green("2.新增节点")
	color.Green("3.修改节点信息")
	color.Green("4.删除节点")
	color.Cyan("请输入：")
	var userChoice string
	fmt.Scan(&userChoice)
	switch userChoice {
	case "1":
		color.Cyan("你选择了查询当前节点信息")
		thisCocoa := &model.CocoaBasic{ConfigImported: true}
		CocoaDataEngine.Cols("otherCocoaSyncer").Where("configImported=true").Get(thisCocoa)
		color.Cyan("___________________________")
		color.Cyan("集群运行状态如下：")
		for seq, info := range thisCocoa.OtherCocoaSyncer {
			color.Green("————————————————————————————————")
			fmt.Println("顺序：", seq+1)
			if info.NodeCode == thisCocoa.NodeCode {
				color.Cyan("本端")
			} else {
				color.Yellow("对端")
			}
			fmt.Println("节点名称：", info.NodeName)
			fmt.Println("节点代号：", info.NodeCode)
			fmt.Println("节点对外地址：", info.NodeAddress)
			fmt.Print("节点工作模式：")
			if info.WorkMode == "leader" {
				color.Cyan(info.WorkMode)
			} else {
				color.Yellow(info.WorkMode)
			}
			fmt.Println("节点通信密钥：", info.CocoaSecret)
			fmt.Println("节点线路：", info.ISPs)
			fmt.Println("节点网络模式：", info.NetworkType)
			fmt.Println("节点状态码（200正常，504失联）：", info.StatusCode)
			fmt.Println("节点上次喂狗时间：", info.WatchDog)
			fmt.Println("节点管理服务：", info.CocoaManagedService)
			color.Green("————————————————————————————————")
		}
		color.Cyan("___________________________")

		color.Green("查询成功!")

	case "2":
		color.Cyan("你选择了新增受管节点")
		var newNode model.NodeInfo

		color.Cyan("请输入节点访问地址，需要带完整协议，非标端口需要加端口号，如“https://www.example.com:8088”")
		color.Cyan("请输入：")
		fmt.Scan(&newNode.NodeAddress)
		fmt.Println("节点访问地址:", newNode.NodeAddress)

		color.Cyan("请输入节点通信密钥")
		color.Cyan("请输入：")
		fmt.Scan(&newNode.CocoaSecret)
		fmt.Println("节点通信密钥:", newNode.CocoaSecret)

		// 请求该节点的 /v1/getClusterInfo , 获取其它数据
		url := newNode.NodeAddress + "/v1/getClusterInfo"
		method := "POST"

		// var apiRequest model.APIRequest
		// apiRequest.CocoaSecret = newNode.CocoaSecret
		apiRequestData, _ := json.Marshal(newNode)
		payload := bytes.NewBuffer(apiRequestData)

		client := &http.Client{}
		req, err := http.NewRequest(method, url, payload)

		if err != nil {
			color.Red("从对端节点请求数据失败，错误原因:", err)
			os.Exit(1)
		}
		req.Header.Add("User-Agent", "CocoaSyncer By Luckykeeper")
		req.Header.Add("Content-Type", "application/json")

		res, err := client.Do(req)
		if err != nil {
			color.Red("从对端节点请求数据失败，错误原因:", err)
			os.Exit(1)
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			color.Red("从对端节点请求数据失败，错误原因:", err)
			os.Exit(1)
		}
		if res.StatusCode == 200 {
			color.Green("网络请求成功")
		} else {
			color.Red("网络请求失败！检查对端地址和密码是否输入正确，端口是否开放，输入的对端地址结尾不能带“/”")
			os.Exit(1)
		}
		var nodeReturn model.ServerReturn
		json.Unmarshal(body, &nodeReturn)
		var newNodeInfo model.NodeInfo
		for _, nodeInfo := range nodeReturn.CocoaBasic.OtherCocoaSyncer {
			newNodeInfo = nodeInfo
		}
		if nodeReturn.StatusCode == 200 {
			color.Red("对端节点已经加入集群，无法重复添加！对端返回:", string(body))
			os.Exit(2)
		} else if nodeReturn.StatusCode == 401 {
			color.Red("节点通信密钥错误！对端返回:", string(body))
			os.Exit(3)
		} else if nodeReturn.StatusCode == 201 {
			json.Unmarshal(body, &newNode)
			color.Green("获取对端节点数据成功！对端数据如下:")
			color.Green("————————————————————————————————")
			fmt.Println("节点名称：", newNodeInfo.NodeName)
			fmt.Println("节点代号：", newNodeInfo.NodeCode)
			fmt.Println("节点对外地址：", newNodeInfo.NodeAddress)
			fmt.Print("节点工作模式：")
			if newNodeInfo.WorkMode == "leader" {
				color.Cyan(newNodeInfo.WorkMode)
			} else {
				color.Yellow(newNodeInfo.WorkMode)
			}
			fmt.Println("节点通信密钥：", newNodeInfo.CocoaSecret)
			fmt.Println("节点线路：", newNodeInfo.ISPs)
			fmt.Println("节点网络模式：", newNodeInfo.NetworkType)
			fmt.Println("节点状态码（200正常，504失联）：", newNodeInfo.StatusCode)
			fmt.Println("节点上次喂狗时间：", newNodeInfo.WatchDog)
			fmt.Println("节点管理服务：", newNodeInfo.CocoaManagedService)
			color.Green("————————————————————————————————")

			color.Yellow("请检查以上信息，是否确认添加(y/n):")
			var userChoice string
			fmt.Scan(&userChoice)
			if userChoice == "y" {
				color.Green("确认添加节点:", newNodeInfo.NodeName, "-", newNodeInfo.NodeCode)

				// 本端数据库操作
				thisCocoa := &model.CocoaBasic{ConfigImported: true}
				CocoaDataEngine.Cols("otherCocoaSyncer").Where("configImported=true").Get(thisCocoa)

				thisCocoa.OtherCocoaSyncer = append(thisCocoa.OtherCocoaSyncer, newNodeInfo)
				CocoaDataEngine.Cols("otherCocoaSyncer").Where("configImported=true").Update(thisCocoa)

				fmt.Println("添加节点成功:", newNodeInfo.NodeName, "-", newNodeInfo.NodeCode)

			} else {
				fmt.Println("取消添加节点:", newNodeInfo.NodeName, "-", newNodeInfo.NodeCode)
				color.Yellow("程序退出")
				os.Exit(0)
			}

		} else if nodeReturn.StatusCode == 0 {
			color.Red("网络请求失败！检查对端地址和密码是否输入正确，端口是否开放，输入的对端地址结尾不能带“/”")
			log.Println(nodeReturn)
			os.Exit(1)
		}

		color.Green("新增成功!")

	case "3":
		color.Cyan("你选择了修改受管节点信息")
		thisCocoa := &model.CocoaBasic{ConfigImported: true}
		CocoaDataEngine.Cols("otherCocoaSyncer").Where("configImported=true").Get(thisCocoa)
		color.Cyan("___________________________")
		color.Cyan("集群运行状态如下：")
		for seq, info := range thisCocoa.OtherCocoaSyncer {
			color.Green("————————————————————————————————")
			fmt.Println("顺序：", seq+1)
			if info.NodeCode == thisCocoa.NodeCode {
				color.Cyan("本端")
			} else {
				color.Yellow("对端")
			}
			fmt.Println("节点名称：", info.NodeName)
			fmt.Println("节点代号：", info.NodeCode)
			fmt.Println("节点对外地址：", info.NodeAddress)
			fmt.Print("节点工作模式：")
			if info.WorkMode == "leader" {
				color.Cyan(info.WorkMode)
			} else {
				color.Yellow(info.WorkMode)
			}
			fmt.Println("节点通信密钥：", info.CocoaSecret)
			fmt.Println("节点线路：", info.ISPs)
			fmt.Println("节点网络模式：", info.NetworkType)
			fmt.Println("节点状态码（200正常，504失联）：", info.StatusCode)
			fmt.Println("节点上次喂狗时间：", info.WatchDog)
			fmt.Println("节点管理服务：", info.CocoaManagedService)
			color.Green("————————————————————————————————")
		}
		color.Cyan("___________________________")
		var (
			nodeChoice                   string
			nodeCurrentInfo, nodeNewInfo model.NodeInfo
			// nodeCurrentInfo model.NodeInfo
		)
		color.Red("请输入需要修改的节点代号:")
		fmt.Scan(&nodeChoice)
		fmt.Println("将要准备修改:", nodeChoice)

		// 拿该节点信息
		// thisCocoa := &model.CocoaBasic{ConfigImported: true}
		// CocoaDataEngine.Cols("otherCocoaSyncer").Where("configImported=true").Get(thisCocoa)
		for _, info := range thisCocoa.OtherCocoaSyncer {
			if info.NodeCode == nodeChoice {
				nodeCurrentInfo = info
			}
		}
		if len(nodeCurrentInfo.NodeCode) == 0 {
			color.Red("节点代号输入有误")
			os.Exit(4)
		}
		log.Println(nodeCurrentInfo)

		color.Green("输入想要修改的内容对应的数字:")
		color.Cyan("___________________________")
		color.Green("1.修改节点对外地址")
		color.Green("2.修改节点线路")
		color.Green("3.修改节点网络模式")
		color.Cyan("___________________________")
		color.Red("请输入:")
		var operationChoice string
		fmt.Scan(&operationChoice)
		switch operationChoice {
		case "1":
			color.Red("你选择了修改节点对外地址")
			fmt.Println("当前节点对外地址:", nodeCurrentInfo.NodeAddress)
			color.Cyan("请输入节点访问地址，需要带完整协议，非标端口需要加端口号，如“https://www.example.com:8088”")
			fmt.Scan(&nodeNewInfo.NodeAddress)
			fmt.Println("新的对外地址:", nodeNewInfo.NodeAddress)

			thisCocoaModified := &model.CocoaBasic{ConfigImported: true}
			// 节点信息，通过遍历循环修改
			for _, info := range thisCocoa.OtherCocoaSyncer {
				if info.NodeCode == nodeChoice {
					info.NodeAddress = nodeNewInfo.NodeAddress
				}
				thisCocoaModified.OtherCocoaSyncer = append(thisCocoaModified.OtherCocoaSyncer, info)
			}

			CocoaDataEngine.Cols("otherCocoaSyncer").Update(thisCocoaModified)

		case "2":
			color.Red("你选择了修改节点线路")
			fmt.Println("当前节点线路:", nodeCurrentInfo.ISPs)
			fmt.Println("请输入新的节点线路代号，节点线路列表见下方，多个线路代号用英文逗号“,”隔开")
			fmt.Println("线路代号如下：")
			color.Cyan("___________________________")
			color.Yellow("default:默认解析线路(必须有一个节点有它)")
			color.Green("CT:中国电信")
			color.Green("CU:中国联通")
			color.Green("CMCC:中国移动")
			color.Green("CNEdu:中国教育网")
			color.Green("OverSea:海外")
			color.Green("None:本节点不进行解析")
			color.Cyan("___________________________")
			var userInputISPs string
			fmt.Scan(&userInputISPs)
			nodeNewInfo.ISPs = strings.Split(userInputISPs, ",")
			fmt.Println("新的节点线路:", nodeNewInfo.ISPs)

			thisCocoaModified := &model.CocoaBasic{ConfigImported: true}
			// 节点信息，通过遍历循环修改
			for _, info := range thisCocoa.OtherCocoaSyncer {
				if nodeChoice == info.NodeCode {
					info.ISPs = nodeNewInfo.ISPs
				}
				thisCocoaModified.OtherCocoaSyncer = append(thisCocoaModified.OtherCocoaSyncer, info)
			}

			CocoaDataEngine.Cols("otherCocoaSyncer").Update(thisCocoaModified)

		case "3":
			color.Red("你选择了修改节点网络模式")
			fmt.Println("当前节点网络模式:", nodeCurrentInfo.NetworkType)
			fmt.Println("请输入新的节点网络模式代号:")
			fmt.Println("网络模式代号如下：")
			color.Cyan("___________________________")
			color.Green("IPV4 :纯 IPV4 线路")
			color.Green("IPV6 :纯 IPV6 线路")
			color.Green("Both :IPv4&IPv6 双栈")
			color.Cyan("___________________________")
			var userInputNetworkType string
			fmt.Scan(&userInputNetworkType)
			fmt.Println("新的网络模式:", userInputNetworkType)

			thisCocoaModified := &model.CocoaBasic{ConfigImported: true}
			// 节点信息，通过遍历循环修改
			for _, info := range thisCocoa.OtherCocoaSyncer {
				if nodeChoice == info.NodeCode {
					info.NetworkType = userInputNetworkType
				}
				thisCocoaModified.OtherCocoaSyncer = append(thisCocoaModified.OtherCocoaSyncer, info)
			}

			CocoaDataEngine.Cols("otherCocoaSyncer").Update(thisCocoaModified)

		default:
			color.Red("未输入或输入错误，你的输入是：", operationChoice)
			os.Exit(1)
		}

		color.Green("修改成功!")

	case "4":
		color.Red("你选择了删除节点信息")
		color.Red("先关闭对应节点的服务，再进行该操作，请务必不要删除 leader 节点（本端）")
		color.Red("再次提醒，危险操作，请特别注意")
		var (
			nodeChoice      string
			nodeCurrentInfo model.NodeInfo
		)
		// 拿该节点信息
		thisCocoa := &model.CocoaBasic{ConfigImported: true}
		CocoaDataEngine.Where("configImported=true").Get(thisCocoa)
		for seq, info := range thisCocoa.OtherCocoaSyncer {
			color.Green("————————————————————————————————")
			fmt.Println("顺序：", seq+1)
			if info.NodeCode == thisCocoa.NodeCode {
				color.Cyan("本端")
			} else {
				color.Yellow("对端")
			}
			fmt.Println("节点名称：", info.NodeName)
			fmt.Println("节点代号：", info.NodeCode)
			fmt.Println("节点对外地址：", info.NodeAddress)
			fmt.Print("节点工作模式：")
			if info.WorkMode == "leader" {
				color.Cyan(info.WorkMode)
			} else {
				color.Yellow(info.WorkMode)
			}
			fmt.Println("节点通信密钥：", info.CocoaSecret)
			fmt.Println("节点线路：", info.ISPs)
			fmt.Println("节点网络模式：", info.NetworkType)
			fmt.Println("节点状态码（200正常，504失联）：", info.StatusCode)
			fmt.Println("节点上次喂狗时间：", info.WatchDog)
			fmt.Println("节点管理服务：", info.CocoaManagedService)
			color.Green("————————————————————————————————")
		}
		color.Red("请输入需要修改的节点代号:")
		fmt.Scan(&nodeChoice)
		fmt.Println("删除:", nodeChoice)

		for _, info := range thisCocoa.OtherCocoaSyncer {
			if info.NodeCode == nodeChoice {
				nodeCurrentInfo = info
			}
		}
		if len(nodeCurrentInfo.NodeCode) == 0 {
			color.Red("节点代号输入有误")
			os.Exit(4)
		}

		thisCocoaModified := &model.CocoaBasic{ConfigImported: true}
		// 节点信息，通过遍历循环删除
		for _, info := range thisCocoa.OtherCocoaSyncer {
			if !(nodeChoice == info.NodeCode) {
				thisCocoaModified.OtherCocoaSyncer = append(thisCocoaModified.OtherCocoaSyncer, info)
			}
		}

		CocoaDataEngine.Cols("otherCocoaSyncer").Update(thisCocoaModified)

		color.Green("删除成功!")

	default:
		color.Red("未输入或输入错误，你的输入是：", userChoice)
		os.Exit(1)
	}
}
