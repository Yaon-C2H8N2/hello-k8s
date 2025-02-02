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

- `/<exampleApp>/manifests/`: The Kubernetes manifests for the applications including deployments, services, and
  ingresses configurations
- `/<exampleApp>/app/.../`: The source code for the applications to deploy
- `/<exampleApp>/app/.../Dockerfile`: Inside every app directory, you will find a Dockerfile to build the application
  image

## Table of Contents

- [Initial Kubernetes setup](#initial-kubernetes-setup)
- [Deploying the stateless app to Kubernetes](#deploying-the-stateless-app-to-kubernetes)
- [Deploying a stateful app to Kubernetes](#deploying-the-stateful-app-to-kubernetes)

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

The ingress service is exposed on port 30736 for HTTP and 32186 for HTTPS and can be accessed from any of the nodes in
the cluster.

The cluster is now set up and ready to deploy applications.

## Deploying the stateless app to Kubernetes

### Application Overview

The stateless application is a simple Go application that returns an HTML response with the name of the node where the
response was generated.

### Manifests

The deployment to the k8s cluster is done using 3 manifests files :

- `statelessApp/manifests/deployment.yaml`: The deployment manifest for the application. It contains the image to use,
  the environment variables, the resources requests, and the affinity rules. In this case the deployment requests 400Mi
  of memory and 250m of CPU and uses pod anti-affinity rules to ensure that pods are not scheduled on the same node.
- `statelessApp/manifests/service.yaml`: The service manifest for the application. It exposes the application on port 80
  and uses a ClusterIP service type.
- `statelessApp/manifests/ingress.yaml`: The ingress manifest for the application. Once the service is created, the
  ingress resource is created to expose the service to the outside world. The ingress uses the Nginx ingress controller
  to route traffic to the service.

### Environment Variables

The application uses the following environment variables:

- `nodename`: The name of the node where the pod is running, retrieved from the Kubernetes field reference.

### Deploying the application

```bash
kubectl apply -f ./manifests
```

This will create the deployment, service, and ingress resources in the Kubernetes cluster.

The app should now be accessible using the url set in the ingress manifest file (here the hostname of a node inside the
cluster) and the port of the ingress service (30736 in this case).

## Deploying the stateful app to Kubernetes

### Application Overview

The stateful application is a simple todo list app. It's composed of a React frontend, a Go API and a Postgres database.
It will be deployed in 2 parts:

- The frontend and the API will be deployed as a stateless app.
- The Postgres database will be deployed as a stateful app.

### Manifest

The deployment of the stateful app is a bit more complex than the stateless app. It's composed of 6 manifests files:

- `statefulApp/manifests/app-deployment.yaml` : The deployment manifest for the frontend and the API. Pretty similar to
  the stateless app deployment manifest.
- `statefulApp/manifests/service.yaml` : The service manifest for the frontend, the API and the Postgres database. The
  frontend and the API are exposed on port 81 and 8081. The Postgres database is exposed on port 5432. The service type
  of the front and API service is ClusterIP. In this case, the Postgres service should be of type ClusterIP to allow the
  API to connect to the database. However for testing purposes, the service is of type NodePort to allow external access
  to inspect the database.
- `statefulApp/manifests/ingress.yaml` : The ingress manifest for the frontend and the API. A URL rewrite is used to
  route traffic to the API when the path starts with `/api`. The ingress controller acts as a reverse proxy to route the
  traffic to the correct service.
- `statefulApp/manifests/volume-claim.yaml` : The volume claim manifest to request a persistent volume to store the
  Postgres data.
- `statefulApp/manifests/persistent-volume.yaml` : Once the volume claim is created, the persistent volume is created to
  bind the volume to the host path.
- `statefulApp/manifests/pgsql-deployment.yaml` : The deployment manifest for the Postgres database. It contains the
  image to use, the environment variables and the volume mounts to persist the data as well as the volume claim to
  request.

### Environment Variables

In this example there is a bit more environment variables to set:

| Name              | Value                                                            | Used by       |
|-------------------|------------------------------------------------------------------|---------------|
| nodename          | spec.nodename (Injected by kubernetes)                           | API           |
| POSTGRES_USER     | postgres                                                         | API, Postgres |
| POSTGRES_PASSWORD | postgres                                                         | API, Postgres |
| POSTGRES_URL      | hello-k8s-stateful-postgres (Name of the pgsql database service) | API           |
| POSTGRES_PORT     | 5432                                                             | API           |
| POSTGRES_DB       | postgres                                                         | API, Postgres |
| JWT_SECRET_KEY    | abcdef1234567890                                                 | API           |

### Deploying the application

```bash
kubectl apply -f ./manifests
```

This will create the deployment, service, ingress, volume claim, persistent volume, and pgsql deployment resources in
the Kubernetes cluster. Exactly like the stateless app, the app should now be accessible using the url set in the
ingress manifest file (here the hostname of a node inside the cluster) and the port of the ingress service (30736 in
this case).