# windows编译
需要下载交叉编译器，建议：
LLVM MinGW
MSYS2
WinLibs

安装 c++编译工具
下载包
https://github.com/niXman/mingw-builds-binaries/releases
解压后，将其下的bin目录，添加到环境变量path


并设置环境变量
set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=1  

然后编译
go build -ldflags="-H windowsgui"

# windows 运行
安装 WebView2（服务器推荐）：下载Evergreen Standalone Installer（离线安装包，x64 版本）：
https://developer.microsoft.com/zh-cn/microsoft-edge/webview2/?form=MA13LH#download https://developer.microsoft.com/zh-cn/microsoft-edge/webview2/#download-section
以管理员身份运行安装程序，并重启服务器（重要！）

最新的windows系统默认安装了，不需要单独安装，老版本的需要单独安装


