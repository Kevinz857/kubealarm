package kube

import (
	"kubealarm/conf"
	"kubealarm/service"
	"kubealarm/utils"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

/*ClusterInfo 用于存储集群信息*/
type ClusterInfo struct {
	Name       string `json:"name"`
	Domain     string `json:"domain"`
	Token      string `json:"token"`
	Clientsets *kubernetes.Clientset
}

var clusterList []ClusterInfo

/*AddNodeInformerForCluster nodeinforer实现*/
func AddNodeInformerForCluster(cluster *ClusterInfo) {
	utils.Log.Info("[AddNodeInformerForCluster]start init nodes informer.")
	//node static attributes init
	nodeStopper := make(chan struct{})
	defer close(nodeStopper)
	// init informer
	factory := informers.NewSharedInformerFactory(cluster.Clientsets, conf.InformerTimeout)
	nodeInformer := factory.Core().V1().Nodes()
	informer := nodeInformer.Informer()
	// start informer , list && watch
	defer runtime.HandleCrash()

	// 从 apiserver 同步资源，即 list
	go factory.Start(nodeStopper)
	//go nodeSharedInformer.Start(nodeStopCh)

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: nil,
		UpdateFunc: func(oldObj interface{}, newObj interface{}) {
			onode, nnode := oldObj.(*corev1.Node), newObj.(*corev1.Node)
			oldStatus, newStatus := onode.Status.Conditions[len(onode.Status.Conditions)-1].Status, nnode.Status.Conditions[len(nnode.Status.Conditions)-1].Status
			if oldStatus == "True" && newStatus == "Unknown" {
				utils.Log.Warn("[onUpdate] disater node is", nnode.Name)
				infoMesg := service.GetNodePods(cluster.Clientsets, nnode.Name)
				var infoMesgI string
				if len([]rune(infoMesg)) > conf.ImLen {
					for imBeg := 0; imBeg <= len([]rune(infoMesg)); imBeg += conf.ImLen {
						if imBeg >= conf.ImLen {
							infoMesgI = infoMesg[conf.ImLen:]
						} else {
							infoMesgI = infoMesg[imBeg : imBeg+conf.ImLen]
						}
						//im.SendDiscussMessageService(infoMesgI)
						utils.Log.Info("[onUpdate] infoMesg_i is ", infoMesgI)
					}
				}
			}
		},
		DeleteFunc: nil,
	})

	if !cache.WaitForCacheSync(nodeStopper, informer.HasSynced) {
		utils.Log.WithFields(logrus.Fields{
			"ClusterName": cluster.Name,
		}).Error("[AddNodeInformerForCluster]timeout waiting for caches nodes to sync.")
		return
	}

	//var nodeLabel labels.Selector
	nodeLister := nodeInformer.Lister()
	nodes, err := nodeLister.List(labels.Everything())
	if err != nil {
		utils.Log.WithFields(logrus.Fields{
			"ClusterName": cluster.Name,
		}).Error("[AddNodeInformerForCluster]get nodes from nInformer cache failed.")
		return
	}

	utils.Log.WithFields(logrus.Fields{
		"nodesNum": len(nodes),
	}).Info("[AddNodeInformerForCluster]init nodes.")

	utils.Log.Info("[AddNodeInformerForCluster]end init nodes info.")
	<-nodeStopper
}
