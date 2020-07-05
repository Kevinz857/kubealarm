package conf

import (
	"time"

	"go.etcd.io/etcd/clientv3"
)

const (
	/*InfoMesgEnviro IM消息头*/
	InfoMesgEnviro string = "【生产环境】【重要】："
	/*InfoMesg IM消息体*/
	InfoMesg string = " 机器宕机维护，涉及到以下实例会自动迁移到其他机器，请知晓：\n"

	/*GracePeriodSeconds 实例销毁时间*/
	GracePeriodSeconds int64 = 300

	/*EetcdClusterListPath ETCD集群信息获取路径*/
	EetcdClusterListPath string = "/k8s-cluster-list"

	/*EetcdRequestTimeout 请求超时时间*/
	EetcdRequestTimeout time.Duration = 10 * time.Second

	/*K8sNodesURL 获取nodelist url*/
	K8sNodesURL string = "/api/v1/nodes"

	/*InformerTimeout k8s api超时时间*/
	InformerTimeout time.Duration = 10 * time.Second

	/*LogPath 日志目录*/
	LogPath string = "./"

	/*InfoLogFileName info日志文件名*/
	InfoLogFileName string = "kube-alarm-info.log"

	/*ErrLogFileName err日志文件名*/
	ErrLogFileName string = "kube-alarm-error.log"

	/*ImLen IM消息默认长度*/
	ImLen int = 1000
)

/*EtcdConf 定义etcd默认server*/
var EtcdConf clientv3.Config = clientv3.Config{
	Endpoints: []string{"127.0.0.1:2379"},
	//Timeout:   2 * time.Second,
	//UserName: "root",
	//Password: "root",
}
