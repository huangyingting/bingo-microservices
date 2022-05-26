kubectl create namespace mysql
helm install mysql azure-marketplace/mysql \
     --set metrics.enabled=true \
     --set-string metrics.serviceMonitor.labels."release"="prometheus" \
     --set metrics.serviceMonitor.enabled=true \
     --set auth.username="bingo" \
     --set auth.password="P@ssw0rd" -n mysql