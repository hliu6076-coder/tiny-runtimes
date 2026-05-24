# tiny-runtimes
## ---最小容器实现
路径：/var/lib/minidocker/images/
```bash
# 1. 创建目标存放目录
sudo mkdir -p /var/lib/minidocker/images/alpine/rootfs
# 2. 进入临时目录，准备下载
cd /tmp
# 3. 下载 Alpine minirootfs
sudo wget https://dl-cdn.alpinelinux.org/alpine/v3.21/releases/x86_64/alpine-minirootfs-3.21.11-x86_64.tar.gz
# 将下载的压缩包解压到你的镜像仓库
sudo tar -xzf /tmp/alpine-minirootfs-3.21.11-x86_64.tar.gz -C /var/lib/minidocker/images/alpine/rootfs
```
#### 1. 定义镜像仓库路径
```go
const imageBase = "/var/lib/minidocker/images"
```

#### 2. 编译与运行
```bash
go mod init tiny-runtimes   # 初始化模块（若需要）
go build -o minidocker .    # 编译
sudo ./minidocker images    # 查询镜像
sudo ./minidocker run alpine /bin/sh  # 启动容器
```