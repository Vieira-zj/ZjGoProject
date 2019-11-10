#!/bin/bash
set -ex
# -x, show run commands with arguments
# -e, tell bash exit script if any statement returns a non-true value
echo "gorun.sh"

# ENV VAR SET
# source $QBOXROOT/kodo/env.sh
# source $QBOXROOT/base/env.sh
ZJ_GOPRJ="${HOME}/Workspaces/zj_go_project"
# if current golang project is not in system path
# GOPATH=${ZJ_GOPRJ}:${GOPATH}


# GO MAIN
# Go learn doc: https://github.com/gopl-zh/gopl-zh.github.com.git
# Go fmt: https://github.com/golang/go/wiki/CodeReviewComments
# Effective Go: https://golang.org/doc/effective_go.html
if [ -z $1 ]; then
    go run src/demo.hello/main/main.go
    exit 0
fi

if [[ $1 == "main" ]]; then
    go run src/demo.hello/main/main.go -args hello world
    # go run src/demo.hello/main/main.go -period 3s
    # go run src/demo.hello/main/main.go -h
    # go run src/demo.hello/main/main.go -p 7890 -c 404
    exit 0
fi

if [[ $1 == "app" ]]; then
    # go run src/demo.app/main/main.go
    go run src/tools.app/apps/k8sio/main.go
    exit 0
fi

if [[ $1 == "utest" ]]; then
    go run src/tools.app/apps/utilstest/main.go
    exit 0
fi


# BUILD TOOLS BIN
function scp_remote() {
    local bin_path="$1"
    local remote_ip="10.200.20.21"
    ping ${remote_ip} -c 1
    if [ $? == 0 ]; then
        cd ${ZJ_GOPRJ}/src/mock.server/main
        scp ${bin_path} qboxserver@${remote_ip}:~/zhengjin/ && rm ${bin_path}
    fi
}

function build_tools_bin() {
    local target=$1
    local main_dir="${ZJ_GOPRJ}/src/tools.app/apps/${target}"
    local bin_path="${HOME}/Downloads/tmp_files/${target}"
    if [[ $2 == "linux" ]]; then
        GOOS=linux GOARCH=amd64 go build -o ${bin_path} ${main_dir}/main.go
    else
        go build -o ${bin_path} ${main_dir}/main.go
    fi
    # scp_remote ${bin_path}
}

if [[ $1 == "tool" ]]; then
    build_tools_bin $2
    # build_tools_bin $2 "linux"
    exit 0
fi

if [[ $1 = "httprouter" ]]; then
    build_tools_bin $1
    cp -r ${ZJ_GOPRJ}/src/tools.app/services/httptemplate/templates ${HOME}/Downloads/tmp_files
    exit 0
fi

if [[ $1 = "grpc" ]]; then
    target=$2 # route_guide
    main_dir="${ZJ_GOPRJ}/src/tools.app/apps/grpc/${target}"
    bin_path="${HOME}/Downloads/tmp_files/${target}"
    go build -o ${bin_path}/server ${main_dir}/server/main.go
    go build -o ${bin_path}/client ${main_dir}/client/main.go
    if [[ -d ${main_dir}/testdata ]]; then
        cp -r ${main_dir}/testdata ${bin_path}
    fi
    exit 0
fi


# BUILD MOCK BIN  ./gorun.sh mock [linux|arm]
function go_build_bin() {
    local target_bin="$1"
    cd ${ZJ_GOPRJ}/src/mock.server/main
    if [ $2 ]; then
        GOOS=linux GOARCH=$2 go build -o ${target_bin} main.go
    else
        go build -o ${target_bin} main.go
    fi
    
    local target_dir="${HOME}/Downloads/tmp_files/mockserver"
    # scp_remote ${target_bin}; scp_remote mock_conf.json
    mv ${target_bin} ${target_dir} && cp mock_conf.json ${target_dir}
}

function build_mock_bin {
    if [[ $1 == "linux" ]]; then
        go_build_bin "${mock_bin}_$1" "amd64"
        return
    fi
    if [[ $1 == "arm" ]]; then
        go_build_bin "${mock_bin}_$1" "arm"
        return
    fi
    go_build_bin "${mock_bin}_mac"
}

mock_bin="mockserver"
if [[ $1 == "mock" ]]; then
    build_mock_bin $2
fi

set +ex # set configs off
