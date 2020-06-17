VERSION=`cat VERSION`
.PHONY: publish docker install up down


all: docker publish install

up:
	kind create cluster --name=kind

down:
	kind delete cluster --name=kind

publish:
	kind load docker-image kubeops:$(VERSION) --name=kind

docker:
	go mod vendor
	docker build -t kubeops:$(VERSION) .

install:
	cd helm && helm install . --generate-name && cd ../

delete:
	helm ls --all --short | xargs -L1 helm delete
