VERSION=`cat VERSION`
.PHONY: publish docker install up down


all: docker publish install

up:
	kind create cluster --name=kind

down:
	kind delete cluster --name=kind

publish:
	export KUBECONFIG="$(kind get kubeconfig-path --name="kind")"
	kind load docker-image --name=kind kubeops

docker:
	docker build -t kubeops .

install:
	cd helm && helm install . --generate-name && cd ../

delete:
	helm ls --all --short | xargs -L1 helm delete
