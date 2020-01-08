#!/bin/bash
kubectl --kubeconfig $KUBECONFIG_PROD get deployment/yx-api -o wide -n default
kubectl --kubeconfig $KUBECONFIG_PROD -n default set image deployment/yx-api yx-api=swr.cn-east-2.myhuaweicloud.com/yb7/yx-api:$1
kubectl --kubeconfig $KUBECONFIG_PROD get deployment/yx-api -o wide -n default
