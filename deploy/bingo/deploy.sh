kubectl create ns bingo
kubectl apply -f deploy.yaml -n bingo
kubectl apply -f service.yaml -n bingo
kubectl apply -f ingress.yaml -n bingo
kubectl apply -f metric.yaml -n bingo