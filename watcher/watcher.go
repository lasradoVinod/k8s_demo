package main

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
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
	nodeName    string
	containerId []string
}

type ControllerInfo struct {
	// map of node name to lightfoot pod ip on that node
	nodes map[string]string
	// map of pod names to container_id
	// These are all pods which have lightfoot:enable
	container map[string]PodInfo
	// List of containers lightfoot is not yet updated about
	pending map[string]bool
	// Channel to trigger new pending
	wake chan struct{}
	// lock for adding to pending state
	mu sync.Mutex
}

type PodEvent int

// Enum for Pod updates
const (
	Add    PodEvent = iota
	Update          = iota
	Delete          = iota
)

var controllerInfo ControllerInfo

func (c *ControllerInfo) SendUpdate(name string) bool {
	var ip string
	var ok bool
	if ip, ok = c.nodes[c.container[name].nodeName]; !ok {
		return false
	}
	var buffer bytes.Buffer
	for _, str := range c.container[name].containerId {
		buffer.WriteString(str)
		buffer.WriteString("\n") // Add a newline if you want separation
	}
	req, err := http.NewRequest("POST", "https://"+ip+":12000/crio-id", bytes.NewReader(buffer.Bytes()))
	if err != nil {
		fmt.Println("Error sending request:", err)
		return false
	}
	req.Header.Set("Content-Type", "text/plain")
	// Create an HTTP client and perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode < 0 || resp.StatusCode > 300 {
		fmt.Println("Request failed with status:", buffer, resp.Status)
		return false
	}
	return true
}

func (c *ControllerInfo) ContactLightroom() {
	timer := time.NewTimer(10 * time.Second)

	for {
		select {
		case <-timer.C:
		case <-c.wake:
			c.mu.Lock()
			for k, _ := range c.pending {
				status := c.SendUpdate(k)
				if status {
					delete(c.pending, k)
				}
			}
			c.mu.Unlock()
		}
	}
}

func handlePodEvent(e PodEvent, pod *v1.Pod) {
	switch e {
	case Add:
		done := false
		// Add the ip to the lightfoot container corresponding to nodename
		name := pod.ObjectMeta.GetName()
		if strings.HasPrefix(name, "lightfoot-daemon") {
			controllerInfo.nodes[pod.Spec.NodeName] = pod.Status.PodIP
			break
		}
		for k, v := range pod.ObjectMeta.GetLabels() {
			if k == "lightfoot" && v == "enable" {
				done = true
				break
			}
		}
		if !done {
			break
		}
		var p PodInfo
		p.nodeName = pod.Spec.NodeName
		for _, containerStatus := range pod.Status.ContainerStatuses {
			p.containerId = append(p.containerId, strings.TrimPrefix(containerStatus.ContainerID, "cri-o://"))
		}
		controllerInfo.container[name] = p
		controllerInfo.mu.Lock()
		controllerInfo.pending[name] = false
		controllerInfo.mu.Unlock()
		controllerInfo.wake <- struct{}{}
	case Update:
	case Delete:
		// Handle these cases later revisions
	default:
		break
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

	controllerInfo.wake = make(chan struct{})

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
