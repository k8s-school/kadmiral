# Pre-requisites

A few Ubuntu instances with `ss` and `sudo` access

# Configure access to the instances

Fill your `~/.ssh/config` and edit `config.yaml` with the master, nodes, and
user information:

```yaml
distrib: ubuntu
master: k8s-master-1
nodes:
  - k8s-node1
  - harbor-node2
user: root
scp: scp
ssh: ssh
```

When `kadmiral rsync` runs, it generates an `env.sh` file from this configuration
on each node before executing the scripts.

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

# Go based CLI

The `kadmiral` command wraps these scripts and runs them in parallel over SSH.
Example usage:

```shell
kadmiral --nodes master1,node1,node2 --user root rsync
kadmiral prereq --os ubuntu
kadmiral init --node master1
kadmiral cni install cilium
kadmiral join node node2 --master master1:6443 --token <token> --ca-hash <hash>
kadmiral reset node2
kadmiral reset --all
```

The CLI uses SSH to copy the repository to each node and execute the relevant
scripts concurrently. Logging is provided via Go's `slog` with a `--log-level`
flag for controlling verbosity.
