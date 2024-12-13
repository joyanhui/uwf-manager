uwf-manager 
 windows下的uwf配置工具
 基于 golang+fyne
 
 因为网络上前人写的已经停止更新，导致无法使用。而一直不怎么擅长记录win命令的我，还是决定写一个小工具来用。
 
 基于golang没有外部依赖，只是一个命令提醒工具。


## 使用
windows 下先启用 uwf 功能。以 Windows 11 IoT 企业版 LTSC 24H2 为例
1  win+r 运行 control 打开控制面板
2 在控制面板中点击程序
3 选择启用或关闭 Windows 功能
4 启用 设备锁定的子项 统一写入筛选器
5 点确定

然后下载 uwf-manager.exe 用管理员权限启动
## 下载地址
临时下载地址

https://github.com/joyanhui/uwf-manager/raw/refs/heads/main/uwf-manager.exe



## FAQ
### 部分参数的修改说明
修改某些参数需要先关闭uwf，然后重启windows系统，修改完成后再启用uwf。

## 构建 - windows 下

1.Install Go
过程略

2.Install GCC 
推荐 msys2 然后 在 msys2 里安装 gcc
```sh
pacman -Syu
pacman -S git mingw-w64-x86_64-toolchain mingw-w64-x86_64-go
echo "export PATH=\$PATH:~/Go/bin" >> ~/.bashrc
```
然后 I:\msys64\mingw64\bin 这个目录加入到windows的环境变量path中

3.

go get fyne.io/fyne/v2
go get fyne.io/fyne/v2/internal/svg@v2.5.2
go get fyne.io/fyne/v2/internal/painter@v2.5.2
go get fyne.io/fyne/v2/storage/repository@v2.5.2
go get fyne.io/fyne/v2/internal/app@v2.5.2
go get fyne.io/fyne/v2/lang@v2.5.2
go get fyne.io/fyne/v2/widget@v2.5.2
go get fyne.io/fyne/v2/internal/driver/glfw@v2.5.2
go get fyne.io/fyne/v2/app
fyne.io/fyne/v2/internal/metadata

go get github.com/fyne-io/image@v0.0.0
go get github.com/go-gl/glfw/v3.3/glfw@v0.0.0

go install fyne.io/fyne/v2/cmd/fyne@latest


## 打包
fyne package -os windows -icon myapp.png
