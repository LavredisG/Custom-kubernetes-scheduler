# Custom scheduler for a Kubernetes cluster

A custom Kubernetes scheduler written in Go. In its first iteration, it assigns pods
based on a key-value pair we define. (WIP)

## Prerequisites

- Docker
- kubectl
- Kubernetes cluster (using `kind`)

## Setting Up a `kind` Cluster

1. Create a `kind` cluster with 3 nodes (1 control plane and 2 worker nodes):
   ```sh
   cat <<EOF | kind create cluster --config=-
   kind: Cluster
   apiVersion: kind.x-k8s.io/v1alpha4
   nodes:
     - role: control-plane
     - role: worker
     - role: worker
   EOF
   ```

2. Label one of the worker nodes with the key/value pair for scheduling
```sh
kubectl label node <worker-node-name> <key>=<value>
```
For our example we have used "example-label-key=example-label-value" but it's easy configurable in the <b>main.go</b>.

## Deploying the scheduler

```sh
kubectl apply -f rbac.yaml
kubectl apply -f deployment.yaml
```

## Testing

```sh
kubectl apply -f pods_deployment.yaml
```
