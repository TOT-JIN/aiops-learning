#!/bin/bash

# 检查依赖
if ! command -v curl &>/dev/null; then
  echo "Error: curl 未安装。请先安装 curl。" >&2
  exit 1
fi

# 环境变量设置
export OLLAMA_HOST=0.0.0.0:8080
export OLLAMA_MODELS=/data/ollama/models

# 安装 Ollama CLI
if [ ! -f /usr/local/bin/ollama ]; then
  curl -fsSL https://ollama.com/install.sh | sh
else
  echo "Ollama 已安装"
fi

# 配置模型存储目录
sudo mkdir -p /data/ollama/models
sudo useradd -r -s /bin/false -U -m -d /data/ollama ollama || echo "用户 ollama 已存在"
sudo usermod -a -G ollama $(whoami)
sudo chown -R ollama:ollama /data/ollama
sudo chmod -R 755 /data/ollama

# 配置服务
service_content="[Unit]
Description=Ollama Service
After=network-online.target

[Service]
Environment=\"OLLAMA_HOST=0.0.0.0:8080\"
Environment=\"OLLAMA_MODELS=/data/ollama/models\"
ExecStart=/usr/local/bin/ollama serve
User=ollama
Group=ollama
Restart=always
RestartSec=3
Environment=\"PATH=\$PATH\"

[Install]
WantedBy=default.target
"

file_path="/etc/systemd/system/ollama.service"
echo "$service_content" | sudo tee "$file_path" > /dev/null
sudo chmod 644 "$file_path"

# 启用服务
sudo systemctl daemon-reload
sudo systemctl enable ollama
sudo systemctl restart ollama

# 检查服务状态
sudo systemctl status ollama || {
  echo "Error: Ollama 服务启动失败。" >&2
  exit 1
}

## 预加载模型
#if ! command -v ollama &>/dev/null; then
#  echo "Error: Ollama CLI 未正确安装。" >&2
#  exit 1
#fi
#
#echo "Ollama 预加载模型"
#OLLAMA_HOST=0.0.0.0:8080 ollama pull qwen2:0.5b || {
#  echo "Error: 模型预加载失败。" >&2
#  exit 1
#}
#
#echo "Ollama 安装和配置完成。"
#