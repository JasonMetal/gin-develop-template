#!/bin/bash
show_usage="args:[-e: develop,bvt,master] [--name:Your Project Name] [-p: Port]"
# 环境
env=""
# 项目名
project=""
# 端口
port1=0
port2=0
# 协议
protocol="http"
# 仓库地址
repository_url="testdocker-registry.web-site.com:5000"

#获取参数
while [ -n "$1" ]
do
        case "$1" in
                -e) env=$2; shift 2;;
                --name) project=$2; shift 2;;
                -p) port1=$2; shift 2;;
                --) break ;;
                *) echo -e "unsupported arg:$1,$2","\n$show_usage"; break ;;
        esac
done

if [[ -z $env || -z $project || -z $port1 ]]; then
        echo "$show_usage"
        exit 0
fi

runEnv=$env
if [ "$env" == "develop" ]; then
  runEnv="test"
fi
if [ "$env" == "master" ]; then
  runEnv="prod"
fi
if [ "$env" == "release" ]; then
  runEnv="bvt"
fi

port2=`expr $port1 + 1`

git checkout "$env" -f
git pull
git submodule update --remote --init
git submodule foreach --recursive "git checkout $env -f && git pull"

# 构建程序
test="export GO111MODULE=on && CGO_ENABLED=1 GOOS=linux GOARCH=amd64; go mod tidy; go build  -o 'http-server' http-server.go;"
bash -c "$test"
# 构建执行脚本
`echo -e "#!/bin/sh\nexec ./http-server -e $runEnv " > run-http.sh`

# 构建docker image
tag_name=$project-$protocol/$env
docker build -t $tag_name .

# 构建docker container
if [ "$env" == "develop" ]; then

  # 删除旧容器1
  docker stop $project-$protocol-$port1;docker rm $project-$protocol-$port1
  # 新建容器1
  docker run --restart=unless-stopped -itd -v /newdata/logs/golang/:/data/log/ -p $port1:50051 --name $project-$protocol-$port1 $tag_name:latest
  # 删除旧容器2
  docker stop $project-$protocol-$port2;docker rm $project-$protocol-$port2
  # 新建容器2
  docker run --restart=unless-stopped -itd -v /newdata/logs/golang/:/data/log/ -p $port2:50051 --name $project-$protocol-$port2 $tag_name:latest
else
  # 镜像标签镜像
  docker tag $tag_name $repository_url/$tag_name
  # 推送镜像
  docker push $repository_url/$tag_name
fi