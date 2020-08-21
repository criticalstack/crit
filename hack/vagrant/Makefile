.PHONY: base

build: ../../bin/crit .pki/ca.crt .pki/etcd/ca.crt destroy
	@vagrant up
	$(MAKE) -s import-ctx

../../bin/crit:
	@$(MAKE) -s -C ../../ build

.pki/ca.crt:
	@../../bin/crit pki init --cert-dir .pki/

.pki/etcd/ca.crt:
	@../../bin/crit e2d pki init --ca-cert .pki/etcd/ca.crt --ca-key .pki/etcd/ca.key

destroy:
	@vagrant destroy -f

base:
	@$(MAKE) -C ./base build

CLUSTER=critical-stack
USER=kubernetes-admin
CTX=$(USER)@$(CLUSTER)
KUBECONFIG=$(HOME)/.kube/config

import-ctx:
	@kubectl --kubeconfig=$(KUBECONFIG) config set clusters.$(CLUSTER).server $$(kubectl --kubeconfig=$(PWD)/.kube/config config view -o jsonpath='{.clusters[0].cluster.server}')
	@kubectl --kubeconfig=$(KUBECONFIG) config set clusters.$(CLUSTER).certificate-authority-data $$(kubectl --kubeconfig=$(PWD)/.kube/config config view --raw -o jsonpath='{.clusters[0].cluster.certificate-authority-data}')
	@kubectl --kubeconfig=$(KUBECONFIG) config set users.$(USER).client-certificate-data $$(kubectl --kubeconfig=$(PWD)/.kube/config config view --raw -o jsonpath='{.users[?(@.name == "$(USER)")].user.client-certificate-data}')
	@kubectl --kubeconfig=$(KUBECONFIG) config set users.$(USER).client-key-data $$(kubectl --kubeconfig=$(PWD)/.kube/config config view --raw -o jsonpath='{.users[?(@.name == "$(USER)")].user.client-key-data}')
	@kubectl --kubeconfig=$(KUBECONFIG) config set-context $(CTX) --cluster=$(CLUSTER) --user=$(USER) --namespace=kube-system
	@kubectl --kubeconfig=$(KUBECONFIG) config use-context $(CTX)

remove-ctx:
	-@kubectl --kubeconfig=$(KUBECONFIG) config delete-cluster $(CLUSTER) 2>/dev/null || true
	-@kubectl --kubeconfig=$(KUBECONFIG) config unset users.$(USER) 2>&1 >/dev/null || true
	-@kubectl --kubeconfig=$(KUBECONFIG) config delete-context $(CTX) 2>/dev/null || true
ifeq ($(shell kubectl --kubeconfig=$(KUBECONFIG) config current-context), $(CTX))
	@kubectl --kubeconfig=$(KUBECONFIG) config use-context $$(kubectl --kubeconfig=$(KUBECONFIG) config get-contexts -o name | head -n 1)
endif
