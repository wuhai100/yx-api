#### 创建pv（本PV来自于华为的云硬盘存储卷）
```bash
kubectl create -f pv.yaml
kubectl get pv
```

#### 创建configmap,用来存放redis的配置文件
```bash
kubectl create configmap redis-conf --from-file=redis.conf
kubectl get cm
```

### 创建headless service
```bash
kubectl create -f headless-svc.yaml
kubectl get sts
```

### 创建用于访问service
```bash
kubectl create -f svc.yaml
kubectl get svc
```

#### 教程来源于https://cloud.tencent.com/developer/article/1392872
