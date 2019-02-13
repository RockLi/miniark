#!/bin/bash

images=(
    'gcr.io/k8s-minikube/storage-provisioner:v1.8.1'
    'k8s.gcr.io/coredns:1.2.2'
    'k8s.gcr.io/etcd:3.2.24'
    'k8s.gcr.io/kube-addon-manager:v8.6'
    'k8s.gcr.io/kube-apiserver:v1.12.4'
    'k8s.gcr.io/kube-controller-manager:v1.12.4'
    'k8s.gcr.io/kube-proxy:v1.12.4'
    'k8s.gcr.io/kube-scheduler:v1.12.4'
    'k8s.gcr.io/kubernetes-dashboard-amd64_v1.10.1'
    'k8s.gcr.io/pause:3.1'
)

for image in "${images[@]}"
do
    docker pull ${image}
done