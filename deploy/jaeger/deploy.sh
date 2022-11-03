kubectl create ns observability
kubectl create -f https://github.com/jaegertracing/jaeger-operator/releases/download/v1.38.0/jaeger-operator.yaml -n observability
kubectl create ns jaeger
kubectl apply -f jaeger.yaml -n jaeger