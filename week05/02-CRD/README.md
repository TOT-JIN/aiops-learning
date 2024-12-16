# CRD && dynamicClient

## 定义AIOps

文件：*crd.yaml*
* Group：aiops.geektime.com
* Version: v1alpha1
* Kind: AIOps
* 复数: aiopses
* 单数: aiops
* 简短名称: aiops
* 以命名空间隔离
* 定义两个字段
  * classCount: 已有课堂数
  * isCompleted: 是否已完结

## 定义aiops资源

文件：*aiops.yaml*
* 命名空间：aiops
* 名称：aiops-course
* 16节课
* 课程已完结

## 使用dynamicClient获取该资源

文件：*main.go*
* 命令行参数解析
* k8s集群连接配置加载
* Dynamic Client 创建：Dynamic client 允许用户访问一个动态的 API 路径，而不需要预先定义的客户端接口。
* REST Mapper 获取：使用 Kubernetes Discovery API 获取所有可用的 API 资源分组，并创建一个 RESTMapper。用于将 GVK 对象映射到 GVR 。
* Kind 转换为 GVR：提取通过命令行传入的资源种类（Kind），并使用 RESTMapper 将其映射到相应的 GVR 。
* 资源获取：使用上一步得到的 GVR 和 Dynamic client 来访问特定的 Kubernetes 资源。
* 资源输出

## 创建自定义类型资源，并获取
* 创建aiops资源类型，并查看
```bash
$ kubectl apply -f crd.yaml
customresourcedefinition.apiextensions.k8s.io/aiopses.aiops.geektime.com created
$ kubectl api-resources|grep aiopses
aiopses               aiops        aiops.geektime.com/v1alpha1            true         AIOps
```
* 创建aiops资源，并查看
```bash
$ kubectl create ns aiops
namespace/aiops created 
$ kubectl apply -f aiops.yaml
aiops.aiops.geektime.com/aiops-course created
$ kubectl -n aiops get aiops -o yaml
apiVersion: v1
items:
- apiVersion: aiops.geektime.com/v1alpha1
  kind: AIOps
  metadata:
    annotations:
      kubectl.kubernetes.io/last-applied-configuration: |
        {"apiVersion":"aiops.geektime.com/v1alpha1","kind":"AIOps","metadata":{"annotations":{},"name":"aiops-course","namespace":"aiops"},"spec":{"classCount":16,"isCompleted":true}}
    creationTimestamp: "2024-12-16T10:06:35Z"
    generation: 1
    name: aiops-course
    namespace: aiops
    resourceVersion: "303790"
    uid: b34cac0c-fc3b-4fab-a555-911b3ab9d478
  spec:
    classCount: 16
    isCompleted: true
kind: List
metadata:
  resourceVersion: ""
```

* 使用dynamicClient获取aiops资源，并输出
```bash
$ go run main.go get aiops
课程: aiops-course, 课堂数: 16, 是否完结: true
```