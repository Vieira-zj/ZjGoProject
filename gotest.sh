#!/bin/bash
set -ex
echo "gotest.sh"

ZJ_GOPRJ="${HOME}/Workspaces/zj_go_project"
# if current golang project is not in system path
# GOPATH=${ZJ_GOPRJ}:${GOPATH}

function go_test_main() {
    # go run src/demo.tests/main/main.go -h
    # go run src/demo.tests/main/main.go -n arg1 arg2 arg3
    go run src/demo.tests/main/main.go -n -s "," arg1 arg2 arg3
}

# GO TEST
function go_test_help() {
    go help test
}

# compile the test binary to pkg.test but do not run it.
# the file name can be changed with the -o flag.
function go_compile_test() {
    go test -c
}

# go tests, root_dir = $GOPATH
function go_func_test() {
    go test -v src/demo.tests/gotests/
    # go test -v src/demo.tests/gotests/word_test.go
    # go test -v -run="TestEcho" demo.tests/gotests/
}

# go tests, code coverage
function go_coverage_test() {
    # go test -v -run="IsPalindrome" -cover -coverprofile=c.out demo.tests/gotests/
    # go test -v -cover -coverprofile=c.out demo.tests/gotests/
    go tool cover -html=c.out
}

# go tests, benchmark
function go_benchmark_test() {
    # go test -v -bench=. src/demo.tests/gotests/word_bmark_test.go
    go test -v -bench=. -benchmem src/demo.tests/gotests/word_bmark_test.go
}

# go unit tests
function go_func_unittest() {
    local test_dir="${GOPATH}/src/demo.tests/gotests"
    local target="calculation"
    go test -v ${test_dir}/${target}.go ${test_dir}/${target}_test.go
}

function ring_unittest() {
    local test_dir="${GOPATH}/src/demo.hello/sort"
    local target="struct_ring"
    go test -v ${test_dir}/${target}.go ${test_dir}/${target}_test.go
}


# TOOLS TEST
function tool_svc_test() {
    local test_dir=$1
    if [ -z $2 ]; then
        go test -v -count=1 src/tools.app/services/${test_dir}
    else
        test_file=$2
        go test -v -count=1 src/tools.app/services/${test_dir}/${test_file}
    fi
}

function tool_utils_test() {
    if [ -z $1 ]; then
        go test -v src/tools.app/utils
    else
        test_file=$1
        go test -v src/tools.app/utils/${test_file}
    fi
}

function tool_utils_benchmark_test() {
    go test -v -bench=. -benchmem src/tools.app/utils/$1
}


# MOCK TEST
function mock_server_test() {
    if [ -z $1 ]; then
        go test -v src/tools.app/utils
    else
        test_file=$1
        go test -v src/mock.server/handlers/${test_file}
    fi
}


# BDD TEST
# ginkgo: http://onsi.github.io/ginkgo/
# gomega: http://onsi.github.io/gomega/
function go_bdd_test_01() {
    # bddtest="${ZJ_GOPRJ}/bin/ginkgo"
    # ginkgo -v -focus="test.asserter.suite02" src/demo.tests/bddtests/
    # ginkgo -v -focus="suite03.case04" src/demo.tests/bddtests/
    ginkgo -v --focus="suite04.case02" src/demo.tests/bddtests/ -- -myFlag="ginkgo test"
}

function go_bdd_test_02() {
    ginkgo -v -focus="test.share.suite12" src/demo.tests/bddtests/
    # ginkgo -v -focus="suite13.context01" src/demo.tests/bddtests/
    # ginkgo -v -focus="suite11.case11" src/demo.tests/bddtests/
}

function go_bdd_benchmark_test() {
    ginkgo -v --focus="suite14.measure01" src/demo.tests/bddtests/
}


# SHELL TEST
# EX01, field check
# https://blog.csdn.net/longyinyushi/article/details/50728049

# EX01-00, number comparison
function shell_test() {
    echo "current dir: $(pwd)"
    
    # remove leading spaces => sed ‘s/^[ \t]*//g'
    count=`ls | wc -l | sed "s/^[ \t]*//g"`
    if [ ${count} -gt 0 ]; then
        echo "file count: ${count}"
    fi
    
    # string comparison
    name="zhengjin"
    if [ ${name} == "$(whoami)" ]; then
        echo "cur user is zhengjin."
    fi
}

# EX01-01, field exist check
function shell_test_0101() {
    is_exist="test"
    if [ ${is_exist} ]; then
        echo "[] check, is exist."
    fi
    if [[ "${is_exist}" ]]; then
        echo "[[]] check, is exist."
    fi
}

# EX01-02, field length check
function shell_test_0102() {
    empty_str=""
    if [ -z $empty_str ]; then
        echo "string length = 0."
    else
        echo "string length > 0."
    fi
}

# EX01-03, field empty check
function shell_test_0103() {
    test_str="test"
    if [[ -n $test_str ]]; then
        echo 'string not empty.'
    else
        echo 'string empty.'
    fi
}

# EX01-04, file exist check
function shell_test_0104() {
    test_path="./c.out"
    
    if [ -f ${test_path} ]; then
        echo "file ${test_path} exist."
    else
        echo "file ${test_path} NOT exist."
    fi
    
    while [ ! -f ${test_path} ]; do
        echo 'checking file ${test_path} ...';sleep 3
    done
    echo "file ${test_path} exist."
    
    for (( i = 0; i < 10; i++ )); do
        echo 'checking file ${test_path} ...';sleep 3
        if [ -f ${test_path} ]; then
            echo "file ${test_path} exist."
            break
        fi
    done
}

# EX02-01, array
function shell_test_0201() {
    tmp_array1=("ele1")
    tmp_array2=("ele2" "ele3")
    tmp_array3=(${tmp_array1[@]} ${tmp_array2[@]})
    
    for ele in ${tmp_array3[@]}; do
        echo $ele
    done
    echo ${tmp_array3[@]}  # echo all elements
    echo "array length: ${#tmp_array3[@]}"
}

# EX02-02
function shell_test_0202() {
    focus_pkg=()
    temp_pkg1=("a1" "a2")
    focus_pkg=(${focus_pkg[@]} ${temp_pkg1[@]}) # append
    echo ${focus_pkg[@]}
    temp_pkg2=("a3" "a4" "a5")
    focus_pkg=(${focus_pkg[@]} ${temp_pkg2[@]})
    echo ${focus_pkg[@]}
}

# EX02-03
function shell_test_0203() {
    focus_pkg=("a1" "a2" "a3" "a4")
    skip_pkg="a1,a2,a3,a4,a5,a6,a7"
    for v in ${focus_pkg[*]}; do
        skip_pkg=${skip_pkg/${v},/} # replace {$v}, with ""
    done
    echo "focus packages => ${focus_pkg[*]}"
    echo "skip packages => ${skip_pkg}"
}

# EX03, if-else with regexp
function shell_test_03() {
    node_name="go1.9_fix"
    if [[ ($node_name =~ "go1.9") && ($node_name =~ "fix") ]]; then
        echo 'version check ok.'
    else
        echo 'version should be go1.9 with fix.'
    fi
}

# EX04, ${var} usage
function shell_test_04() {
    tmp_file="/dir1/dir2/dir3/my.file.txt"
    # sub string start 0, len = 5
    echo ${tmp_file:0:5}
    # sub string start 5, len = 5
    echo ${tmp_file:5:5}
    # replace first "dir" with "path"
    echo ${tmp_file/dir/path}
    # replace all "dir" with "path"
    echo ${tmp_file//dir/path}
}

# EX05, +=
function shell_test_05() {
    src="hello"
    src=${src}" world"
    echo ${src}
    src="test, ${src}"
    echo ${src}
}

# EX05-01, calculation
function shell_test_0501() {
    i=5
    ((i++))
    ret=$((i+5*2))
    echo "results: $ret"
}

# EX06, iterator
function shell_test_06() {
    for i in $(seq 1 3); do
        echo "index ${i}"
    done
    
    for (( i = 0; i < 3; i++ )); do
        echo "index ${i}"
    done
    
    for f in $(ls ~/Downloads/tmp_files/test.*); do
        echo "test file: ${f}"
    done
}

# EX07, run download parallel
function shell_test_07() {
    url="http://7zkl9d.com1.z1.glb.clouddn.com/slowResponse"
    for (( i=0; i<20; i++)); do
        echo "run at: $i"
        curl -v ${url} -x iovip-z1.qbox.me:80 > /dev/null &
        sleep 2s
    done
    sleep 15m
    ps -ef | grep "curl" | grep -v "grep" | awk '{print $2}' | xargs kill -9
}

# EX08, custom functions
echoEnv() { echo "TEST_ENV=$TEST_ENV"; echo "TEST_ZONE=$TEST_ZONE";}
setEnv() { export TEST_ENV=$1; echo "TEST_ENV=$TEST_ENV";}
setZone() { export TEST_ZONE=$1; echo "TEST_ZONE=$TEST_ZONE";}

findStr() { grep "$1" ./*;}
findStrAll() { grep -r "$1" ./;}

function shell_test_08() {
    echoEnv
    # findStrAll "search_text"
}

# EX09, shell exit with ret code
function shell_test_09() {
    echo "test exit with error code 1."
    exit 1
}


# MAIN
# go_test_main

# go_test_help
# go_func_test
# go_benchmark_test

go_func_unittest
# ring_unittest

# tool_svc_test diskusage filestree_test.go
# tool_utils_test mails_test.go
# tool_utils_benchmark_test ioutil_bmark_test.go

# mock_server_test mockdemo_test.go

# go_bdd_test_01
# go_bdd_test_02
# go_bdd_benchmark_test

# shell_test
# shell_test_0203

echo "golang test DONE."
set +ex
