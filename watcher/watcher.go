package main

import (
	"fmt"
	"time"

	//	"os"

	"strings"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

type PodInfo struct {
  nodeName string
  containerId string
}
type Info struct {
   // map of node name to lightfoot pod ip on that node
   nodes  map[string] string
   // map of pod names to container_id
   // These are all pods which have lightfoot:enable
   container   map[string] PodInfo
}

type PodEvent int
// Enum for Pod updates
const (
	Add PodEvent = iota
	Update = iota
	Delete = iota
)

func printPodInfo(pod *v1.Pod) {
	fmt.Println(pod.Spec.NodeName, pod.Status.PodIP, pod.ObjectMeta.GetName(), pod.ObjectMeta.GetLabels())
	for _, containerStatus := range pod.Status.ContainerStatuses {
                containerID := containerStatus.ContainerID
                if strings.HasPrefix(containerID, "cri-o://") {
                        containerID = containerID[8:] 
                }
                fmt.Printf("Container: %s, ID: %s\n", containerStatus.Name, containerID)
        }

}

func handlePodEvent(e PodEvent, pod *v1.Pod){
   switch e {

   }
}

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
			handlePodEvent(Add, pod)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			pod := oldObj.(*v1.Pod)
			handlePodEvent(Update, pod)
		},
		DeleteFunc: func(obj interface{}) {
			//Do nothing we don't care about deletes for now
		},
	})
	stopCh := make(chan struct{})
	defer close(stopCh)

	factory.Start(stopCh)
	// Keep the program running.
	select {}
}

