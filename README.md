# Pre-requisites

A few Ubuntu instances with `ssh` and `sudo` access

# Configure access to the instances

Fill your  `~/.ssh/config` and then edit `MASTER`, `NODES`and `USER` in env.sh

```shell
MASTER="k8s-master-1"
NODES="k8s-node1 harbor-node2"

```

# Manage the cluster

```shell
# Create the cluster
./create.sh

# Reset the cluster
./reset.sh

# Connect to k8s
ssh k8s-master-1
kubectl get pods -A
```
