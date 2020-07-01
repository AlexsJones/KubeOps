VERSION=`cat VERSION`
.PHONY: publish docker install up down

rename: 
	gomove github.com/AlexsJones/kubeops $(NAME)

all: docker publish install

run-builtin-example:
	go run examples/builtin/main.go --kubeconfig=$(HOME)/.kube/config

run-crd-example:
	go run examples/crd/main.go --kubeconfig=$(HOME)/.kube/config

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

list:
	@$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$'

