# Miniark

Miniark is an (offline) k8s MacOSX installer to quickly setup local dev k8s environment.

For full featured cluster version, we will opensource it in the future.

## Background

因为科学上网的问题，造成初学者在学习K8S的时候需要具备一定的技能才能搭建好一个k8s环境，于是萌生了本项目。
同时miniark也支持离线安装，以便加快mini k8s环境的安装速度。

## Version

1. minikube: v0.32.0
2. kubernetes: v1.12.4
3. helm: v2.12.3


## Platform

1. MacOSX
2. Linux(TBD)

## Components Included

1. Kubernetes(minikube)
2. Helm

## Prerequisites

Please manually install these dependencies before run miniark.

1. [Docker](https://www.docker.com/products/docker-desktop)
2. [VirtualBox](https://www.virtualbox.org/wiki/Downloads)


## Install

### Download binary 

1. Follow the [URL](https://github.com/rockl2e/miniark/releases/download/v1.0.0/miniark) to download the `miniark` binary
2. `chmod +x miniark`
3. move the `miniark` binary to your PATH and rename to `miniark`

We support both online and offline installation.

### Online Installation

Run `miniark` directly, it will download and install the kubernetes environment automatically.


### Offline Installation(Recommended)

1. Download the pre-packed files, click link here [URL](https://minio.longguikeji.com/ark/v1.0/miniark-offline.tar.gz)
2. Execute command `tar zxvf miniark-offline.tar.gz -C ~/` to extract files to `$HOME/.miniark`
3. Run `miniark`, everything should works fine now in a minute.


## Warning

Since we need the offline installation feature, we made a modified minikube, you can checkout the [commit](https://github.com/RockLi/minikube/commit/06856df3a8a8af8a0893abc5fb9375bd770bfd74) in case you have security concerns.

