// CocoaSyncer - 心爱酱多节点智能解析平台 - 推对端数据
// @CreateTime : 2023/8/30 15:10
// @LastModified : 2023/8/30 15:10
// @Author : Luckykeeper
// @Email : luckykeeper@luckykeeper.site
// @Project : CocoaSyncer

package subFunction

import (
	"bytes"
	"cocoaSyncer/model"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/fatih/color"
)

// 推送数据到对端
func pushClusterInfo(accessInfo model.NodeInfo, clusterInfo model.CocoaBasic) bool {
	url := accessInfo.NodeAddress + "/v1/updateClusterInfo"
	method := "POST"

	var apiRequest model.APIRequest
	apiRequest.CocoaSecret = accessInfo.CocoaSecret

	// 节点信息
	apiRequest.CocoaBasic.OtherCocoaSyncer = clusterInfo.OtherCocoaSyncer

	// DNS 账号信息
	apiRequest.CocoaBasic.CloudPlatformInfo = clusterInfo.CloudPlatformInfo

	// 服务信息
	apiRequest.CocoaBasic.CocoaManagedService = clusterInfo.CocoaManagedService

	apiRequestData, _ := json.Marshal(apiRequest)
	payload := bytes.NewBuffer(apiRequestData)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		color.Red(err.Error())
		return false
	}
	req.Header.Add("User-Agent", "CocoaSyncer By Luckykeeper")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		color.Red(err.Error())
		return false
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		color.Red(err.Error())
		return false
	}

	if res.StatusCode != 200 {
		color.Red("HTTP 请求失败!")
		log.Println("详细信息:", body)
		return false
	}

	var serverReturn model.ServerReturn
	json.Unmarshal(body, &serverReturn)

	if serverReturn.StatusCode != 200 {
		color.Red("与节点 " + accessInfo.NodeCode + " 通信失败，原因为：")
		color.Red(serverReturn.StatusString)
		return false
	}

	return true
}
