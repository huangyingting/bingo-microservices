# https://kubernetes.github.io/ingress-nginx/user-guide/monitoring/

helm upgrade --install ingress-nginx ingress-nginx \
--repo https://kubernetes.github.io/ingress-nginx \
--namespace ingress-nginx --create-namespace \
--set controller.service.externalTrafficPolicy=Local \
--set controller.metrics.enabled=true \
--set controller.metrics.serviceMonitor.enabled=true \
--set-string controller.metrics.serviceMonitor.additionalLabels."release"="prometheus"