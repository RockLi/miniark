# Miniark

Miniark is an (offline) k8s installer supposed to quickly setup local dev k8s environment.

For the full featured cluster version, we will opensource in the future.

## Background

因为科学上网的问题，造成初学者在学习K8S的时候需要具备一定的技能才能搭建好一个k8s环境，于是萌生了本项目。
同时miniark也支持离线安装，以便加快mini k8s环境的安装速度。

## Version

1. minikube: v0.32.0
2. kubernetes: v1.12.4
3. helm: v2.12.3

Normally we will use the last stable version of k8s.

## Platform

1. MacOSX
2. Linux(TBD)

## Components Included

1. Kubernetes(minikube)
2. Helm

## Prerequisites

1. Docker
2. VirtualBox


## Install

### Download binary 

1. Follow the [URL](https://minio.longguikeji.com/ark/v1.0/miniark-osx) to download the `miniark` binary
2. `chmod +x miniark-osx`
3. move the `miniark` binary to your PATH and rename to `miniark`

We support both online and offline installation.

### Online Installation

`miniark`

This command will download and install the kubernetes environment automatically.


### Offline Installation

1. Download the pre-packed files, click link [URL](https://minio.longguikeji.com/ark/v1.0/miniark-offline.tar.gz)
2. You can execute command `tar zxvf miniark-offline.tar.gz -C ~/` to extract the files to `$HOME/.miniark`
3. Run `miniark`, everything should works fine now.


## Warning

Due to the feature of offline installation, we make a modified minikube, you can checkout the [commit](https://github.com/RockLi/minikube/commit/06856df3a8a8af8a0893abc5fb9375bd770bfd74). 

