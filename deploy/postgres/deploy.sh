kubectl create ns postgres
helm install postgres bitnami/postgresql \
     --set architecture=replication \
     --set metrics.enabled=true \
     --set-string metrics.serviceMonitor.labels."release"="prometheus" \
     --set metrics.serviceMonitor.enabled=true \
     --set auth.database="bingo" \
     --set auth.username="bingo" \
     --set auth.password="microservices" \
     --set auth.replicationPassword="microservices" -n postgres