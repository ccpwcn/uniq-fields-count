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
        echo �������ɹ�����ʼ����...
        set GOPATH = %env_path%
        go install main
        move /y main.exe %product_name%.exe > NUL
        echo ������ɣ��������ƣ�%product_name%.exe
    ) else (
        echo "��ȡ��������Ŀ¼ʧ��"
    )
) else (
    echo "��ȡ��Ŀ����ʧ��"
)