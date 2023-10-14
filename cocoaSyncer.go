// CocoaSyncer - 心爱酱多节点智能解析平台 - 程序入点
//	@CreateTime		: 2023/8/30 15:10
//	@LastModified	: 2023/9/19 17:06
//	@Author			: Luckykeeper
//	@Email			: luckykeeper@luckykeeper.site
//	@Project		: CocoaSyncer
//	@contact.name	Luckykeeper
//	@contact.url	https://github.com/luckykeeper
//	@contact.email	luckykeeper@luckykeeper.site
//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// 相关参考文档：
// 阿里云 DNS API：
// https://help.aliyun.com/document_detail/2355661.html?spm=a2c4g.39863.0.0.79c31355NpSART
// https://api.aliyun.com/api-tools/demo/Alidns/7a52d9f9-59b0-43bd-96fc-8d475220f56f
// https://api.aliyun.com/api/Alidns/2015-01-09/UpdateDomainRecord?spm=api-workbench.CodeSample%20Detail%20Page.0.0.318359978Y2vdu
// https://help.aliyun.com/document_detail/29807.html?spm=api-workbench.api_explorer.0.0.46b114b5xHKXr3
// Viper: https://github.com/spf13/viper
// gin-metrics: https://github.com/penglongli/gin-metrics

// 运作方式：首先应搭建 leader 节点，follower 节点仅配置服务端口，友好名称，代号，访问地址，工作模式，密钥，数据库，然后在 leader 上添加节点，同步数据
// leader 功能：同步数据、更改数据（服务、线路、节点名称）、切换角色
// follower 功能：接收数据、切换角色
// 数据推送：DDNS 是定时的，为了防止同时在一条线路发起多次解析，应当约定一个时间，各节点统一在该时间执行时间变更，在此之前检测节点状态
// 服务端和客户端共用一套代码

package main

import (
	"cocoaSyncer/model"
	cocoaSyncerRouter "cocoaSyncer/router"
	subFunction "cocoaSyncer/subfunction"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

// 基础变量设置
var (
	CocoaBasic model.CocoaBasic
)

// CLI
func CocoaSyncerCLI() {
	CocoaSyncer := &cli.App{
		Name: "CocoaSyncer",
		Usage: "CocoaSyncer - 心爱酱多节点智能解析平台" +
			"\nPowered By Luckykeeper <luckykeeper@luckykeeper.site | https://luckykeeper.site>" +
			"\n————————————————————————————————————————" +
			"\n注意：使用前需要先填写同目录下 config.yaml !",
		Version: "1.0.1_build20230919",
		Commands: []*cli.Command{
			{
				Name:    "runProd",
				Aliases: []string{"r"},
				Usage:   "启动 CocoaSyncer (生产环境)",
				Action: func(cCtx *cli.Context) error {
					CocoaSyncer(false)
					return nil
				},
			},
			{
				Name:    "runDebug",
				Aliases: []string{"rd"},
				Usage:   "启动 CocoaSyncer (调试)",
				Action: func(cCtx *cli.Context) error {
					CocoaSyncer(true)
					return nil
				},
			},
			{
				Name:    "showStatus",
				Aliases: []string{"status"},
				Usage:   "输出当前运行状态",
				Action: func(cCtx *cli.Context) error {
					color.Cyan("Welcome To Use CocoaSyncer!")
					color.Cyan("Reading Configuration...")
					readConfig(true)
					subFunction.InitializeDatabase(CocoaBasic, false)
					subFunction.CLIShowStatus()
					return nil
				},
			},
			{
				Name:    "updateCloudPlatformInfo",
				Aliases: []string{"dns"},
				Usage:   "更新云平台信息（更新 DNS API 账户）",
				Action: func(cCtx *cli.Context) error {
					color.Cyan("Welcome To Use CocoaSyncer!")
					color.Cyan("Reading Configuration...")
					readConfig(true)
					subFunction.InitializeDatabase(CocoaBasic, false)
					subFunction.CLIManageCloudPlatformInfo()
					return nil
				},
			},
			{
				Name:    "updateManagedServices",
				Aliases: []string{"service"},
				Usage:   "更新服务信息（要 DDNS 的域名）",
				Action: func(cCtx *cli.Context) error {
					color.Cyan("Welcome To Use CocoaSyncer!")
					color.Cyan("Reading Configuration...")
					readConfig(true)
					subFunction.InitializeDatabase(CocoaBasic, false)
					subFunction.CLIManageServices()
					return nil
				},
			},
			{
				Name:    "updateNodeInfo",
				Aliases: []string{"node"},
				Usage:   "更新节点信息(CocoaSyncer)",
				Action: func(cCtx *cli.Context) error {
					color.Cyan("Welcome To Use CocoaSyncer!")
					color.Cyan("Reading Configuration...")
					readConfig(true)
					subFunction.InitializeDatabase(CocoaBasic, false)
					subFunction.CLINodeManage()
					return nil
				},
			},
		},
		Copyright: "Luckykeeper <luckykeeper@luckykeeper.site | https://luckykeeper.site> | https://github.com/luckykeeper>",
	}

	if err := CocoaSyncer.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func main() {
	CocoaSyncerCLI()
}

// 程序入口
func CocoaSyncer(debugMode bool) {
	subFunction.CronTaskInit()
	color.Cyan("Welcome To Use CocoaSyncer!")
	color.Cyan("Reading Configuration...")
	readConfig(debugMode)
	if debugMode {
		log.Println("调试模式：ON!")
	}
	subFunction.InitializeDatabase(CocoaBasic, debugMode)
	subFunction.ImportNodeInfo(CocoaBasic)
	cocoaSyncerAPIService(CocoaBasic, debugMode)
}

// API 服务入口
func cocoaSyncerAPIService(CocoaBasic model.CocoaBasic, debugMode bool) {
	if !debugMode {
		gin.SetMode(gin.ReleaseMode)
	}
	cocoaSyncerRouter.CocoaSyncerRouter = gin.Default()
	if debugMode {
		cocoaSyncerRouter.GinRouter(cocoaSyncerRouter.CocoaSyncerRouter, true)
	} else {
		cocoaSyncerRouter.GinRouter(cocoaSyncerRouter.CocoaSyncerRouter, false)
	}

	// 使用反代
	if CocoaBasic.ReverseProxy {
		cocoaSyncerRouter.CocoaSyncerRouter.ForwardedByClientIP = true
	} else {
		// 不使用反代
		cocoaSyncerRouter.CocoaSyncerRouter.ForwardedByClientIP = false
	}
	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(CocoaBasic.APIPort),
		Handler: cocoaSyncerRouter.CocoaSyncerRouter,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}

	}()

	// 等待中断信号以优雅地关闭服务器（设置 15 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("等待最长15秒处理完剩余连接后关闭服务")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("服务器关闭超时，强制退出，原因:", err)
	}
	log.Println("服务成功关闭,mafumayumayumayu~~~~")
}

// read config.yaml
func readConfig(debugMode bool) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	// 校验配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalln("没有找到配置文件，请检查 ./config.yaml 是否存在！")
		} else {
			log.Fatalln("配置文件校验失败，请检查 ./config.yaml 是否存在语法错误")
		}
	}
	// 如果需要嵌套解析 yaml ，需要把 struct 也做成嵌套的
	// 参考：https://blog.csdn.net/weixin_42586723/article/details/121162029
	viper.Unmarshal(&CocoaBasic)

	log.Println("读取配置成功，当前配置如下：")
	log.Println("____________________________")
	log.Println("API 端口:", CocoaBasic.APIPort)
	log.Println("节点友好名称:", CocoaBasic.NodeName)
	log.Println("节点代号:", CocoaBasic.NodeCode)
	log.Println("工作模式:", CocoaBasic.WorkMode)
	log.Println("超时时间(分):", CocoaBasic.Timeout)
	log.Println("数据库类型:", CocoaBasic.DataBaseType)
	if debugMode {
		log.Println("通信密钥:", CocoaBasic.CocoaSecret)
		log.Println("数据库连接信息:", CocoaBasic.Dsn)
	} else {
		log.Println("通信密钥:", CocoaBasic.CocoaSecret[0:len(CocoaBasic.CocoaSecret)/8]+"**************")
		log.Println("数据库连接信息:", CocoaBasic.Dsn[0:len(CocoaBasic.Dsn)/10]+"**************")
	}
	log.Println("____________________________")
}
