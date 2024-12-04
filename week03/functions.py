import json
'''
预定义 function
'''
def modify_config(service_name, key, value):
    print(f"exec:  sed -i 's/--{key}=[^ ]*/--{key}={value}/' {service_name}.service")
    return json.dumps({service_name: service_name, key: value})


def restart_service(service_name):
    print(f"exec:  systemctl restart {service_name}.service")
    return json.dumps({service_name: service_name})

def apply_manifest(resource_type, image):
    print(f"exec:  kubectl create {resource_type} {image}-{resource_type} --image={image}")
    return json.dumps({resource_type: image})
tools = [
    {
        "type": "function",
        "function": {
            "name": "modify_config",
            "description": "修改某个服务的某个配置项",
            "parameters": {
                "type": "object",
                "properties": {
                    "service_name": {
                        "type": "string",
                        "description": "服务名称，例如：帮我修改 server01 的配置，则对应服务名称为 server01",
                    },
                    "key": {
                        "type": "string",
                        "description": "配置项，例如：key1 修改为 value1，则对应配置项为 key1",
                    },
                    "value": {
                        "type": "string",
                        "description": "配置值，例如：key1 修改为 value1，则对于配置值为 value1",
                    }
                },
                "required": ["service_name", "key", "value"]
            }
        }
    },
    {
        "type": "function",
        "function": {
            "name": "restart_service",
            "description": "重启某个服务",
            "parameters": {
                "type": "object",
                "properties": {
                    "service_name": {
                        "type": "string",
                        "description": "服务名称，例如：帮我重启 server01 服务，则对应服务名称为 server01",
                    }
                },
                "required": ["service_name"]
            }
        }
    },
    {
        "type": "function",
        "function": {
            "name": "apply_manifest",
            "description": "选择想要的应用类型，结合应用参数来部署对应类型的应用",
            "parameters": {
                "type": "object",
                "properties": {
                    "resource_type": {
                        "type": "string",
                        "description": "资源类型，例如：帮我部署一个 resourceA，则对应的资源类型为 resourceA",
                    },
                    "image": {
                        "type": "string",
                        "description": "镜像名称，例如：帮我部署一个 resouceA，镜像是 resouce_a，则镜像名称为 resouce_a",
                    }
                },
                "required": ["resource_type", "image"]
            }
        }
    }
]