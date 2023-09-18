// CocoaSyncer - 心爱酱多节点智能解析平台 - CLI 展示运行状态
// @CreateTime : 2023/8/30 15:10
// @LastModified : 2023/8/30 15:10
// @Author : Luckykeeper
// @Email : luckykeeper@luckykeeper.site
// @Project : CocoaSyncer

package subFunction

import (
	"cocoaSyncer/model"
	"fmt"
	"net/http"
	"strconv"

	"github.com/fatih/color"
)

// CLI 查看节点状态
func CLIShowStatus() {
	color.Cyan("以下是 CocoaSyncer 当前的运行状态：")
	lineStatusOk := color.New(color.FgGreen)
	lineStatusOops := color.New(color.FgRed)
	lineNormal := color.New(color.FgBlue)
	line := color.New(color.FgBlue)
	line = line.Add(color.Bold)
	line = line.Add(color.BgCyan)

	// 彩色CLI 输出示例
	// line.Println("___________________________")
	// lineNormal.Print("节点启动状态：")
	// lineStatusOk.Print("Ok | ")
	// lineStatusOops.Println("寄!")
	// line.Println("___________________________")

	// 获取配置信息
	thisCocoa := &model.CocoaBasic{ConfigImported: true}
	CocoaDataEngine.Where("configImported=true").Get(thisCocoa)
	// log.Println(len(thisCocoa.OtherCocoaSyncer))

	line.Println("___________________________")

	// 检测 CocoaSyncer 是否上线
	if checkCocoaSyncer(*thisCocoa) {
		lineNormal.Print("节点启动状态：")
		lineStatusOk.Println("启动且正常运行")
	} else {
		lineNormal.Print("节点启动状态：")
		lineStatusOops.Println("寄！")
	}
	fmt.Println("本节点其它信息：")
	fmt.Println("API端口：", thisCocoa.APIPort)
	fmt.Println("使用反代：", thisCocoa.ReverseProxy)
	fmt.Println("超时时间(min)：", thisCocoa.Timeout)

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
	color.Cyan("云平台信息如下：")
	for seq, info := range thisCocoa.CloudPlatformInfo {
		color.Green("————————————————————————————————")
		fmt.Println("顺序：", seq+1)
		fmt.Println("云平台：", info.PlatformProvider)
		fmt.Println("AccessKey：", info.AccessKey)
		fmt.Println("AccessSecret：", info.AccessSecret)
		color.Green("————————————————————————————————")
	}
	line.Println("___________________________")

}

// 检测 CocoaSyncer 是否上线
func checkCocoaSyncer(CocoaBasic model.CocoaBasic) bool {
	url := "http://127.0.0.1:" + strconv.Itoa(CocoaBasic.APIPort) + "/v1/status"
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return false
	}
	req.Header.Add("User-Agent", "CocoaSyncer By Luckykeeper")

	res, err := client.Do(req)
	if err != nil {
		return false
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return false
	}
	return true
}
