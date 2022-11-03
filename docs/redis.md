```
helm install redis bitnami/redis \
     --set replica.replicaCount=3 \
     --set sentinel.enabled=true \
     --set sentinel.masterSet=bingo \
     --set sentinel.quorum=3 \
     --set metrics.enabled=true \
     --set-string metrics.serviceMonitor.additionalLabels."release"="prometheus" \
     --set metrics.serviceMonitor.enabled=true \
     --set image.tag=latest -n redis
NAME: redis
LAST DEPLOYED: Mon May 30 09:21:12 2022
NAMESPACE: redis
STATUS: deployed
REVISION: 1
TEST SUITE: None
NOTES:
CHART NAME: redis
CHART VERSION: 16.9.5
APP VERSION: 6.2.7

** Please be patient while the chart is being deployed **

Redis&trade; can be accessed via port 6379 on the following DNS name from within your cluster:

    redis.redis.svc.cluster.local for read only operations

For read/write operations, first access the Redis&trade; Sentinel cluster, which is available in port 26379 using the same domain name above.



To get your password run:

    export REDIS_PASSWORD=$(kubectl get secret --namespace redis redis -o jsonpath="{.data.redis-password}" | base64 --decode)

To connect to your Redis&trade; server:

1. Run a Redis&trade; pod that you can use as a client:

   kubectl run --namespace redis redis-client --restart='Never'  --env REDIS_PASSWORD=$REDIS_PASSWORD  --image docker.io/bitnami/redis:latest --command -- sleep infinity

   Use the following command to attach to the pod:

   kubectl exec --tty -i redis-client \
   --namespace redis -- bash

2. Connect using the Redis&trade; CLI:
   REDISCLI_AUTH="$REDIS_PASSWORD" redis-cli -h redis -p 6379 # Read only operations
   REDISCLI_AUTH="$REDIS_PASSWORD" redis-cli -h redis -p 26379 # Sentinel access

To connect to your database from outside the cluster execute the following commands:

    kubectl port-forward --namespace redis svc/redis 6379:6379 &
    REDISCLI_AUTH="$REDIS_PASSWORD" redis-cli -h 127.0.0.1 -p 6379
```

```
k logs redis-node-0 -n redis 
Defaulted container "redis" out of: redis, sentinel, metrics
 01:21:44.93 INFO  ==> about to run the command: REDISCLI_AUTH=$REDIS_PASSWORD timeout 220 redis-cli -h redis.redis.svc.cluster.local -p 26379 sentinel get-master-addr-by-name bingo
Could not connect to Redis at redis.redis.svc.cluster.local:26379: Connection refused
 01:21:50.98 INFO  ==> Configuring the node as master
1:C 30 May 2022 01:21:51.005 # oO0OoO0OoO0Oo Redis is starting oO0OoO0OoO0Oo
1:C 30 May 2022 01:21:51.005 # Redis version=6.2.7, bits=64, commit=00000000, modified=0, pid=1, just started
1:C 30 May 2022 01:21:51.005 # Configuration loaded
1:M 30 May 2022 01:21:51.005 * monotonic clock: POSIX clock_gettime
1:M 30 May 2022 01:21:51.011 # A key '__redis__compare_helper' was added to Lua globals which is not on the globals allow list nor listed on the deny list.
1:M 30 May 2022 01:21:51.011 * Running mode=standalone, port=6379.
1:M 30 May 2022 01:21:51.011 # Server initialized
1:M 30 May 2022 01:21:51.011 * Ready to accept connections
1:M 30 May 2022 01:22:38.218 * Replica redis-node-1.redis-headless.redis.svc.cluster.local:6379 asks for synchronization
1:M 30 May 2022 01:22:38.218 * Full resync requested by replica redis-node-1.redis-headless.redis.svc.cluster.local:6379
1:M 30 May 2022 01:22:38.218 * Replication backlog created, my new replication IDs are 'cd2340f1336c3fa4437a00d1cac4955f5010b4f7' and '0000000000000000000000000000000000000000'
1:M 30 May 2022 01:22:38.218 * Starting BGSAVE for SYNC with target: disk
1:M 30 May 2022 01:22:38.219 * Background saving started by pid 198
198:C 30 May 2022 01:22:38.227 * DB saved on disk
198:C 30 May 2022 01:22:38.227 * RDB: 2 MB of memory used by copy-on-write
1:M 30 May 2022 01:22:38.288 * Background saving terminated with success
1:M 30 May 2022 01:22:38.288 * Synchronization with replica redis-node-1.redis-headless.redis.svc.cluster.local:6379 succeeded
1:M 30 May 2022 01:23:27.818 * Replica redis-node-2.redis-headless.redis.svc.cluster.local:6379 asks for synchronization
1:M 30 May 2022 01:23:27.818 * Full resync requested by replica redis-node-2.redis-headless.redis.svc.cluster.local:6379
1:M 30 May 2022 01:23:27.818 * Starting BGSAVE for SYNC with target: disk
1:M 30 May 2022 01:23:27.818 * Background saving started by pid 442
442:C 30 May 2022 01:23:27.827 * DB saved on disk
442:C 30 May 2022 01:23:27.827 * RDB: 2 MB of memory used by copy-on-write
1:M 30 May 2022 01:23:27.923 * Background saving terminated with success
1:M 30 May 2022 01:23:27.923 * Synchronization with replica redis-node-2.redis-headless.redis.svc.cluster.local:6379 succeeded
```

```
k logs redis-node-2 -n redis 
Defaulted container "redis" out of: redis, sentinel, metrics
 01:23:22.72 INFO  ==> about to run the command: REDISCLI_AUTH=$REDIS_PASSWORD timeout 220 redis-cli -h redis.redis.svc.cluster.local -p 26379 sentinel get-master-addr-by-name bingo
 01:23:27.76 INFO  ==> about to run the command: REDISCLI_AUTH=$REDIS_PASSWORD timeout 220 redis-cli -h redis.redis.svc.cluster.local -p 26379 sentinel get-master-addr-by-name bingo
 01:23:27.79 INFO  ==> Current master: REDIS_SENTINEL_INFO=(redis-node-0.redis-headless.redis.svc.cluster.local,6379)
 01:23:27.79 INFO  ==> Configuring the node as replica
1:C 30 May 2022 01:23:27.806 # oO0OoO0OoO0Oo Redis is starting oO0OoO0OoO0Oo
1:C 30 May 2022 01:23:27.806 # Redis version=6.2.7, bits=64, commit=00000000, modified=0, pid=1, just started
1:C 30 May 2022 01:23:27.806 # Configuration loaded
1:S 30 May 2022 01:23:27.807 * monotonic clock: POSIX clock_gettime
1:S 30 May 2022 01:23:27.807 # A key '__redis__compare_helper' was added to Lua globals which is not on the globals allow list nor listed on the deny list.
1:S 30 May 2022 01:23:27.807 * Running mode=standalone, port=6379.
1:S 30 May 2022 01:23:27.808 # Server initialized
1:S 30 May 2022 01:23:27.808 * Ready to accept connections
1:S 30 May 2022 01:23:27.809 * Connecting to MASTER redis-node-0.redis-headless.redis.svc.cluster.local:6379
1:S 30 May 2022 01:23:27.812 * MASTER <-> REPLICA sync started
1:S 30 May 2022 01:23:27.814 * Non blocking connect for SYNC fired the event.
1:S 30 May 2022 01:23:27.814 * Master replied to PING, replication can continue...
1:S 30 May 2022 01:23:27.817 * Partial resynchronization not possible (no cached master)
1:S 30 May 2022 01:23:27.818 * Full resync from master: cd2340f1336c3fa4437a00d1cac4955f5010b4f7:10628
1:S 30 May 2022 01:23:27.924 * MASTER <-> REPLICA sync: receiving 176 bytes from master to disk
1:S 30 May 2022 01:23:27.924 * MASTER <-> REPLICA sync: Flushing old data
1:S 30 May 2022 01:23:27.931 * MASTER <-> REPLICA sync: Loading DB in memory
1:S 30 May 2022 01:23:27.939 * Loading RDB produced by version 6.2.7
1:S 30 May 2022 01:23:27.940 * RDB age 0 seconds
1:S 30 May 2022 01:23:27.940 * RDB memory usage when created 1.93 Mb
1:S 30 May 2022 01:23:27.940 # Done loading RDB, keys loaded: 0, keys expired: 0.
1:S 30 May 2022 01:23:27.940 * MASTER <-> REPLICA sync: Finished with success
1:S 30 May 2022 01:23:27.940 * Background append only file rewriting started by pid 35
1:S 30 May 2022 01:23:27.971 * AOF rewrite child asks to stop sending diffs.
35:C 30 May 2022 01:23:27.972 * Parent agreed to stop sending diffs. Finalizing AOF...
35:C 30 May 2022 01:23:27.972 * Concatenating 0.00 MB of AOF diff received from parent.
35:C 30 May 2022 01:23:27.972 * SYNC append only file rewrite performed
35:C 30 May 2022 01:23:27.972 * AOF rewrite: 0 MB of memory used by copy-on-write
1:S 30 May 2022 01:23:28.014 * Background AOF rewrite terminated with success
1:S 30 May 2022 01:23:28.014 * Residual parent diff successfully flushed to the rewritten AOF (0.00 MB)
1:S 30 May 2022 01:23:28.014 * Background AOF rewrite finished successfully
```