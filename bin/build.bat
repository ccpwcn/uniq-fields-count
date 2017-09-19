@echo off

set product_name=""
for /f %%i in ("%cd%\..") do (
    set product_name=%%~ni
)

set env_path=""
for /f "delims=" %%i in ("%cd%") do (
    set env_path=%%~dpi
)

if not %product_name%=="" (
    if not %env_path%=="" (
        echo 环境检测成功，开始编译...
        set GOPATH = %env_path%
        go install main
        move /y main.exe %product_name%.exe > NUL
        echo 编译完成，程序名称：%product_name%.exe
    ) else (
        echo "获取构建环境目录失败"
    )
) else (
    echo "获取项目名称失败"
)