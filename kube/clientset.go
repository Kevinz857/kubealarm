// Copyright 2019 Kube Capacity Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kube

import (
	"flag"
	"os"
	"path/filepath"
	"sync"

	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	// Required for GKE, OIDC, and more
	"kubealarm/utils"
	//_ "k8s.io/client-go/plugin/pkg/client/auth"
)

var clientset *kubernetes.Clientset
var globalMutex sync.RWMutex

/*InitClientset 初始化*/
func InitClientset() {
	globalMutex.Lock()
	for _, v := range clusterList {
		var cluster *ClusterInfo = new(ClusterInfo)
		cluster.Name = v.Name

		var cfg *rest.Config = new(rest.Config)
		cfg.Host, cfg.BearerToken = "https://"+v.Domain, v.Token

		var tls rest.TLSClientConfig
		tls.Insecure = true

		cfg.TLSClientConfig = tls

		var err error
		cluster.Clientsets, err = kubernetes.NewForConfig(cfg)
		if err != nil {
			utils.Log.WithFields(logrus.Fields{
				"ClusterName": v.Name,
				"error":       err,
			}).Error("[InitClientset] k8s config new failed.")
			continue
		}

		AddNodeInformerForCluster(cluster)
	}
	globalMutex.Unlock()
}

/*InitClientsetV1 初始化clientset*/
func InitClientsetV1() *kubernetes.Clientset {
	if flag.Lookup("kubeconfig") != nil {
		return clientset
	}
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", filepath.Join("etc", "kubernetes", "kubeconfig"), "absolute path to the kubeconfig file")
	}
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	utils.Log.Info("[InitClientset] clientset initial done.")
	return clientset
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
