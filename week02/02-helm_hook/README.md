# Helm Hook使用
> 练习以下四个hook
> - pre-install-hook
> - post-install-hook
> - pre-upgrade-hook
> - post-delete-hook

## 自定义信息和格式化
通过 Helm 的 `values.yaml` 文件来让这些信息更具动态性。将输出信息放在 `values.yaml` 中，方便修改或重用。

#### 在 `values.yaml` 中定义：

```yaml
hookMessages:
  preInstall: "Starting installation..."
  postInstall: "Installation complete!"
  preUpgrade: "Starting upgrade..."
  postDelete: "Cleanup after deletion!"
```

#### 在 Hook 中引用 `values.yaml`：

```yaml
command:
  - sh
  - -c
  - |
    echo "Helm Hook | pre-install | {{ .Release.Name }} | {{ .Values.hookMessages.preInstall }}"
```

#### 执行结果：

先后执行 安装、升级、卸载 操作

```bash
#执行install操作
[root@k8s-master01 demo]# helm install nginx . -n nginx --create-namespace
NAME: nginx
LAST DEPLOYED: Fri Sep 20 15:35:27 2024
NAMESPACE: nginx
STATUS: deployed
REVISION: 1
NOTES:
1. Get the application URL by running these commands:
  export POD_NAME=$(kubectl get pods --namespace nginx -l "app.kubernetes.io/name=demo,app.kubernetes.io/instance=nginx" -o jsonpath="{.items[0].metadata.name}")
  export CONTAINER_PORT=$(kubectl get pod --namespace nginx $POD_NAME -o jsonpath="{.spec.containers[0].ports[0].containerPort}")
  echo "Visit http://127.0.0.1:8080 to use your application"
  kubectl --namespace nginx port-forward $POD_NAME 8080:$CONTAINER_PORT

#已执行pre-install-hook和post-install-hook
[root@k8s-master01 ~]# kubectl -n nginx get pods -o wide
NAME                            READY   STATUS      RESTARTS   AGE   IP           NODE           NOMINATED NODE   READINESS GATES
nginx-demo-757c8bd58-njzlx      1/1     Running     0          55s   10.244.5.7   k8s-worker03   <none>           <none>
nginx-post-install-hook-lf6mw   0/1     Completed   0          55s   10.244.4.5   k8s-worker02   <none>           <none>
nginx-pre-install-hook-pxsxm    0/1     Completed   0          78s   10.244.5.6   k8s-worker03   <none>           <none>

#查看输出
[root@k8s-master01 ~]# kubectl -n nginx logs pod/nginx-pre-install-hook-pxsxm
Helm Hook | pre-install | nginx | Starting installation...
[root@k8s-master01 ~]# kubectl -n nginx logs pod/nginx-post-install-hook-lf6mw
Helm Hook | post-install | nginx | Installation complete!
```

```bash
#执行upgrade操作
[root@k8s-master01 demo]# helm upgrade --install nginx . -n nginx
Release "nginx" has been upgraded. Happy Helming!
NAME: nginx
LAST DEPLOYED: Fri Sep 20 15:39:30 2024
NAMESPACE: nginx
STATUS: deployed
REVISION: 2
NOTES:
1. Get the application URL by running these commands:
  export POD_NAME=$(kubectl get pods --namespace nginx -l "app.kubernetes.io/name=demo,app.kubernetes.io/instance=nginx" -o jsonpath="{.items[0].metadata.name}")
  export CONTAINER_PORT=$(kubectl get pod --namespace nginx $POD_NAME -o jsonpath="{.spec.containers[0].ports[0].containerPort}")
  echo "Visit http://127.0.0.1:8080 to use your application"
  kubectl --namespace nginx port-forward $POD_NAME 8080:$CONTAINER_PORT

#pre-upgrade-hook已执行
[root@k8s-master01 ~]# kubectl -n nginx get pods -o wide
NAME                            READY   STATUS      RESTARTS   AGE     IP           NODE           NOMINATED NODE   READINESS GATES
nginx-demo-757c8bd58-njzlx      1/1     Running     0          4m25s   10.244.5.7   k8s-worker03   <none>           <none>
nginx-post-install-hook-lf6mw   0/1     Completed   0          4m25s   10.244.4.5   k8s-worker02   <none>           <none>
nginx-pre-install-hook-pxsxm    0/1     Completed   0          4m48s   10.244.5.6   k8s-worker03   <none>           <none>
nginx-pre-upgrade-hook-7l6vs    0/1     Completed   0          46s     10.244.4.6   k8s-worker02   <none>           <none>

#查看输出
[root@k8s-master01 ~]# kubectl -n nginx logs pod/nginx-pre-upgrade-hook-7l6vs
Helm Hook | pre-upgrade | nginx | Starting upgrade...
```

```bash
#执行uninstall操作
[root@k8s-master01 demo]# helm uninstall nginx -n nginx
release "nginx" uninstalled

#demo负载已删除，新增post-delete-hook
[root@k8s-master01 ~]# kubectl -n nginx get pods -o wide
NAME                            READY   STATUS      RESTARTS   AGE     IP           NODE           NOMINATED NODE   READINESS GATES
nginx-post-delete-hook-slrbg    0/1     Completed   0          49s     10.244.6.5   k8s-worker04   <none>           <none>
nginx-post-install-hook-lf6mw   0/1     Completed   0          6m16s   10.244.4.5   k8s-worker02   <none>           <none>
nginx-pre-install-hook-pxsxm    0/1     Completed   0          6m39s   10.244.5.6   k8s-worker03   <none>           <none>
nginx-pre-upgrade-hook-7l6vs    0/1     Completed   0          2m37s   10.244.4.6   k8s-worker02   <none>           <none>

#查看输出
[root@k8s-master01 ~]# kubectl -n nginx logs pod/nginx-post-delete-hook-slrbg
Helm Hook | post-delete | nginx | Cleanup after deletion!
```

