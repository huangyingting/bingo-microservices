helm install redis bitnami/redis \
     --set image.registry=ghcr.io \
     --set image.repository=huangyingting/redis \
     --set replica.replicaCount=0 \
     --set image.tag=latest -n redis