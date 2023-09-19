// CocoaSyncer - 心爱酱多节点智能解析平台 - 数据模型
//	@CreateTime		: 2023/8/30 15:10
//	@LastModified	: 2023/9/19 17:06
//	@Author			: Luckykeeper
//	@Email			: luckykeeper@luckykeeper.site
//	@Project		: CocoaSyncer

package model

import "time"

// xorm Tag 文档：https://xorm.io/zh/docs/chapter-02/4.columns/

// ________________________________________________________________
// 设置信息
type CocoaBasic struct {
	// 平台相关
	APIPort             int                   `json:"apiPort" yaml:"APIPort" xorm:"'apiPort' comment('CocoaSyncer API 服务端口')"`                 // CocoaSyncer API 服务端口
	NodeName            string                `json:"nodeName" yaml:"NodeName" xorm:"'nodeName' comment('CocoaSyncer 节点友好名称（WEB 展示及记忆用）')"`    // CocoaSyncer 节点友好名称（WEB 展示及记忆用）
	NodeCode            string                `json:"nodeCode" yaml:"NodeCode" xorm:"'nodeCode' comment('CocoaSyncer 节点代号（WEB 不展示，建议填写地区）')"`  // CocoaSyncer 节点代号（WEB 不展示，建议填写地区）
	NodeAddress         string                `json:"nodeAddress" yaml:"NodeAddress" xorm:"'nodeAddress' comment('CocoaSyncer 节点的访问地址')"`      // CocoaSyncer 节点的访问地址
	WorkMode            string                `json:"workMode" yaml:"WorkMode" xorm:"'workMode' comment('CocoaSyncer 工作模式（leader|follower）')"` // CocoaSyncer 工作模式（leader|follower）
	CocoaSecret         string                `json:"cocoaSecret" yaml:"CocoaSecret" xorm:"'cocoaSecret' comment('CocoaSyncer 通信密钥，各节点可不同')"`  // CocoaSyncer 通信密钥，各节点可不同
	ReverseProxy        bool                  `json:"reverseProxy" yaml:"ReverseProxy" xorm:"'reverseProxy' comment('CocoaSyncer 互联的其它节点')"`   // 反向代理(true=使用)
	OtherCocoaSyncer    []NodeInfo            `json:"otherCocoaSyncer" xorm:"'otherCocoaSyncer' json comment('CocoaSyncer 互联的其它节点')"`          // CocoaSyncer 互联的节点【包含本机】（不在配置文件里面配置）
	CocoaManagedService []CocoaManagedService `json:"cocoaManagedService" xorm:"'cocoaManagedService' json comment('CocoaSyncer 管理的服务')"`      // CocoaSyncer 管理的服务（全部）（不在配置文件里面配置）
	CloudPlatformInfo   []CloudPlatformInfo   `json:"cloudPlatformInfo" xorm:"'cloudPlatformInfo' json comment('CocoaSyncer 对接的云平台信息')"`       // CocoaSyncer 对接的云平台信息（不在配置文件里面配置）
	Timeout             int                   `json:"timeout" yaml:"Timeout" xorm:"'timeout' comment('超时时间')"`                                 // 超时时间

	// 数据库
	DataBaseType string `yaml:"DataBaseType" xorm:"'dataBaseType' comment('数据库类型（mysql,postgres,tidb,sqlite3,sqlite,mssql,oracle,cockroach）')"` // 数据库类型（mysql,postgres,tidb,sqlite3,sqlite,mssql,oracle,cockroach）
	Dsn          string `yaml:"DSN" xorm:"'dsn' comment('dsn ，数据库连接信息')"`                                                                       // dsn ，数据库连接信息

	// 表自身的属性
	CreatedAt time.Time `json:"createdAt" xorm:"'createdAt' created comment('数据创建时间')"`
	UpdatedAt time.Time `json:"updatedAt" xorm:"'updatedAt' updated comment('数据最后一次更新时间')"`

	ConfigImported bool `json:"configImported" xorm:"'configImported' comment('设置数据导入状态')"`
}

// 节点数据模型
type NodeInfo struct {
	NodeName            string                `json:"nodeName"`            // CocoaSyncer 节点友好名称（WEB 展示及记忆用）
	NodeCode            string                `json:"nodeCode"`            // CocoaSyncer 节点代号（WEB 不展示，建议填写地区）
	NodeAddress         string                `json:"nodeAddress"`         // CocoaSyncer 节点的访问地址
	WorkMode            string                `json:"workMode"`            // CocoaSyncer 工作模式（leader|follower）
	CocoaSecret         string                `json:"cocoaSecret"`         // CocoaSyncer 通信密钥，各节点可不同
	ISPs                []string              `json:"ISPs"`                // 域名解析线路 (default|CT|CU|CMCC|CNEdu|OverSea|None)[默认|中国电信|中国联通|中国移动|中国教育网|境外|不解析]
	NetworkType         string                `json:"networkType"`         // 节点网络类型（IPV4|IPV6|Both）[仅使用IPV4|仅使用IPV6|IPV4&IPV6双栈]
	StatusCode          int                   `json:"statusCode"`          // 节点状态码，200 正常，504 失联
	WatchDog            time.Time             `json:"watchDog"`            // 看门狗，判断master是否存活
	CocoaManagedService []CocoaManagedService `json:"cocoaManagedService"` // 节点管理的服务
}

// 服务模型
type CocoaManagedService struct {
	ServiceName      string `json:"serviceName"`      // 服务友好名称
	Domain           string `json:"domain"`           // 服务域名名称(e.g. luckykeeper.site)
	Record           string `json:"record"`           // 主机名称(e.g. www)
	PlatformProvider string `json:"platformProvider"` // 云平台服务商
}

// DNS 平台信息
// 目前支持：阿里云
type CloudPlatformInfo struct {
	PlatformProvider string `json:"platformProvider"` // 云平台服务商
	AccessKey        string `json:"accessKey"`        // DNS API账号的 AccessKey
	AccessSecret     string `json:"accessSecret"`     // DNS API账号的 AccessSecret
}

// ________________________________________________________________
// 返回数据
type ServerReturn struct {
	StatusCode   int        `json:"statusCode"`   // 业务码 （操作成功200，Token错误401）
	StatusString string     `json:"statusString"` // 业务码的文字说明
	CocoaBasic   CocoaBasic `json:"cocoaBasic"`   // 业务数据
}

// ________________________________________________________________
// 监控接口 - 请求数据
type APIRequest struct {
	CocoaSecret string     `json:"cocoaSecret"` // 节点的请求密钥
	CocoaBasic  CocoaBasic `json:"cocoaBasic"`  // 更新节点信息
}

// ________________________________________________________________
// 业务接口
