# Hello-k8s

Hello-k8s is a simple project to demonstrate the setup of a self-hosted Kubernetes cluster and the deployment of a
stateless Go application to the cluster.

## Prerequisites

- Docker
- Kubernetes
- kubectl
- Go

## Kubernetes setup

The k8s cluster is set up using kubeadm and is composed of :
- 1 control-plane node
- 2 worker nodes

The plugins used are:
- Flannel for networking
- Nginx for ingress

## Project Structure

- `/manifests`: The Kubernetes manifests for the application including deployment, service, and ingress configurations
- `/app`: The source code for the Go application
- `/app/Dockerfile`: The Dockerfile for the Go application

## Deploying to Kubernetes

To deploy the application to Kubernetes, run the following commands:

```bash
kubectl apply -f ./manifests
```

This will create the deployment, service, and ingress resources in the Kubernetes cluster.

## Environment Variables

The application uses the following environment variables:

- `nodename`: The name of the node where the pod is running, retrieved from the Kubernetes field reference.

## Resources

The deployment requests the following resources:

- Memory: 400Mi
- CPU: 250m

## Affinity Rules

The deployment uses pod anti-affinity rules to ensure that pods are not scheduled on the same node.