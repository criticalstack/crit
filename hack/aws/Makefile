SHELL := /bin/bash

.PHONY: apply destroy update .kube/config

TOOLS_DIR     := ../tools
TFCONFIG=$(TOOLS_DIR)/tfconfig/bin/tfconfig
CRIT=../../bin/crit

apply: $(CRIT) .pki/ca.crt .pki/etcd/ca.crt config-confirm tfinit tfvalidate tfapply import-ctx cluster-ready
update: tftaint apply
destroy: tfdestroy

$(CRIT):
	@$(MAKE) -s -C ../../ crit

.pki/ca.crt: $(CRIT)
	@[[ -f $@ ]] || $(CRIT) certs init --cert-dir .pki/

.pki/etcd/ca.crt: $(CRIT)
	@[[ -f $@ ]] || $(CRIT) e2d pki init --ca-cert .pki/etcd/ca.crt --ca-key .pki/etcd/ca.key


import-ctx: .kube/config
	@$(CRIT) config import $(PWD)/.kube/config

CONTROL_PLANE_ENDPOINT ?= $$(terraform output -json | jq -r '.control_plane_ips.value[0]')

.kube/config:
	@mkdir -p .kube
	@$(CRIT) generate kubeconfig $(PWD)/.kube/config \
		--cert-dir $(PWD)/.pki \
		--server https://$(CONTROL_PLANE_ENDPOINT):6443

cluster-ready:
	@echo -n "Waiting for controlplane to be ready ..."
	@until kubectl get no >/dev/null 2>&1; do echo -ne "\rWaiting for controlplane to be ready ..."; sleep 2; done
	@echo " Ready!"

$(TFCONFIG):
	@cd $(TOOLS_DIR); $(MAKE) -s -C tfconfig tfconfig

config-confirm: $(TFCONFIG)
	@$(TFCONFIG) create ./

clean:
	@rm -f .pki/ca.crt
	@rm -f .pki/ca.key
	@rm -f .pki/etcd/ca.crt
	@rm -f .pki/etcd/ca.key

tfinit:
	@terraform init
tfvalidate:
	@terraform validate
tfapply:
	@terraform apply -auto-approve
tfdestroy:
	@terraform destroy -auto-approve
tftaint:
	@terraform taint --allow-missing aws_s3_bucket_object.crit_bin_object
	@terraform state list | grep aws_instance | xargs -L 1 terraform taint
