package main

import (
	"fmt"
	"time"

	//	"os"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

func main() {
	// Create a Kubernetes client using the provided kubeconfig.
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	factory := informers.NewSharedInformerFactory(clientset, 2*time.Second)
	podInformer := factory.Core().V1().Pods()

	podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod := obj.(*v1.Pod)
			fmt.Println("Pod Added", pod.Name, pod.Status.Phase)
		},
		UpdateFunc: func(oldObj, newObj interface{}) { /* handle pod updated */ },
		DeleteFunc: func(obj interface{}) { /* handle pod deleted */ },
	})
	stopCh := make(chan struct{})
	defer close(stopCh)

	// Keep the program running.
	select {}
}
