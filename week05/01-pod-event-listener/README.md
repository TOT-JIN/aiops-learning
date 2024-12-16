# 使用 Informer + RateLimitingQueue 监听 Pod 事件
> 仿照课程中代码编写脚本使用 Informer + RateLimitingQueue 监听 Pod 事件
> 
> 根据pod的不同事件类型做不同的队列处理，以便输出信息时更好的操作
> 
> 使用 deployment 部署pod，并进行扩缩容以及更新。
> 

## 定义事件类型结构体，并添加到队列中

* 定义事件类型结构体，封装事件类型和对应的pod的key
```golang
type Event struct {
	key       string
	eventType watch.EventType
}
```

* 任何需要定义workqueue.TypedRateLimitingInterface[T]类型的queue的地方指定 T 为 Event

```golang
type Controller struct {
	indexer  cache.Indexer
	queue    workqueue.TypedRateLimitingInterface[Event]
	informer cache.Controller
}
```

* 各EventHandler方法中，实例化Event，并添加到队列中
```golang
func onAddPod(obj interface{}, queue workqueue.TypedRateLimitingInterface[Event]) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err == nil {
		queue.Add(Event{key: key, eventType: watch.Added})
	}
}
```

* 对于 queue 的相关方法，如 Add, AddRateLimited, Done, Forget, Get, NumRequeues 等，传参皆为 Event 类型。
* 对于 indexer.GetByKey 方法，传参为pod的key，此处为从queue中取出类型为Event的元素的key字段
```golang
func (c *Controller) syncToStdout(event Event) error {
    obj, exists, err := c.indexer.GetByKey(event.key)
    ...
} 
```

## 运行观察输出
* 创建并删除一个简单的deployment
```bash
$ kubectl create deployment nginx --image nginx
ADDED 事件：Pod nginx-8f458dc5b-hn6zf 在命名空间 default，状态：Pending
MODIFIED 事件：Pod nginx-8f458dc5b-hn6zf 在命名空间 default，状态：Pending
MODIFIED 事件：Pod nginx-8f458dc5b-hn6zf 在命名空间 default，状态：Pending
MODIFIED 事件：Pod nginx-8f458dc5b-hn6zf 在命名空间 default，状态：Running 

$ kubectl delete deploy nginx
MODIFIED 事件：Pod nginx-8f458dc5b-hn6zf 在命名空间 default，状态：Running
MODIFIED 事件：Pod nginx-8f458dc5b-hn6zf 在命名空间 default，状态：Running
MODIFIED 事件：Pod nginx-8f458dc5b-hn6zf 在命名空间 default，状态：Running
Pod default/nginx-8f458dc5b-hn6zf 已经不存在了

```






