kubectl create ns rabbitmq
helm install rabbitmq bitnami/rabbitmq -n rabbitmq \
     --set metrics.enabled=true \
     --set-string metrics.serviceMonitor.labels."release"="prometheus" \
     --set auth.username=bingo \
     --set auth.password=microservices \
     --set metrics.serviceMonitor.enabled=true