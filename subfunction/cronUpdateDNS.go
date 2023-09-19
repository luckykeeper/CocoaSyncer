// CocoaSyncer - 心爱酱多节点智能解析平台 - 定时更新DNS解析信息
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

	cocoaAliDNS "github.com/denverdino/aliyungo/dns"

	"github.com/fatih/color"
	"github.com/honwen/golibs/cip"
)

// 根据线路设置，定时更新 DNS 解析
func CronUpdateDNS() {
	color.Cyan("___________________________")
	log.Println("CronTaskCenter - CronUpdateDNS - Start")
	log.Println("定时任务中心 - 更新 DNS 解析记录 - 开始")

	// 拿本地数据
	thisCocoa := model.CocoaBasic{ConfigImported: true}
	CocoaDataEngine.Get(&thisCocoa)

	var (
		currentServiceInfo []model.CocoaManagedService
		currentNeworkType  string
		currentISPs        []string
	)

	// 获取线路，网络模式，云平台账号（thisCocoa）和服务（仅本端）
	for _, info := range thisCocoa.OtherCocoaSyncer {
		if info.NodeCode == thisCocoa.NodeCode {
			currentServiceInfo = info.CocoaManagedService
			currentNeworkType = info.NetworkType
			currentISPs = info.ISPs
		}
	}
	// log.Println(currentServiceInfo)
	// log.Println(currentNeworkType)
	// log.Println(currentISPs)

	// failSafe 机制，对端节点故障时临时将对端节点的线路解析到本端
	for _, nodeInfo := range thisCocoa.OtherCocoaSyncer {
		// 对端节点
		if nodeInfo.NodeCode != thisCocoa.NodeCode {
			if nodeInfo.StatusCode == 504 {
				currentISPs = append(currentISPs, nodeInfo.ISPs...)
				color.Yellow("failSafe!!!")
				color.Yellow("特别注意：检测到对端节点故障，将对端节点的线路转移到本端解析!")
			}
		}
	}

	// 获取公网 IP 地址，使用 cip 库
	publicIPV4 := cip.MyIPv4()
	publicIPV6 := cip.MyIPv6()

	var clientHaveIPV4 bool
	if len(publicIPV4) == 0 {
		color.Red("没有获取到公网 IPV4 地址，网络连接可能异常!")
		clientHaveIPV4 = false
	} else {
		color.Cyan("当前公网 IPV4:", publicIPV4)
		clientHaveIPV4 = true
	}

	var clientHaveIPV6 bool
	if len(publicIPV6) == 0 {
		color.Yellow("没有获取到公网 IPV6 地址，你可能没有 IPV6 连接!")
		clientHaveIPV6 = false
	} else {
		color.Cyan("当前公网 IPV6:", publicIPV6)
		clientHaveIPV6 = true
	}

	var cocoaAliClient *cocoaAliDNS.Client
	// 当前仅支持阿里云
	for _, provider := range thisCocoa.CloudPlatformInfo {
		if provider.PlatformProvider == "Ali" {
			cocoaAliClient = cocoaAliDNS.NewClient(provider.AccessKey, provider.AccessSecret)
		}
	}

	// 执行 IPV4 DDNS
	if clientHaveIPV4 {
		// 检查节点是否配置了该解析线路
		if currentNeworkType == "IPV4" || currentNeworkType == "Both" {
			color.Green("开始节点 IPV4 解析")
			for _, service := range currentServiceInfo {
				color.Green("开始节点 IPV4 解析 - 阿里云")
				// 当前仅支持阿里云
				if service.PlatformProvider == "Ali" {
					domainInfoToBeDesribe := cocoaAliDNS.DescribeSubDomainRecordsArgs{SubDomain: service.Record + "." + service.Domain}
					domainInfo, err := cocoaAliClient.DescribeSubDomainRecords(&domainInfoToBeDesribe)
					if err != nil {
						color.Red("获取域名解析信息失败，详细错误信息如下：")
						log.Println(err)
						continue
					} else {
						// 解析线路与线上信息一致，更新解析信息
						for _, isp := range currentISPs {
							var aliLine string
							// 运营商与阿里对应转换
							switch isp {
							case "default":
								// 默认线路
								aliLine = "default"
							case "CT":
								// 中国电信
								aliLine = "telecom"

							case "CU":
								// 中国联通
								aliLine = "unicom"

							case "CMCC":
								// 中国移动
								aliLine = "mobile"

							case "CNEdu":
								// 中国教育网
								aliLine = "edu"

							case "OverSea":
								aliLine = "oversea"
								// 海外线路

							case "None":
								// 不解析
								continue

							default:
								// 不解析
								continue
							}

							// 根据解析线路更新解析信息
							for _, domainRecordSpecInfo := range domainInfo.DomainRecords.Record {
								// 要解析的域名必须启用
								if domainRecordSpecInfo.Status == "ENABLE" {
									// 解析信息是 A 类型
									if domainRecordSpecInfo.Type == "A" {
										if aliLine == domainRecordSpecInfo.Line {
											// dev: 匹配到的详细信息
											// log.Println("dev!!!")
											// log.Println("domainRecordSpecInfo:", domainRecordSpecInfo)

											// 做个锁，不能同时执行
											var usingTag bool
											usingTag = false
											for usingTag {
												time.Sleep(1 * time.Second)
												log.Print("休眠一秒等待上个域名解析任务完成", service.ServiceName, isp)
											}
											usingTag = true
											if domainRecordSpecInfo.Value != publicIPV4 {

												// 更新解析信息
												var domainRecordUpdateInfo cocoaAliDNS.UpdateDomainRecordArgs
												domainRecordUpdateInfo = cocoaAliDNS.UpdateDomainRecordArgs{RecordId: domainRecordSpecInfo.RecordId, RR: domainRecordSpecInfo.RR,
													Type: domainRecordSpecInfo.Type, Value: publicIPV4, Line: domainRecordSpecInfo.Line}
												_, err := cocoaAliClient.UpdateDomainRecord(&domainRecordUpdateInfo)

												if err != nil {
													color.Red("更新解析记录失败，详细信息如下：")
													log.Println("domainRecordUpdateInfo:", domainRecordUpdateInfo.Line, domainRecordUpdateInfo.RR,
														domainRecordUpdateInfo.Type, domainRecordUpdateInfo.Value)
													log.Println("err:", err)
												} else {
													color.Green("更新：", domainRecordUpdateInfo)
												}
												time.Sleep(3 * time.Second)
												usingTag = false
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	// // 执行 IPV6 DDNS
	if clientHaveIPV6 {
		// 检查节点是否配置了该解析线路
		if currentNeworkType == "IPV6" || currentNeworkType == "Both" {
			color.Green("开始节点 IPV6 解析")
			for _, service := range currentServiceInfo {
				color.Green("开始节点 IPV6 解析 - 阿里云")
				// 当前仅支持阿里云
				if service.PlatformProvider == "Ali" {
					domainInfoToBeDesribe := cocoaAliDNS.DescribeSubDomainRecordsArgs{SubDomain: service.Record + "." + service.Domain}
					domainInfo, err := cocoaAliClient.DescribeSubDomainRecords(&domainInfoToBeDesribe)
					if err != nil {
						color.Red("获取域名解析信息失败，详细错误信息如下：")
						log.Println(err)
						continue
					} else {
						// 解析线路与线上信息一致，更新解析信息
						for _, isp := range currentISPs {
							var aliLine string
							// 运营商与阿里对应转换
							switch isp {
							case "default":
								// 默认线路
								aliLine = "default"
							case "CT":
								// 中国电信
								aliLine = "telecom"

							case "CU":
								// 中国联通
								aliLine = "unicom"

							case "CMCC":
								// 中国移动
								aliLine = "mobile"

							case "CNEdu":
								// 中国教育网
								aliLine = "edu"

							case "OverSea":
								aliLine = "oversea"
								// 海外线路

							case "None":
								// 不解析
								continue

							default:
								// 不解析
								continue
							}

							// 根据解析线路更新解析信息
							for _, domainRecordSpecInfo := range domainInfo.DomainRecords.Record {
								// 要解析的域名必须启用
								if domainRecordSpecInfo.Status == "ENABLE" {
									// 解析信息是 A 类型
									if domainRecordSpecInfo.Type == "AAAA" {
										if aliLine == domainRecordSpecInfo.Line {
											// dev: 匹配到的详细信息
											// log.Println("dev!!!")
											// log.Println("domainRecordSpecInfo:", domainRecordSpecInfo)

											// 做个锁，不能同时执行
											var usingTag bool
											usingTag = false
											for usingTag {
												time.Sleep(1 * time.Second)
												log.Print("休眠一秒等待上个域名解析任务完成", service.ServiceName, isp)
											}
											usingTag = true
											if domainRecordSpecInfo.Value != publicIPV6 {

												// 更新解析信息
												var domainRecordUpdateInfo cocoaAliDNS.UpdateDomainRecordArgs
												domainRecordUpdateInfo = cocoaAliDNS.UpdateDomainRecordArgs{RecordId: domainRecordSpecInfo.RecordId, RR: domainRecordSpecInfo.RR,
													Type: domainRecordSpecInfo.Type, Value: publicIPV6, Line: domainRecordSpecInfo.Line}
												_, err := cocoaAliClient.UpdateDomainRecord(&domainRecordUpdateInfo)

												if err != nil {
													color.Red("更新解析记录失败，详细信息如下：")
													log.Println("domainRecordUpdateInfo:", domainRecordUpdateInfo.Line, domainRecordUpdateInfo.RR,
														domainRecordUpdateInfo.Type, domainRecordUpdateInfo.Value)
													log.Println("err:", err)
												} else {
													color.Green("更新：", domainRecordUpdateInfo)
												}
												time.Sleep(3 * time.Second)
												usingTag = false
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	log.Println("CronTaskCenter - CronUpdateDNS - Finished")
	log.Println("定时任务中心 - 更新 DNS 解析记录 - 完成")
	color.Cyan("___________________________")
}
