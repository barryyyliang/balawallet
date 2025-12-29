@echo off

cd tool

:: set1：compile tools.go，
echo  "compile tools.go..."
go build -o tools.exe tools.go
if errorlevel 1 (
    echo "error：compile tools.go error！"
    pause
    exit /b 1
)

::  步骤2：执行 tools.exe（Windows 原生写法，无需 ./）
echo "execute tools.exe..."
tools.exe
if errorlevel 1 (
    echo "error： execute tools.exe error"
    pause
    exit /b 1
)

:: 步骤3：仅当确认需要切换目录时才执行 cd ..（建议注释/删除，除非明确需要）
:: 警告：请确认是否真的需要跳到上级目录编译，否则注释下面这行！
cd ..

:: 步骤4：编译主程序，修正 ldflags 写法
echo "compile wallet"
go build -ldflags=-H=windowsgui 
if errorlevel 1 (
    echo "error: compile wallet error"
    pause
    exit /b 1
)

echo " finished compile！"

move .\balawallet.exe  .\bin\
