BASE_BOX = crit-base.box

build: clean build-tools
	@vagrant up
	@vagrant package --output=$(BASE_BOX)
	@vagrant destroy -f
	@vagrant box add -f --name crit $(BASE_BOX)
	@$(MAKE) -s clean

build-tools:
	@$(MAKE) -s -C ../../tools/containerd-sync build

clean:
	@rm -f $(BASE_BOX)
