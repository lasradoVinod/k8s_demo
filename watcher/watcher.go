package main

import (
	"bytes"
	"fmt"
	"log"
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
		log.Println("lightfoot instance not found for node ", c.container[name].nodeName)
		return false
	}
	var buffer bytes.Buffer
	for _, str := range c.container[name].containerId {
		buffer.WriteString(str)
		buffer.WriteString("\n") // Add a newline if you want separation
	}
	req, err := http.NewRequest("POST", "http://"+ip+":12000/crio-id", bytes.NewReader(buffer.Bytes()))
	if err != nil {
		fmt.Println("Error sending request:", err)
		return false
	}
	req.Header.Set("Content-Type", "text/plain")
	// Create an HTTP client and perform the request
	log.Println("Sending ", c.container[name].nodeName, ip, " pod ", name, buffer)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 300 {
		fmt.Println("Request failed with status:", buffer, resp.Status)
		return false
	}
	return true
}

func (c *ControllerInfo) ManagePending() {
	c.mu.Lock()
	for k := range c.pending {
		if c.SendUpdate(k) {
			delete(c.pending, k)
		}
	}
	c.mu.Unlock()
}

func (c *ControllerInfo) ContactLightfoot() {
	timer := time.NewTicker(10 * time.Second)
	for {
		select {
		case _ = <-timer.C:
			c.ManagePending()
		case <-c.wake:
			c.ManagePending()
		}
	}
}

func CheckPodLabels(labels map[string]string) bool {
	for k, v := range labels {
		if k == "lightfoot" && v == "enable" {
			return true
		}
	}
	return false
}

func AddPod(name string, p PodInfo) {
	controllerInfo.container[name] = p
	controllerInfo.mu.Lock()
	controllerInfo.pending[name] = true
	controllerInfo.mu.Unlock()
	controllerInfo.wake <- struct{}{}
}

func handlePodEvent(e PodEvent, pod *v1.Pod) {
	switch e {
	case Add:
		// Add the ip to the lightfoot container corresponding to nodename
		name := pod.ObjectMeta.GetName()
		log.Println("Adding ", name)
		if strings.HasPrefix(name, "lightfoot-daemon") {
			controllerInfo.nodes[pod.Spec.NodeName] = pod.Status.PodIP
			break
		}
		if !CheckPodLabels(pod.ObjectMeta.GetLabels()) {
			break
		}
		if _, ok := controllerInfo.container[name]; ok {
			break
		}
		var p PodInfo
		p.nodeName = pod.Spec.NodeName
		for _, containerStatus := range pod.Status.ContainerStatuses {
			p.containerId = append(p.containerId, strings.TrimPrefix(containerStatus.ContainerID, "cri-o://"))
		}
		AddPod(name, p)
	case Delete:
		name := pod.ObjectMeta.GetName()
		if strings.HasPrefix(name, "lightfoot-daemon") {
			// Delete lightfoot-ip
			nodeName := pod.Spec.NodeName
			delete(controllerInfo.nodes, nodeName)
			// Add all pods on this node to pending
			controllerInfo.mu.Lock()
			for k, v := range controllerInfo.container {
				fmt.Println(v.nodeName, nodeName)
				if v.nodeName == nodeName {
					controllerInfo.pending[k] = true
				}
			}
			controllerInfo.mu.Unlock()
			break
		}
		if CheckPodLabels(pod.ObjectMeta.GetLabels()) {
			controllerInfo.mu.Lock()
			delete (controllerInfo.container, name)
			delete (controllerInfo.pending, name)
			controllerInfo.mu.Unlock()
		}

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

	controllerInfo.container = make(map[string]PodInfo)
	controllerInfo.pending = make(map[string]bool)
	controllerInfo.nodes = make(map[string]string)
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
			handlePodEvent(Add, pod)
		},
		DeleteFunc: func(obj interface{}) {
			//Do nothing we don't care about deletes for now
		},
	})

	go controllerInfo.ContactLightfoot()
	stopCh := make(chan struct{})
	defer close(stopCh)

	factory.Start(stopCh)
	// Keep the program running.
	select {}
}
