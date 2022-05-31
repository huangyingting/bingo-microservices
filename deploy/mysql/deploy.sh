kubectl create namespace mysql
helm install mysql bitnami/mysql \
     --set architecture=replication \
     --set metrics.enabled=true \
     --set-string metrics.serviceMonitor.labels."release"="prometheus" \
     --set metrics.serviceMonitor.enabled=true \
     --set auth.database="bingo" \
     --set auth.username="bingo" \
     --set auth.password="P@ssw0rd" \
     --set auth.replicationPassword="P@ssw0rd" -n mysql