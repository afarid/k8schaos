cnf ?= app.env
include $(cnf)
export $(shell sed 's/=.*//' $(cnf))

dpl ?= deploy.env
include $(dpl)
export $(shell sed 's/=.*//' $(dpl))

build:
	docker build -t k8schaos:$(VERSION) .

load-image:
	kind load docker-image k8schaos:$(VERSION) k8schaos:$(VERSION)

deploy: build load-image
	helm upgrade --install  --namespace $(DEPLOYMENT_NAMESPACE) --set image.tag=$(VERSION) --set config.namespace=$(NAMESPACE) \
    --set config.timePeriod=$(TIME_PERIOD)  k8schaos ./deploy/k8schaos --kube-context $(K8S_CONTEXT)

remove:
	helm uninstall k8schaos --namespace $(DEPLOYMENT_NAMESPACE) --kube-context $(K8S_CONTEXT)

.PHONY: build load-image deploy