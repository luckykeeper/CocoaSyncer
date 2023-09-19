// CocoaSyncer - 心爱酱多节点智能解析平台 - 数据库初始化
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

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
)

var CocoaDataEngine *xorm.Engine

// 初始化数据库
func InitializeDatabase(CocoaBasic model.CocoaBasic, debugMode bool) {
	// 防止空指针问题，这样声明 err
	// 见：https://stackoverflow.com/questions/56396386/xorm-example-not-working-runtime-error-invalid-memory-address-or-nil-pointer
	var err error
	CocoaDataEngine, err = xorm.NewEngine(CocoaBasic.DataBaseType, CocoaBasic.Dsn)
	if err != nil {
		log.Fatalln("数据库初始化失败：", err)
	}
	// 不使用缓存，小项目没有必要使用缓存
	// cacher := caches.NewLRUCacher(caches.NewMemoryStore(), 1000)
	// CocoaDataEngine.SetDefaultCacher(cacher)

	// 不需要 Close ，orm 会自己判断
	// defer CocoaDataEngine.Close()

	if debugMode {
		CocoaDataEngine.ShowSQL(true)
	}
	err = CocoaDataEngine.Ping()
	if err != nil {
		log.Panicln("数据库连接失败！检查数据库信息是否正确，程序返回原因为：", err)
	} else {
		CocoaDataEngine.SetConnMaxLifetime(time.Second * 60) // 最大连接存活时间
		CocoaDataEngine.SetMaxOpenConns(100)                 // 最大连接数
		CocoaDataEngine.SetMaxIdleConns(3)                   // 最大空闲连接数
		log.Println("连接到远程数据库成功！")
	}
	log.Println("同步数据表中……")
	err = CocoaDataEngine.Sync(new(model.CocoaBasic))
	if err != nil {
		log.Fatalln("同步数据表失败：", err)
	}
	log.Println("数据库同步成功！")
}

// 写入本节点信息
func ImportNodeInfo(CocoaBasic model.CocoaBasic) {
	log.Println("判断数据库内是否已经写入配置数据……")
	testConfig := new(model.CocoaBasic)
	_, err := CocoaDataEngine.Cols("configImported").Get(testConfig)
	if err != nil {
		log.Fatalln("判断数据库内是否已经写入配置数据失败:", err)
	}
	// log.Println(testConfig.ConfigImported)
	if !testConfig.ConfigImported {
		log.Println("初次启动，写入配置文件到数据库……")
		CocoaBasic.ConfigImported = true

		// 同时写入节点数据
		thisCocoaSyncerInfo := &model.NodeInfo{
			NodeName:    CocoaBasic.NodeName,
			NodeCode:    CocoaBasic.NodeCode,
			NodeAddress: CocoaBasic.NodeAddress,
			WorkMode:    CocoaBasic.WorkMode,
			CocoaSecret: CocoaBasic.CocoaSecret,
		}
		CocoaBasic.OtherCocoaSyncer = append(CocoaBasic.OtherCocoaSyncer, *thisCocoaSyncerInfo)

		CocoaDataEngine.Insert(CocoaBasic)
	} else {
		log.Println("数据库存在配置，跳过初始化流程")
	}
}
