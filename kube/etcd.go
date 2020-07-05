package kube

import (
	"context"
	"encoding/json"
	"k8s_announce/conf"
	"k8s_announce/utils"

	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/clientv3"
)

/*PullK8sClusterListFromEtcd 从etcd拉取cluster list信息*/
func PullK8sClusterListFromEtcd() {
	cli, err := clientv3.New(conf.EtcdConf)
	if err != nil {
		utils.Log.WithFields(logrus.Fields{
			"error": err,
		}).Error("[PullK8sClusterListFromEtcd]etcd client new failed.")
		return
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), conf.EetcdRequestTimeout)
	resp, err := cli.Get(ctx, conf.EetcdClusterListPath)
	cancel()
	if err != nil {
		utils.Log.WithFields(logrus.Fields{
			"error": err,
		}).Error("[PullK8sClusterListFromEtcd]etcd get k8s cluster list failed.")
		return
	}
	for _, ev := range resp.Kvs {
		if string(ev.Key) == conf.EetcdClusterListPath {
			utils.Log.WithFields(logrus.Fields{
				"key":   string(ev.Key),
				"value": string(ev.Value),
			}).Info("[PullK8sClusterListFromEtcd]get k8s cluster list conf.")

			var clusterListTmp []ClusterInfo
			err := json.Unmarshal(ev.Value, &clusterListTmp)
			if err != nil {
				utils.Log.WithFields(logrus.Fields{
					"error": err,
				}).Error("[PullK8sClusterListFromEtcd]cluster list value unmarshaled failed.")
				return
			}
			clusterList = clusterListTmp
			for _, v := range clusterList {
				utils.Log.WithFields(logrus.Fields{
					"name":   v.Name,
					"domain": v.Domain,
					"token":  v.Token,
				}).Info("[PullK8sClusterListFromEtcd]cluster list details.")
			}
		}
	}
}