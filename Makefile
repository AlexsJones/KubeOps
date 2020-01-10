VERSION=`cat VERSION`
.PHONY: publish docker install up down


all: docker publish install

up:
	kind create cluster --config kind/config.yaml

down:
	kind delete cluster

publish:
	kind load docker-image kubeops:$(VERSION)

docker:
	docker build -t kubeops:$(VERSION) .

install:
	cd helm && helm install . --generate-name && cd ../

delete:
	helm ls --all --short | xargs -L1 helm delete