kubectl create ns etcd
helm install etcd bitnami/etcd -n etcd \
     --set auth.rbac.rootPassword=microservices \
     --set metrics.enabled=true

cat <<EOF | kubectl apply -f -
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: etcd
  namespace: etcd
  labels:
    release: prometheus
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: etcd
  endpoints:
  - port: client
    path: /metrics
  namespaceSelector:
    matchNames:
    - etcd
EOF