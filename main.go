package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"strings"
	"syscall"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/sys/windows"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

func main() {
	if !amAdmin() { //判断是否以管理员权限运行
		runMeElevated()
		return

	}
	err := os.Setenv("LANG", "zh_CN.UTF-8")
	if err != nil {
		log.Println("setenv lang err")
		return
	}

	myApp := app.New()
	myWindow := myApp.NewWindow("uwf-manager ")

	inputCMD := widget.NewEntry()
	inputCMD.SetPlaceHolder("...")

	inputInfo := widget.NewEntry()
	inputInfo.SetPlaceHolder("...")

	contentRun := container.NewVBox(canvas.NewText("命令 : ", color.White), inputCMD, widget.NewButton("执行", func() {
		log.Println("Content was:", inputCMD.Text)
	}))

	BaseConfig := container.New(layout.NewHBoxLayout(), canvas.NewText("基本配置 : ", color.White), widget.NewButton("查看配置", func() {
		inputCMD.SetText("uwfmgr Get-Config")
		inputInfo.SetText("查看配置详情里面有下一次启动的配置注意检查哦")

	}), widget.NewButton("启用UWF", func() {
		inputCMD.SetText("uwfmgr filter enable")
		inputInfo.SetText("启用命令可能需要执行多次，并重启电脑才有效")
	}), widget.NewButton("禁用UWF", func() {
		inputCMD.SetText("uwfmgr filter disable")
		inputInfo.SetText("禁用命令可能需要执行多次，并重启电脑才有效")
	}), layout.NewSpacer())

	SavePatch := container.New(layout.NewHBoxLayout(), canvas.NewText("写入过滤 : ", color.White), widget.NewButton("内存", func() {
		inputCMD.SetText("uwfmgr overlay Set-Type RAM")
		inputInfo.SetText("可能在禁用uwf并重启后才可以修改这个参数,保存到内存有助于保护硬盘")

	}), widget.NewButton("硬盘", func() {
		inputCMD.SetText("uwfmgr overlay Set-Type DISK")
		inputInfo.SetText("可能在禁用uwf并重启后才可以修改这个参数,保存到硬盘可以配置更大的overlay容量")
	}), layout.NewSpacer(), canvas.NewText("写入位置", color.Opaque))

	cacheSize := container.New(layout.NewHBoxLayout(), canvas.NewText("缓存设置 : ", color.White), widget.NewButton("最大缓存", func() {
		inputCMD.SetText("uwfmgr overlay set-size 10240")
		inputInfo.SetText("单位是MB，如果是内存模式建议5120以上，硬盘模式建议20480以上")

	}), widget.NewButton("警告阈值", func() {
		inputCMD.SetText("uwfmgr overlay set-warningthreshold 5024")
		inputInfo.SetText("建议为最大缓存的50-70%左右,超过这个容量后电脑会弹出一条提醒消息")
	}), widget.NewButton("严重阈值", func() {
		inputCMD.SetText("uwfmgr overlay set-criticalthreshold 8192")
		inputInfo.SetText("建议为最大缓存的80-95%左右,超过这个容量后电脑随时都可能要求重启")
	}), layout.NewSpacer(), canvas.NewText("单位MB", color.Opaque))

	DiskPartSet := container.New(layout.NewHBoxLayout(), canvas.NewText("分区保护 : ", color.White), widget.NewButton("启用C盘", func() {
		inputCMD.SetText("uwfmgr volume protect C:")
		inputInfo.SetText("建议开启uwf后并保护windows所在分区，如果你要保护的盘符不是C自行修改即可")

	}), widget.NewButton("移除C盘", func() {
		inputCMD.SetText("uwfmgr Unprotect protect C:")
		inputInfo.SetText("盘符可以自行修改")
	}), widget.NewButton("启用所有分区", func() {
		inputCMD.SetText("uwfmgr volume protect all")
		inputInfo.SetText("建议开启uwf后并保护windows所在分区，如果你要保护的盘符不是C自行修改即可")

	}), widget.NewButton("禁用所有分区", func() {
		inputCMD.SetText("uwfmgr Unprotect protect all")
		inputInfo.SetText("这样会导致uwf虽然开启了服务但是其实是失效状态")
	}))

	windowsUpdate := container.New(layout.NewHBoxLayout(), canvas.NewText("windows更新 : ", color.White), widget.NewButton("允许绕过uwf", func() {
		inputCMD.SetText("uwfmgr servicing Update-Windows")
		inputInfo.SetText("允许windows更新来更新受保护的系统")

	}), widget.NewButton("禁止绕过uwf", func() {
		inputCMD.SetText("uwfmgr servicing disable")
		inputInfo.SetText("禁止windows更新来更新受保护的系统（这条命令并不能禁用windows更新的运行只是更新后重启无效而已）")
	}), layout.NewSpacer(), canvas.NewText("", color.Opaque))

	FileExclusion := container.New(layout.NewHBoxLayout(), canvas.NewText("排除目录 : ", color.White), widget.NewButton("排除目录", func() {
		inputCMD.SetText("uwfmgr file add-exclusion  Path C:\\Users\\yh\\Downloads")
		inputInfo.SetText("重启后才有效")

	}), widget.NewButton("不再排除目录", func() {
		inputCMD.SetText("uwfmgr file remove-exclusion  Path C:\\Users\\yh\\Downloads")
		inputInfo.SetText("重启后才有效")
	}), widget.NewButton("排除文件", func() {
		inputCMD.SetText("uwfmgr file add-exclusion  Filename C:\\Users\\yh\\Downloads\\demo.txt")
		inputInfo.SetText("重启后才有效")

	}), widget.NewButton("不再排除文件", func() {
		inputCMD.SetText("uwfmgr file remove-exclusion  Filename C:\\Users\\yh\\Downloads\\demo.txt")
		inputInfo.SetText("重启后才有效")
	}), widget.NewButton("保存文件", func() {
		inputCMD.SetText("uwfmgr file commit  Filename C:\\Users\\yh\\Downloads\\demo.txt")
		inputInfo.SetText("保存一个文件的内容修改和更新")
	}), widget.NewButton("确定删除", func() {
		inputCMD.SetText("uwfmgr file commit-delete  Filename C:\\Users\\yh\\Downloads\\demo.txt")
		inputInfo.SetText("这个文件你已经删除，不希望重启后恢复")
	}), widget.NewButton("查看配置", func() {
		inputCMD.SetText("uwfmgr file get-exclusions")
		inputInfo.SetText("显示针对当前会话和下次会话的具体文件排除配置信息")
	}), layout.NewSpacer(), canvas.NewText("", color.Opaque))
	registrySet := container.New(layout.NewHBoxLayout(), canvas.NewText("注册表排除 : ", color.White), widget.NewButton("排除路径", func() {
		inputCMD.SetText("uwfmgr registry add-exclusion HKLM\\Software\\Microsoft\\Windows\\run")
		inputInfo.SetText("重启后才有效")

	}), widget.NewButton("不再排除", func() {
		inputCMD.SetText("uwfmgr registry remove-exclusion HKLM\\Software\\Microsoft\\Windows\\run")
		inputInfo.SetText("重启后才有效")
	}), widget.NewButton("保存一个注册表值", func() {
		inputCMD.SetText("uwfmgr.exe registry commit HKLM\\Software\\Test TestValue")
		inputInfo.SetText("保存一个注册表键值的修改")
	}), widget.NewButton("确定删除", func() {
		inputCMD.SetText("uwfmgr  registry commit-delete HKLM\\Software\\Test TestValue")
		inputInfo.SetText("这个键值你已经删除，不希望重启后恢复")
	}), widget.NewButton("查看配置", func() {
		inputCMD.SetText("uwfmgr registry get-exclusions")
		inputInfo.SetText("显示针对当前会话和下次会话的具体注册表排除配置信息")
	}), layout.NewSpacer(), canvas.NewText("", color.Opaque))
	text4 := canvas.NewText("---- 命令 和 提醒  ----", color.White)
	centered := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), text4, layout.NewSpacer())

	myWindow.SetContent(container.New(layout.NewVBoxLayout(), BaseConfig, SavePatch, cacheSize, DiskPartSet, windowsUpdate, FileExclusion, registrySet, centered, canvas.NewText("注意事项和提醒 : ", color.White), inputInfo, contentRun))
	myWindow.Resize(fyne.NewSize(700, 500))
	myWindow.ShowAndRun()
}

// 判断管理员权限相关
func runMeElevated() {
	verb := "runas"
	exe, _ := os.Executable()
	cwd, _ := os.Getwd()
	args := strings.Join(os.Args[1:], " ")

	verbPtr, _ := syscall.UTF16PtrFromString(verb)
	exePtr, _ := syscall.UTF16PtrFromString(exe)
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
	argPtr, _ := syscall.UTF16PtrFromString(args)

	var showCmd int32 = 1 //SW_NORMAL

	err := windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, showCmd)
	if err != nil {
		fmt.Println(err)
		// os.Exit(0) //如果重新以admin打开程序失败，可以退出当前程序
	}
}

func amAdmin() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		fmt.Println("admin no")
		return false
	}
	fmt.Println("admin yes")
	return true
}
