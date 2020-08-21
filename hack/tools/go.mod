module github.com/criticalstack/crit/hack/tools

go 1.13

require (
	github.com/criticalstack/crit v0.0.0-00010101000000-000000000000
	github.com/golangci/golangci-lint v1.27.0
	github.com/shurcooL/httpfs v0.0.0-20190707220628-8d4bc4ba7749
	github.com/shurcooL/vfsgen v0.0.0-20181202132449-6a9ea43bcacd
	github.com/spf13/cobra v1.0.0
	k8s.io/code-generator v0.18.3
	sigs.k8s.io/controller-tools v0.3.0
	sigs.k8s.io/testing_frameworks v0.1.2
)

replace github.com/criticalstack/crit => ../../
