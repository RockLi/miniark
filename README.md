# Miniark

miniark is an (offline) k8s installer

## Prerequisites

1. Docker
2. VirtualBox

You must install docker & Virtualbox first.


## Install

### Download binary 

1. Follow the [URL](https://minio.longguikeji.com/ark/v1.0/miniark-osx) to download the `miniark` binary
2. `chmod +x miniark-osx`
3. move the `miniark` binary to your PATH and rename to `miniark`

### Online Install

`miniark`

This command will download and install the kubernetes environment.


### Offline Install

1. download the package files, put [URL](https://minio.longguikeji.com/ark/v1.0/miniark-offline.tar.gz) here
2. Extract the package file to `$HOME/.miniark`
3. Run `miniark`, everything should works fine now.


## Components

1. Kubernetes
2. Helm