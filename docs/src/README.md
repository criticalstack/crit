<div style="float: right; position: relative; display: inline; margin: 30px;"><img src="images/crit-md.png" width=200></div>

# Crit

**Crit** is a command-line tool for bootstrapping Kubernetes clusters. It handles the initial configuration of Kubernetes control plane components, and adding workers to the cluster.

It is designed to be used within automated scripting (i.e. non-interactive). Many providers of virtual infrastructure allow user-defined customization via shell script, which ensures Crit composes well with provider provisioning tools (e.g. AWS Cloudformation).


## Installation

The easiest way to install:

```sh
curl -sSfL https://get.crit.sh | sh
```

Pre-built binaries are also available in [Releases](https://github.com/criticalstack/crit/releases/latest). Crit is written in Go so it is also pretty simple to install via go:

```sh
go get -u github.com/criticalstack/crit/cmd/crit
```

RPM/Debian packages are also available via [packagecloud.io](https://packagecloud.io/criticalstack/public).

## Requirements

Crit is a standalone binary, however, there are implied requirements that aren't as straight-forward. Be sure to check out the [Getting Started](getting-started.md).
