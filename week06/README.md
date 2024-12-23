# 补全chatgpt.go，实现deleteResource方法

> 需补充两段代码：
> 
> * callFunction方法中调用deleteResource前，解析大模型返回的json model，并传参
> * 实现deleteResource方法

## callFunction
```go
	if name == "deleteResource" {
		params := struct {
			Namespace    string `json:"namespace"`
			ResourceType string `json:"resource_type"`
			ResourceName string `json:"resource_name"`
		}{}
		if err := json.Unmarshal([]byte(arguments), &params); err != nil {
			return "", fmt.Errorf("failed to parse function call name=%s arguments=%s", name, arguments)
		}
		return deleteResource(params.Namespace, params.ResourceType, params.ResourceName)
	}
```

## deleteResource

```go
func deleteResource(namespace, resourceType, resourceName string) (string, error) {
	clientGo, err := utils.NewClientGo(kubeconfig)
	resourceType = strings.ToLower(resourceType)
	var gvr schema.GroupVersionResource
	switch resourceType {
	case "deployment":
		gvr = schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	case "service":
		gvr = schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}
	case "pod":
		gvr = schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}
	default:
		return "", fmt.Errorf("unsupported resource type: %s", resourceType)
	}
	// 先获取一下资源看看是否存在
	resourceRe, err := clientGo.DynamicClient.Resource(gvr).Namespace(namespace).Get(context.TODO(), resourceName, metav1.GetOptions{})

	if err != nil {
		return fmt.Sprintf("资源 %s 已不存在", resourceRe), fmt.Errorf("对应资源已不存在: %w", err)
	}
	// 执行删除操作
	if err := clientGo.DynamicClient.Resource(gvr).Namespace(namespace).Delete(context.TODO(), resourceName, metav1.DeleteOptions{}); err != nil {
		return fmt.Sprintf("资源 %s 删除失败", resourceRe), fmt.Errorf("删除失败: %w", err)
	}
	return fmt.Sprintf("%s 删除成功", resourceRe.GetName()), nil
}
```

## 执行测试
* 执行前，资源已生成
```bash
kubectl get deploy nginx-deployment
NAME               READY   UP-TO-DATE   AVAILABLE   AGE
nginx-deployment   3/3     3            3           2m38s 
```
* 执行后，资源已删除
```bash
./k8scopilot ask chatgpt           
我是 K8s Copilot，有什么可以帮助你：
> 帮我删除一下default NS下，名为nginx-deployment的deploy
nginx-deployment 删除成功
```
```bash
kubectl get deploy nginx-deployment
Error from server (NotFound): deployments.apps "nginx-deployment" not found
```
* 再次尝试执行删除
```bash
> 帮我删除一下default NS下，名为nginx-deployment的deploy
Error calling function: 对应资源已不存在: deployments.apps "nginx-deployment" not found 
```