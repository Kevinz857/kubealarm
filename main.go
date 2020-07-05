package main

import (
	"time"

	"k8s_announce/conf"
	"k8s_announce/kube"
	"k8s_announce/utils"
)

func init() {
	Log := utils.InitLogger(conf.LogPath, time.Hour*24*7, time.Hour*24)
	Log.Info("[init]logger init successfully")
}

func main() {
	kube.PullK8sClusterListFromEtcd()
	kube.InitClientset()
}
