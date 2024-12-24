# CronHPA 添加 ConfigMap

> 在 CronHPASpec 中添加一个新的字段 ConfigMapData，用于存储 ConfigMap 的键值对
> 
> 在 Reconcile 中添加 ConfigMap 的处理逻辑，当 ConfigMap 发生变化时，更新 CronHPASpec 中的 ConfigMapData 字段
> 

## 修改 cronhpa_types.go

```go
type CronHPASpec struct {
    ScaleTargetRef ScaleTargetReference `json:"scaleTargetRef"`
    Jobs           []JobSpec            `json:"jobs"`

    // ConfigMapData holds the key-value pairs for the ConfigMap
    ConfigMapData map[string]string `json:"configMapData,omitempty"`
}

```

## 修改 cronhpa_controller.go

* 定义 `createOrUpdateConfigMap` 方法，用于 ConfigMap 的创建和更新

```go
func (r *CronHPAReconciler) createOrUpdateConfigMap(ctx context.Context, cronhpa *autoscalingv1.CronHPA) error {
	log := log.FromContext(ctx)
	configMapName := fmt.Sprintf("%s-config", cronhpa.Name)

	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: cronhpa.Namespace,
		},
		Data: cronhpa.Spec.ConfigMapData,
	}

	err := r.Client.Create(ctx, configMap)
	if err != nil {
		if errors.IsAlreadyExists(err) {
			existingConfigMap := &corev1.ConfigMap{}
			if err := r.Client.Get(ctx, types.NamespacedName{Name: configMapName, Namespace: cronhpa.Namespace}, existingConfigMap); err != nil {
				log.Error(err, "Failed to get existing ConfigMap", "configMap", configMapName)
				return err
			}
			existingConfigMap.Data = cronhpa.Spec.ConfigMapData
			if err := r.Client.Update(ctx, existingConfigMap); err != nil {
				log.Error(err, "Failed to update ConfigMap", "configMap", configMapName)
				return err
			}
		} else {
			log.Error(err, "Failed to create ConfigMap", "configMap", configMapName)
			return err
		}
	}
	return nil
}
```

* 在 `Reconcile` 中添加 ConfigMap 的处理调用
```go
	var cronhpa autoscalingv1.CronHPA
	if err := r.Get(ctx, req.NamespacedName, &cronhpa); err != nil {
		if errors.IsNotFound(err) {
			log.Info("CronHPA resource not found. Ignoring since object must be deleted")
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	// 调用创建或更新 ConfigMap 的函数
	if err := r.createOrUpdateConfigMap(ctx, &cronhpa); err != nil {
		return reconcile.Result{}, err
	}

```

## 运行测试

* 编辑 `cronhpa/config/samples/autoscaling_v1_cronhpa.yaml`，定义声明变量
* 执行CRD安装
```bash
$ make manifests
~/cronhpa/bin/controller-gen rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases
$ make install
~/cronhpa/bin/controller-gen rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases

~/cronhpa/bin/kustomize build config/crd | kubectl apply -f -
customresourcedefinition.apiextensions.k8s.io/cronhpas.autoscaling.aiops.com created
$ kubectl api-resources|grep cronhpa
cronhpas                                       autoscaling.aiops.com/v1               true         CronHPA
```
* 运行operator
```bash
$ kubectl get deploy nginx #已有一个副本
NAME    READY   UP-TO-DATE   AVAILABLE   AGE
nginx   1/1     1            1           18s
$ make run
~/cronhpa/bin/controller-gen rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases
~/cronhpa/bin/controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."
2024-12-24T21:10:06+08:00       INFO    setup   starting manager
2024-12-24T21:10:06+08:00       INFO    starting server {"name": "health probe", "addr": "[::]:8081"}
2024-12-24T21:10:06+08:00       INFO    Starting EventSource    {"controller": "cronhpa", "controllerGroup": "autoscaling.aiops.com", "controllerKind": "CronHPA", "source": "kind source: *v1.CronHPA"}
2024-12-24T21:10:06+08:00       INFO    Starting Controller     {"controller": "cronhpa", "controllerGroup": "autoscaling.aiops.com", "controllerKind": "CronHPA"}
2024-12-24T21:10:06+08:00       INFO    Starting workers        {"controller": "cronhpa", "controllerGroup": "autoscaling.aiops.com", "controllerKind": "CronHPA", "worker count": 1}
2024-12-24T21:10:06+08:00       INFO    Reconciling CronHPA     {"controller": "cronhpa", "controllerGroup": "autoscaling.aiops.com", "controllerKind": "CronHPA", "CronHPA": {"name":"cronhpa-sample","namespace":"default"}, "namespace": "default", "name": "cronhpa-sample", "reconcileID": "3c8e974a-e5d1-4e60-b4c2-536d10b676c2"}
2024-12-24T21:10:06+08:00       INFO    Job info        {"controller": "cronhpa", "controllerGroup": "autoscaling.aiops.com", "controllerKind": "CronHPA", "CronHPA": {"name":"cronhpa-sample","namespace":"default"}, "namespace": "default", "name": "cronhpa-sample", "reconcileID": "3c8e974a-e5d1-4e60-b4c2-536d10b676c2", "name": "scale-up", "lastRunTime": "0001-01-01 00:00:00 +0000 UTC", "nextScheduledTime": "0001-01-01T00:01:00Z", "now": "2024-12-24T21:10:06+08:00"}
2024-12-24T21:10:06+08:00       INFO    Updating deployment replicas    {"controller": "cronhpa", "controllerGroup": "autoscaling.aiops.com", "controllerKind": "CronHPA", "CronHPA": {"name":"cronhpa-sample","namespace":"default"}, "namespace": "default", "name": "cronhpa-sample", "reconcileID": "3c8e974a-e5d1-4e60-b4c2-536d10b676c2", "name": "nginx", "targetSize": 3}
2024-12-24T21:10:07+08:00       INFO    Successfully updated deployment replicas        {"controller": "cronhpa", "controllerGroup": "autoscaling.aiops.com", "controllerKind": "CronHPA", "CronHPA": {"name":"cronhpa-sample","namespace":"default"}, "namespace": "default", "name": "cronhpa-sample", "reconcileID": "3c8e974a-e5d1-4e60-b4c2-536d10b676c2", "deployment": {"name":"nginx","namespace":"default"}, "replicas": 3}
2024-12-24T21:10:07+08:00       INFO    Requeue after   {"controller": "cronhpa", "controllerGroup": "autoscaling.aiops.com", "controllerKind": "CronHPA", "CronHPA": {"name":"cronhpa-sample","namespace":"default"}, "namespace": "default", "name": "cronhpa-sample", "reconcileID": "3c8e974a-e5d1-4e60-b4c2-536d10b676c2", "time": "53.175338s"}
2024-12-24T21:10:07+08:00       INFO    Reconciling CronHPA     {"controller": "cronhpa", "controllerGroup": "autoscaling.aiops.com", "controllerKind": "CronHPA", "CronHPA": {"name":"cronhpa-sample","namespace":"default"}, "namespace": "default", "name": "cronhpa-sample", "reconcileID": "b35c6e3f-2d10-41db-a084-d9d905c5a141"}
2024-12-24T21:10:07+08:00       INFO    Job info        {"controller": "cronhpa", "controllerGroup": "autoscaling.aiops.com", "controllerKind": "CronHPA", "CronHPA": {"name":"cronhpa-sample","namespace":"default"}, "namespace": "default", "name": "cronhpa-sample", "reconcileID": "b35c6e3f-2d10-41db-a084-d9d905c5a141", "name": "scale-up", "lastRunTime": "2024-12-24 21:10:06 +0800 CST", "nextScheduledTime": "2024-12-24T21:11:00+08:00", "now": "2024-12-24T21:10:07+08:00"}
2024-12-24T21:10:07+08:00       INFO    Requeue after   {"controller": "cronhpa", "controllerGroup": "autoscaling.aiops.com", "controllerKind": "CronHPA", "CronHPA": {"name":"cronhpa-sample","namespace":"default"}, "namespace": "default", "name": "cronhpa-sample", "reconcileID": "b35c6e3f-2d10-41db-a084-d9d905c5a141", "time": "52.785667s"}
$ kubectl get deploy nginx
NAME    READY   UP-TO-DATE   AVAILABLE   AGE
nginx   3/3     3            3           42s   #副本增大到3个
#查看是否已挂载configmap volume
$ kubectl get deploy nginx -o yaml|grep -C3 config-volume
        terminationMessagePolicy: File
        volumeMounts:
        - mountPath: /etc/config
          name: config-volume
      dnsPolicy: ClusterFirst
      restartPolicy: Always
      schedulerName: default-scheduler
--
      - configMap:
          defaultMode: 420
          name: cronhpa-sample-config
        name: config-volume                         # 已挂载 configmap
```
