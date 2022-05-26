kubectl create ns redis
helm install redis bitnami/redis \
     --set image.registry=ghcr.io \
     --set image.repository=huangyingting/redis \
     --set replica.replicaCount=0 \
     --set metrics.enabled=true \
     --set-string metrics.serviceMonitor.additionalLabels."release"="prometheus" \
     --set metrics.serviceMonitor.enabled=true \
     --set image.tag=latest -n redis