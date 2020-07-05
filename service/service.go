package service

import (
	"context"

	"kubealarm/conf"
	"kubealarm/utils"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
)

/*GetNodePods 获取指定节点pod*/
func GetNodePods(clientset *kubernetes.Clientset, nodeName string) string {
	var podList []string

	pods, err := clientset.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{
		FieldSelector: "spec.nodeName=" + nodeName,
	})
	if err != nil {
		panic(err.Error())
	}
	for _, v := range pods.Items {
		podList = append(podList, v.Name)
	}
	var infoStr string = conf.InfoMesg
	infoStr = conf.InfoMesgEnviro + nodeName + conf.InfoMesg
	for _, v := range podList {
		v += "\n"
		infoStr = infoStr + v
	}
	utils.Log.Info("[GetNodePods] infoMesg is ", infoStr)
	return infoStr
}
