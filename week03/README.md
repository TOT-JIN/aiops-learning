# 实践 Function Calling

## 修改配置:
```bash
请输入您想操作的内容: 帮我修改 gateway 的配置，vendor 修改为 alipay

ChatGPT want to call function:  [ChatCompletionMessageToolCall(id='call_RWNvbFHTMLTOimdLt5ZoFpk6', function=Function(arguments='{"key":"vendor","service_name":"gateway","value":"alipay"}', name='modify_config', parameters=None), type='function')]
Function name: modify_config, arguments: {'key': 'vendor', 'service_name': 'gateway', 'value': 'alipay'}
exec:  sed -i 's/--vendor=[^ ]*/--vendor=alipay/' gateway.service
Function exec result: {"gateway": "gateway", "vendor": "alipay"}
LLM Res: 已成功将 gateway 的配置中的 vendor 修改为 alipay。
```

## 重启服务
```bash
请输入您想操作的内容: 帮我重启 gateway 服务

ChatGPT want to call function:  [ChatCompletionMessageToolCall(id='call_MqbazMup4FTlMiFbusgyEyDK', function=Function(arguments='{"service_name":"gateway"}', name='restart_service', parameters=None), type='function')]
Function name: restart_service, arguments: {'service_name': 'gateway'}
exec:  systemctl restart gateway.service
Function exec result: {"gateway": "gateway"}
LLM Res: 我已经成功地重启了 `gateway` 服务。如果还有其他需要帮助的地方，请告诉我！
```

## 部署服务
```bash
请输入您想操作的内容: 帮我部署一个 deployment，镜像是 nginx

ChatGPT want to call function:  [ChatCompletionMessageToolCall(id='call_ICIAUsBWEWnuGnKEIMGrCvT8', function=Function(arguments='{"resource_type":"deployment","image":"nginx"}', name='apply_manifest', parameters=None), type='function')]
Function name: apply_manifest, arguments: {'resource_type': 'deployment', 'image': 'nginx'}
exec:  kubectl create deployment nginx-deployment --image=nginx
Function exec result: {"deployment": "nginx"}
LLM Res: nginx Deployment 已成功部署。如果还有其他需要，请告诉我！
```