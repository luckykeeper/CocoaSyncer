// CocoaSyncer - 心爱酱多节点智能解析平台 - 拉对端数据
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
	"io"
	"log"
	"net/http"

	"github.com/fatih/color"
)

// 拉对端狗数据
func pullClusterInfo(accessInfo model.NodeInfo) (nodeInfo model.NodeInfo, operationStatus bool) {
	url := accessInfo.NodeAddress + "/v1/getClusterInfo"
	method := "POST"

	var apiRequest model.APIRequest
	apiRequest.CocoaSecret = accessInfo.CocoaSecret
	apiRequestData, _ := json.Marshal(apiRequest)
	payload := bytes.NewBuffer(apiRequestData)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		color.Red(err.Error())
		return nodeInfo, false
	}
	req.Header.Add("User-Agent", "CocoaSyncer By Luckykeeper")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		color.Red(err.Error())
		return nodeInfo, false
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		color.Red(err.Error())
		return nodeInfo, false
	}

	if res.StatusCode != 200 {
		color.Red("HTTP 请求失败!")
		log.Println("详细信息:", body)
		return nodeInfo, false
	}

	var serverReturn model.ServerReturn
	json.Unmarshal(body, &serverReturn)

	if serverReturn.StatusCode != 200 {
		color.Red("与节点 " + accessInfo.NodeCode + " 通信失败，原因为：")
		color.Red(serverReturn.StatusString)
		return nodeInfo, false
	}

	// 仅返回对端节点的数据
	for _, nodeInfo := range serverReturn.CocoaBasic.OtherCocoaSyncer {
		if nodeInfo.NodeCode == accessInfo.NodeCode {
			return nodeInfo, true
		}
	}

	return nodeInfo, false
}
