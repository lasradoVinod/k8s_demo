package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/homedir"
)

var (
	kubeconfig *string
)

func main() {
	// Parse command-line flags to specify the kubeconfig file.
	kubeconfig = flag.String("kubeconfig", filepath.Join(homedir.HomeDir(), ".kube", "config"), "absolute path to the kubeconfig file")
	flag.Parse()

	// Create a Kubernetes client using the provided kubeconfig.
	config, err := rest.InClusterConfig()
	if err != nil {
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// Set up a pod informer to watch pods with the label "lightfoot:enable".
	podInformer := cache.NewListWatchFromClient(
		clientset.CoreV1().RESTClient(),
		"pods",
		"",
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				pod := obj.(*v1.Pod)
				labels := pod.GetLabels()
				if val, ok := labels["lightfoot:enable"]; ok && val == "true" {
					writePodNameToFile(pod.Name)
				}
			},
		},
	)

	podInformerController := cache.NewController(
		cache.NewConfig(),
		podInformer,
		cache.ResourceEventHandlerFuncs{},
	)

	stopCh := make(chan struct{})
	defer close(stopCh)

	go podInformerController.Run(stopCh)

	// Keep the program running.
	select {}
}

func writePodNameToFile(podName string) {
	fileName := "pods_with_lightfoot_enable.txt"
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(fmt.Sprintf("%s\n", podName)); err != nil {
		fmt.Println("Error writing to file:", err)
	}
}
