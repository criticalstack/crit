## cinder create node

Creates a new cinder worker

### Synopsis

Creates a new cinder worker

```
cinder create node [flags]
```

### Options

```
  -c, --config string       cinder configuration file
  -h, --help                help for node
      --image string        node image (default "criticalstack/cinder:v1.0.0-beta.1")
      --kubeconfig string   sets kubeconfig path instead of $KUBECONFIG or $HOME/.kube/config
      --name string         cluster name (default "cinder")
```

### Options inherited from parent commands

```
  -v, --verbose count   log output verbosity
```

### SEE ALSO

* [cinder create](cinder-create.md)	 - Create cinder resources

