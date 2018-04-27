DOCKER_IMAGE = fabianbaier/dcos_metrics_snapshot
DOCKER_TAG := $(shell git rev-parse HEAD)

.PHONY: sadeps saclean sa docker

.DEFAULT_GOAL := sa

docker: docker
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-w' -o dcos_crawler main.go
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
	docker tag $(DOCKER_IMAGE):$(DOCKER_TAG) $(DOCKER_IMAGE):latest
	docker push $(DOCKER_IMAGE):$(DOCKER_TAG)
	docker push $(DOCKER_IMAGE):latest

sadeps:
	dcos package install dcos-enterprise-cli

saclean: sadeps
	dcos security org service-accounts delete dcos_metrics_snapshot || true
	dcos security secrets delete dcos_metrics_snapshot_secret || true

sa: sadeps
	dcos security org service-accounts keypair /tmp/private.pem /tmp/public.pem
	dcos security org service-accounts create -p /tmp/public.pem -d "service account for dcos_metrics_snapshot" dcos_metrics_snapshot
	dcos security secrets create-sa-secret /tmp/private.pem dcos_metrics_snapshot dcos_metrics_snapshot_secret
	curl --fail -k -X PUT \
		--cacert $(shell dcos config show | grep core.ssl_verify | awk '{print $$2}')  \
		-H "Authorization: token=$(shell dcos config show core.dcos_acs_token)" \
		$(shell dcos config show core.dcos_url)/acs/api/v1/acls/dcos:superuser/users/dcos_metrics_snapshot/full