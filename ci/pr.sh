#!/bin/bash

K3S_VERSION=1.0.0

sudo curl -Lo /usr/local/bin/k3s https://github.com/rancher/k3s/releases/download/v${K3S_VERSION}/k3s && sudo chmod a+x /usr/local/bin/k3s
sudo ln -sf k3s /usr/local/bin/kubectl
sudo k3s server --docker --write-kubeconfig-mode 0644 &>/dev/null &

# Build
make test
make build
make image

# Run functional tests
make functest
