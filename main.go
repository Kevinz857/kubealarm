package main

import (
	"time"

	"kubealarm/conf"
	"kubealarm/kube"
	"kubealarm/utils"
	//"kubealarm/utils"
)

func init() {
	Log := utils.InitLogger(conf.LogPath, time.Hour*24*7, time.Hour*24)
	Log.Info("[init]logger init successfully")
}

func main() {
	kube.PullK8sClusterListFromEtcd()
	kube.InitClientset()
}
