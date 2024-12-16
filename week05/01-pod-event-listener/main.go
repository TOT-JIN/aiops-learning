package main

import (
	"flag"
	"fmt"
	"k8s.io/apimachinery/pkg/watch"
	"path/filepath"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/client-go/util/workqueue"
)

// 定义事件类型
type Event struct {
	key       string
	eventType watch.EventType
}

type Controller struct {
	indexer  cache.Indexer
	queue    workqueue.TypedRateLimitingInterface[Event]
	informer cache.Controller
}

func NewController(queue workqueue.TypedRateLimitingInterface[Event], indexer cache.Indexer, informer cache.Controller) *Controller {
	return &Controller{
		informer: informer,
		indexer:  indexer,
		queue:    queue,
	}
}

func (c *Controller) processNextItem() bool {
	event, quit := c.queue.Get()
	if quit {
		return false
	}
	defer c.queue.Done(event)

	//fmt.Printf(" queue's event key: %s\n", event.key)
	err := c.syncToStdout(event)
	c.handleErr(err, event)
	return true
}

func (c *Controller) syncToStdout(event Event) error {
	obj, exists, err := c.indexer.GetByKey(event.key)
	if err != nil {
		fmt.Printf("从存储中获取 key %s 的对象失败，错误：%v\n", event.key, err)
		return err
	}

	if !exists {
		fmt.Printf("Pod %s 已经不存在了\n", event.key)
	} else {
		pod := obj.(*corev1.Pod)
		fmt.Printf("%s 事件：Pod %s 在命名空间 %s，状态：%s\n", event.eventType, pod.Name, pod.Namespace, pod.Status.Phase)
	}
	return nil
}

func (c *Controller) handleErr(err error, event Event) {
	if err == nil {
		c.queue.Forget(event)
		return
	}

	// 如果重试次数小于 5 次，重新加入队列
	if c.queue.NumRequeues(event) < 5 {
		fmt.Printf("Retry %d for key %s\n", c.queue.NumRequeues(event), event.key)
		c.queue.AddRateLimited(event)
		return
	}

	// 超过 5 次后，放弃进一步的重试
	c.queue.Forget(event)
	fmt.Printf("Dropping pod %q out of the queue: %v\n", event.key, err)
}

func main() {
	var err error
	var config *rest.Config

	var kubeconfig *string

	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "[Optional] Absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "Absolute path to the kubeconfig file")
	}

	if config, err = rest.InClusterConfig(); err != nil {
		if config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig); err != nil {
			panic(err.Error())
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	informerFactory := informers.NewSharedInformerFactory(clientset, time.Minute*10)

	queue := workqueue.NewTypedRateLimitingQueue(workqueue.DefaultTypedControllerRateLimiter[Event]())

	podInformer := informerFactory.Core().V1().Pods()
	informer := podInformer.Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    func(obj interface{}) { onAddPod(obj, queue) },
		UpdateFunc: func(old, new interface{}) { onUpdatePod(new, queue) },
		DeleteFunc: func(obj interface{}) { onDeletePod(obj, queue) },
	})

	controller := NewController(queue, podInformer.Informer().GetIndexer(), informer)

	stopper := make(chan struct{})
	defer close(stopper)

	informerFactory.Start(stopper)
	informerFactory.WaitForCacheSync(stopper)

	go func() {
		for {
			if !controller.processNextItem() {
				break
			}
		}
	}()

	<-stopper
}

func onAddPod(obj interface{}, queue workqueue.TypedRateLimitingInterface[Event]) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err == nil {
		queue.Add(Event{key: key, eventType: watch.Added})
	}
}

func onUpdatePod(new interface{}, queue workqueue.TypedRateLimitingInterface[Event]) {
	key, err := cache.MetaNamespaceKeyFunc(new)
	if err == nil {
		queue.Add(Event{key: key, eventType: watch.Modified})
	}
}

func onDeletePod(obj interface{}, queue workqueue.TypedRateLimitingInterface[Event]) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err == nil {
		queue.Add(Event{key: key, eventType: watch.Deleted})
	}
}
