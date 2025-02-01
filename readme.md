# Hello-k8s

Hello-k8s is a simple project to demonstrate the setup of a self-hosted Kubernetes cluster and the deployment of a
stateless Go application to the cluster.

## Prerequisites

To understand this project and especially the Kubernetes setup, you need to have a basic understanding of the following:

- Containerization and more specifically Docker
- Networking
- Unix systems
- Any programming language (Go in this case)

## Projects Structures

Inside every example, you will find the following structure:

- `/<exampleApp>/manifests/`: The Kubernetes manifests for the applications including deployments, services, and ingresses configurations
- `/<exampleApp>/app/.../`: The source code for the applications to deploy
- `/<exampleApp>/app/.../Dockerfile`: Inside every app directory, you will find a Dockerfile to build the application image

## Table of Contents

- [Initial Kubernetes setup](#initial-kubernetes-setup)
- [Deploying the stateless app to Kubernetes](#deploying-the-stateless-app-to-kubernetes)
- Deploying a stateful app to Kubernetes (Coming soon)

## Initial Kubernetes setup

The k8s cluster is set up using kubeadm and is, at first, composed of :

- 1 control-plane node
- 2 worker nodes

The plugins used are:

- Flannel for networking
- Nginx for ingress

The cluster was initialized using `kubeadm` with a configuration file to specify the cgroup driver and the container
runtime.

Here's an overview of the configuration file:

```yaml
apiVersion: kubeadm.k8s.io/v1beta4
kind: InitConfiguration
nodeRegistration:
  criSocket: unix:///var/run/cri-dockerd.sock # The container runtime socket
  taints: null
---
apiVersion: kubeadm.k8s.io/v1beta4
imageRepository: registry.k8s.io
kind: ClusterConfiguration
kubernetesVersion: 1.32.0
proxy: { }
---
kind: KubeletConfiguration
apiVersion: kubelet.config.k8s.io/v1beta1
cgroupDriver: systemd # The cgroup driver
```

Before initializing the cluster, we need to disable swap on all nodes and enable the br_netfilter kernel module:

```bash
swapoff -a
modprobe br_netfilter
```

Then using `kubeadm init` with the config file we can initialize the control-plane node:

```bash
kubeadm init -f kubeadm-config.yaml
```

After the control-plane node is initialized, a join command is generated. Run we can add the worker nodes to the cluster
by running the following command on each worker node:
```bash
kubeadm join <control-plane-host>:<control-plane-port> --token <token>
```

After the worker nodes are added, we can install the Flannel networking plugin:

```bash
kubectl apply -f https://github.com/flannel-io/flannel/releases/latest/download/kube-flannel.yml
```

One the networking plugin is installed, we can install the Nginx ingress controller :
```bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.12.0/deploy/static/provider/baremetal/deploy.yaml
```

After the ingress controller is installed, we can run the following command to get the nodeport of the ingress service:

```bash
kubectl get services -n nginx-ingress -o wide

# Output
NAME                                 TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)                      AGE
ingress-nginx-controller             NodePort    10.107.116.103   <none>        80:30736/TCP,443:32186/TCP   12m
ingress-nginx-controller-admission   ClusterIP   10.110.149.189   <none>        443/TCP                      12m
```

The ingress service is exposed on port 30736 for HTTP and 32186 for HTTPS and can be accessed from any of the nodes in the cluster.

The cluster is now set up and ready to deploy applications.

## Deploying the stateless app to Kubernetes
### Environment Variables

The application uses the following environment variables:

- `nodename`: The name of the node where the pod is running, retrieved from the Kubernetes field reference.

### Resources

The deployment requests the following resources:

- Memory: 400Mi (400 mebibytes)
- CPU: 250m (250 millicores, equivalent to 0.25 cores)

### Affinity Rules

The deployment uses pod anti-affinity rules to ensure that pods are not scheduled on the same node.

To deploy the application to Kubernetes, run the following commands:

```bash
kubectl apply -f ./manifests
```

This will create the deployment, service, and ingress resources in the Kubernetes cluster.