kubectl create ns redis

# sentinel
helm install redis bitnami/redis \
     --set architecture=replication \
     --set image.registry=ghcr.io \
     --set image.repository=huangyingting/redis \
     --set replica.replicaCount=3 \
     --set sentinel.enabled=true \
     --set sentinel.masterSet=bingo \
     --set sentinel.quorum=3 \
     --set metrics.enabled=true \
     --set-string metrics.serviceMonitor.additionalLabels."release"="prometheus" \
     --set metrics.serviceMonitor.enabled=true \
     --set image.tag=latest -n redis

# standalone
helm install redis bitnami/redis \
     --set image.registry=ghcr.io \
     --set image.repository=huangyingting/redis \
     --set architecture=standalone \
     --set metrics.enabled=true \
     --set-string metrics.serviceMonitor.additionalLabels."release"="prometheus" \
     --set metrics.serviceMonitor.enabled=true \
     --set auth.password=microservices \
     --set image.tag=latest -n redis
