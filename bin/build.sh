#########################################################################
# File Name: build.sh
# Author: lidawei
# mail: lidawei@sinoiov.com
# Created Time: Tue 22 Aug 2017 09:56:02 PM CST
#########################################################################
#!/bin/bash
#/bin/sh

parentDir()
{
	local this_dir=`pwd`
	dirname "$this_dir"
}
		
go_path=`parentDir`
# echo building on $go_path
export GOPATH=$go_path

target_name=${go_path##*/}
# echo target name $target_name

# echo building...

/opt/web_app/go/bin/go build $go_path/src/main/main.go
mv -f main $target_name
