#!/usr/bin/env bash
re="swr.cn-east-2.myhuaweicloud.com/yb7/([^:]+):([^ ]+)"
imageStr=$(kubectl --kubeconfig $KUBECONFIG_TEST get deploy yx-api -o jsonpath='{..image}')
echo "current version"
if [[ $imageStr =~ $re ]]; then echo ${BASH_REMATCH[2]}; fi

if [ $# -eq 0 ];
then nextVersion=$(./increment_version.sh -p ${BASH_REMATCH[2]});
else nextVersion=$(./increment_version.sh $1 ${BASH_REMATCH[2]});
fi

echo "next version"
echo $nextVersion;

rm -rf yx-api
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o yx-api ../main.go
docker build -t swr.cn-east-2.myhuaweicloud.com/yb7/yx-api:$nextVersion .
docker push swr.cn-east-2.myhuaweicloud.com/yb7/yx-api:$nextVersion
rm -rf yx-api
kubectl --kubeconfig $KUBECONFIG_TEST set image deployment/yx-api yx-api=swr.cn-east-2.myhuaweicloud.com/yb7/yx-api:$nextVersion
sleep 2
kubectl --kubeconfig $KUBECONFIG_TEST get pod -l app=yx-api

git add ../
git commit -m "$nextVersion $1"
git push
echo $(date "+%Y-%m-%d %H:%M:%S")
