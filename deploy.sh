#!/bin/bash
set -e
if ! command -v docker &> /dev/null; then
    echo "❌ 未找到 Docker，请先安装 Docker。"
    echo ""
    echo "Ubuntu/Debian:  sudo apt install docker.io"
    echo "CentOS/RHEL:    sudo yum install docker"
    echo "Arch Linux:     sudo pacman -S docker"
    echo "macOS:          访问 https://docs.docker.com/get-docker/"
    echo ""
    echo "安装后请确保 docker 服务已启动：sudo systemctl start docker"
    exit 1
fi
docker build -t tiny-runtimes:latest -f minidocker/dockerfile minidocker
exec docker run --rm -it --privileged tiny-runtimes:latest "$@"