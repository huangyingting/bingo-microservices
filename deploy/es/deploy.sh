kubectl create -f https://download.elastic.co/downloads/eck/2.2.0/crds.yaml
kubectl create ns es
kubectl apply -f es.yaml -n es