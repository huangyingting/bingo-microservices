kubectl create ns rabbitmq
helm install rabbitmq azure-marketplace/rabbitmq -n rabbitmq \
     --set metrics.enabled=true \
     --set-string metrics.serviceMonitor.labels."release"="prometheus" \
     --set metrics.serviceMonitor.enabled=true