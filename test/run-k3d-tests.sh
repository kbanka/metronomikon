#!/bin/bash
set -e
REPO_DIR=$( cd "$(dirname $( dirname "${BASH_SOURCE[0]}"))" >/dev/null 2>&1 && pwd )
SERVICE_NAME="metronomikon"
GOLANG_IMAGE="golang:1.13"
KUBECTL_IMAGE=${TERRAFORM_IMAGE:-948223650506.dkr.ecr.us-east-1.amazonaws.com/applause/terraform:0.12.19.1}
KUBECTL_BIN="/usr/local/bin/kubectl"

check_k3d() {
  k3d version || { echo 'k3d does not exist'; exit 1; }
}

delete_k3d_if_exists() {
  k3d d || true
}

deploy_k3d() {
  k3d shell -c 'true' || (delete_k3d_if_exists ; k3d create --publish 80:80)
}

deploy_dashboard_on_k3d() {
  make image
  # wait for k3s to start
  while ! k3d get-kubeconfig --name='k3s-default' 2>&1 > /dev/null; do sleep 1; done
  export KUBECONFIG="$(k3d get-kubeconfig --name='k3s-default')"
  export KUBECTL_CMD="docker run --rm -v ${REPO_DIR}:/code -v $(dirname $KUBECONFIG):/root/.kube -e KUBECONFIG=/root/.kube/kubeconfig.yaml -w /code --entrypoint=${KUBECTL_BIN} --network=host ${KUBECTL_IMAGE}"
  k3d import-images applause/${SERVICE_NAME}:latest
}

deploy_test_resources() {
  for manifest in ${REPO_DIR}/test/data/*.yaml; do
    ${KUBECTL_CMD} delete -f test/data/$(basename $manifest) || true
  done
  for manifest in ${REPO_DIR}/example/k8s/*.yaml; do
    ${KUBECTL_CMD} delete -f example/k8s/$(basename $manifest) || true
  done
  for manifest in ${REPO_DIR}/test/data/*.yaml; do
    ${KUBECTL_CMD} apply -f test/data/$(basename $manifest)
  done
  for manifest in ${REPO_DIR}/example/k8s/*.yaml; do
    ${KUBECTL_CMD} apply -f example/k8s/$(basename $manifest)
  done
}

wait_for_dashboard_to_be_ready() {
  echo "Waiting for traefik to be ready"
  while ! curl -s 127.0.0.1/ping > /dev/null; do sleep 1; done
  echo "Waiting for api to be ready"
  while ! curl -s -f 127.0.0.1/v1/jobs > /dev/null; do sleep 0.2; done
}

run_tests() {
  get_endpoint
  docker build -t ${SERVICE_NAME}-test ${REPO_DIR}/test
  for i in $(ls -1 ${REPO_DIR}/test/tests | sort); do
    docker run --rm --env=PYTHONUNBUFFERED=TRUE --network=host ${SERVICE_NAME}-test ${ENDPOINT} tests/${i} || return 1
  done
}

get_endpoint() {
  ENDPOINT="http://127.0.0.1:80"
}

main() {
  check_k3d
  deploy_k3d
  deploy_dashboard_on_k3d
  deploy_test_resources
  wait_for_dashboard_to_be_ready
  run_tests
  delete_k3d_if_exists
}

if [ "${1}" != "--source-only" ]; then
    main "${@}"
fi
