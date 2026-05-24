package main

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "syscall"
)

const imageBase = "/var/lib/minidocker/images"

// 查询本地镜像
func listImages() {
    files, err := os.ReadDir(imageBase)
    if err != nil {
        fmt.Println("无法读取镜像目录:", err)
        return
    }
    for _, f := range files {
        if f.IsDir() {
            fmt.Println(f.Name())
        }
    }
}

// 运行容器：imageName 是镜像目录名，cmdArgs 是要执行的命令
func runContainer(imageName string, cmdArgs []string) {
    imagePath := filepath.Join(imageBase, imageName, "rootfs")
    if _, err := os.Stat(imagePath); os.IsNotExist(err) {
        fmt.Println("镜像不存在:", imageName)
        return
    }

    // 重新调用自身，传入特殊参数，进入新的命名空间
    args := append([]string{"minidocker-init", imagePath}, cmdArgs...)
    cmd := exec.Command("/proc/self/exe", args...)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    // 设置新的命名空间标志
    cmd.SysProcAttr = &syscall.SysProcAttr{
        Cloneflags: syscall.CLONE_NEWUTS | // 主机名隔离
            syscall.CLONE_NEWPID |        // 进程空间隔离
            syscall.CLONE_NEWNS |         // 挂载点隔离
            syscall.CLONE_NEWNET,         // 网络隔离（也可以不加，与宿主机共享网络）
    }
    if err := cmd.Run(); err != nil {
        fmt.Println("容器退出:", err)
    }
}

// 初始化容器环境，只有被 runContainer 调用时才执行
func initContainer(rootfs string, args []string) {
    // 挂载 proc
    syscall.Mount("proc", filepath.Join(rootfs, "proc"), "proc", 0, "")
    // 切换根文件系统
    if err := syscall.Chroot(rootfs); err != nil {
        panic(err)
    }
    // 切换到根目录
    if err := os.Chdir("/"); err != nil {
        panic(err)
    }
    // 执行用户指定的命令，使用 execve 替换当前进程
    if len(args) == 0 {
        args = []string{"/bin/sh"}
    }
    if err := syscall.Exec(args[0], args, os.Environ()); err != nil {
        panic(err)
    }
}

func main() {
    // 根据第一个参数判断是普通命令还是容器初始化
    if len(os.Args) > 1 && os.Args[1] == "minidocker-init" {
        if len(os.Args) < 3 {
            fmt.Println("用法: minidocker-init <rootfs路径> [命令...]")
            os.Exit(1)
        }
        initContainer(os.Args[2], os.Args[3:])
        return
    }

    if len(os.Args) < 2 {
        fmt.Println("可用命令: images | run <镜像名> [命令]")
        return
    }

    switch os.Args[1] {
    case "images":
        listImages()
    case "run":
        if len(os.Args) < 3 {
            fmt.Println("用法: run <镜像名> [命令]")
            return
        }
        runContainer(os.Args[2], os.Args[3:])
    default:
        fmt.Println("未知命令:", os.Args[1])
    }
}