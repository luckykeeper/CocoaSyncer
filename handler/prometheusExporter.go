// CocoaSyncer - 心爱酱多节点智能解析平台 - Prometheus Exporter
// @CreateTime : 2023/8/30 15:10
// @LastModified : 2023/8/30 15:10
// @Author : Luckykeeper
// @Email : luckykeeper@luckykeeper.site
// @Project : CocoaSyncer

package handler

import (
	"cocoaSyncer/model"
	subFunction "cocoaSyncer/subfunction"

	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
)

// Prometheus Exporter
func PrometheusExporter(router *gin.Engine) {
	// get global Monitor object
	monitor := ginmetrics.GetMonitor()

	// +optional set metric path, default /debug/metrics
	monitor.SetMetricPath("/metrics")
	// +optional set slow time, default 5s
	monitor.SetSlowTime(5)
	// +optional set request duration, default {0.1, 0.3, 1.2, 5, 10}
	// used to p95, p99
	monitor.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})

	// cocoaSyncer 相关信息
	thisCocoa := model.CocoaBasic{ConfigImported: true}
	subFunction.CocoaDataEngine.Get(&thisCocoa)

	cocoaSyncerStatus := &ginmetrics.Metric{
		// 注意中间不能有空格
		Type:        ginmetrics.Gauge,
		Name:        "cocoaSyncerStatus",
		Description: "cocoaSyncer运行状态",
		Labels:      []string{"cocoaSyncerStatus"},
	}

	// Add metric to global monitor object
	_ = ginmetrics.GetMonitor().AddMetric(cocoaSyncerStatus)

	_ = ginmetrics.GetMonitor().GetMetric("cocoaSyncerStatus").SetGaugeValue([]string{"NodeNumber"}, float64(len(thisCocoa.OtherCocoaSyncer)))
	_ = ginmetrics.GetMonitor().GetMetric("cocoaSyncerStatus").SetGaugeValue([]string{"CocoaSyncerManagedService"}, float64(len(thisCocoa.CocoaManagedService)))
	_ = ginmetrics.GetMonitor().GetMetric("cocoaSyncerStatus").SetGaugeValue([]string{"CocoaSyncerCloudPlatformInfo"}, float64(len(thisCocoa.CloudPlatformInfo)))

	nodeOnlineCount := 0
	nodeOfflineCount := 0
	for _, node := range thisCocoa.OtherCocoaSyncer {
		if node.StatusCode == 200 {
			nodeOnlineCount++
		} else {
			nodeOfflineCount++
		}
	}

	_ = ginmetrics.GetMonitor().GetMetric("cocoaSyncerStatus").SetGaugeValue([]string{"CocoaSyncerNodeOnlineCount"}, float64(nodeOnlineCount))
	_ = ginmetrics.GetMonitor().GetMetric("cocoaSyncerStatus").SetGaugeValue([]string{"CocoaSyncerNodeOfflineCount"}, float64(nodeOfflineCount))

	// set middleware for gin
	monitor.Use(router)
}
