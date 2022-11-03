```
NAME: mysql
LAST DEPLOYED: Tue May 31 16:40:30 2022
NAMESPACE: mysql
STATUS: deployed
REVISION: 1
TEST SUITE: None
NOTES:
CHART NAME: mysql
CHART VERSION: 9.0.0
APP VERSION: 8.0.29

** Please be patient while the chart is being deployed **

Tip:

  Watch the deployment status using the command: kubectl get pods -w --namespace mysql

Services:

  echo Primary: mysql-primary.mysql.svc.cluster.local:3306
  echo Secondary: mysql-secondary.mysql.svc.cluster.local:3306

Execute the following to get the administrator credentials:

  echo Username: root
  MYSQL_ROOT_PASSWORD=$(kubectl get secret --namespace mysql mysql -o jsonpath="{.data.mysql-root-password}" | base64 --decode)

To connect to your database:

  1. Run a pod that you can use as a client:

      kubectl run mysql-client --rm --tty -i --restart='Never' --image  marketplace.azurecr.io/bitnami/mysql:8.0.29-debian-10-r15 --namespace mysql --env MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD --command -- bash

  2. To connect to primary service (read/write):

      mysql -h mysql-primary.mysql.svc.cluster.local -uroot -p"$MYSQL_ROOT_PASSWORD"

  3. To connect to secondary service (read-only):

      mysql -h mysql-secondary.mysql.svc.cluster.local -uroot -p"$MYSQL_ROOT_PASSWORD"



To access the MySQL Prometheus metrics from outside the cluster execute the following commands:

    kubectl port-forward --namespace mysql svc/mysql-metrics 9104:9104 &
    curl http://127.0.0.1:9104/metrics
```